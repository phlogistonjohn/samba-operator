/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	destPath      = "/var/lib/svcwatch/status.json"
	svcLabelKey   = ""
	svcLabelValue = ""
	svcNamespace  = ""
)

type HostInfo struct {
	Name        string `json:"name"`
	IPv4Address string `json:"ipv4"`
	Target      string `json:"target"`
}

type HostState struct {
	Reference string     `json:"ref"`
	Items     []HostInfo `json:"items"`
}

func (hs HostState) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(hs)
}

func (hs HostState) Differs(other HostState) bool {
	if hs.Reference != other.Reference {
		return true
	}
	if len(hs.Items) != len(other.Items) {
		return true
	}
	for i := range hs.Items {
		if hs.Items[i] != other.Items[i] {
			return true
		}
	}
	return false
}

func serviceToHostState(svc *corev1.Service, nameLabel string) HostState {
	name := svc.Labels[nameLabel]
	if name == "" {
		// name label not found, fall back to object name
		name = svc.Name
	}
	hs := HostState{
		Reference: fmt.Sprintf("k8s: %s service/%s", svc.Namespace, svc.Name),
		Items:     []HostInfo{},
	}
	for i, ig := range svc.Status.LoadBalancer.Ingress {
		var n string
		if i == 0 {
			n = name
		} else {
			n = fmt.Sprintf("%s-%d", name, i)
		}
		hs.Items = append(hs.Items, HostInfo{
			Name: n,
			IPv4Address: ig.IP,
			Target: "external",
		})
	}
	if svc.Spec.ClusterIP != "" {
		hs.Items = append(hs.Items, HostInfo{
			Name: fmt.Sprintf("%s-cluster", name),
			IPv4Address: svc.Spec.ClusterIP,
			Target: "internal",
		})
	}
	return hs
}

func serviceUpdate(prev HostState, svc *corev1.Service, nameLabel string) (HostState, bool) {
	newh := serviceToHostState(svc, nameLabel)
	return newh, prev.Differs(newh)
}

func processUpdates(
	path, nameLabel string,
	updates <-chan *corev1.Service, errors chan<- error) {
	// initialize host state to the zero value
	var hs HostState
	changed := true
	for svc := range updates {
		if svc == nil {
			return
		}
		hs, changed = serviceUpdate(hs, svc, nameLabel)
		if changed {
			err := hs.Save(path)
			if err != nil {
				errors <- err
			}
		}
	}
}

func envConfig(d *string, vn string) {
	ev := os.Getenv(vn)
	if ev != "" {
		*d = ev
	}
}

func main() {
	envConfig(&destPath, "DESTINATION_PATH")
	envConfig(&svcLabelKey, "SERVICE_LABEL_KEY")
	envConfig(&svcLabelValue, "SERVICE_LABEL_VALUE")
	envConfig(&svcNamespace, "SERVICE_NAMESPACE")

	flag.StringVar(&destPath, "destination", "", "JSON file to update")
	flag.StringVar(&svcLabelKey, "label-key", "", "Label key to watch")
	flag.StringVar(&svcLabelValue, "label-value", "", "Label value")
	flag.StringVar(&svcNamespace, "namespace", "", "Namespace")
	flag.Parse()

	l, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("failed to set up logger\n")
		os.Exit(2)
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	kcfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		l.Error("failed to create watch", zap.Error(err))
		os.Exit(1)
	}
	clientset := kubernetes.NewForConfigOrDie(kcfg)

	sel := fmt.Sprintf("%s=%s", svcLabelKey, svcLabelValue)
	w, err := clientset.CoreV1().Services(svcNamespace).Watch(
		context.TODO(),
		metav1.ListOptions{LabelSelector: sel})
	if err != nil {
		l.Error("failed to create watch", zap.Error(err))
		os.Exit(1)
	}

	errors := make(chan error, 1)
	updates := make(chan *corev1.Service)
	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch,
		os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	ready := true
	go processUpdates(destPath, svcLabelKey, updates, errors)
	for ready {
		select {
		case v := <-w.ResultChan():
			l.Info("new result from watch", zap.Any("value", v))
			svc, ok := v.Object.(*corev1.Service)
			if !ok {
				l.Error("got non-service from service watch")
				ready = false
				break
			}
			l.Info("updated service", zap.Any("service", svc))
			updates <- svc
		case <-signalch:
			l.Info("terminating")
			ready = false
			break
		case e := <-errors:
			l.Error("error updating host state", zap.Error(e))
			ready = false
			break
		}
	}
	updates <- nil
	close(updates)
	close(errors)
	w.Stop()
}
