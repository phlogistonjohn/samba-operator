// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types" // nolint:typecheck
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	sambaoperatorv1alpha1 "github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
)

type obj[T any] interface {
	*T
	rtclient.Object
}

type autoClient[T any, O obj[T]] struct {
	client rtclient.Client
	logger Logger
	scheme *runtime.Scheme
}

func clientFromMgr[T any, O obj[T]](m *SmbShareManager) *autoClient[T, O] {
	return &autoClient[T, O]{client: m.client, logger: m.logger}
}

func (f *autoClient[T, O]) Get(
	ctx context.Context,
	name types.NamespacedName) (*T, error) {
	// ---
	t := new(T)
	err := f.client.Get(ctx, name, O(t))
	return t, err
}

func (f *autoClient[T, O]) GetExisting(
	ctx context.Context,
	name types.NamespacedName) (*T, error) {
	// ---
	t := new(T)
	err := f.client.Get(ctx, name, O(t))
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		return nil, nil
	}
	return t, nil
}

func (f *autoClient[T, O]) GetOrCreate(
	ctx context.Context,
	name types.NamespacedName,
	defaultResource func() *T,
	owner rtclient.Object) (*T, bool, error) {
	// ---
	t, err := f.GetExisting(ctx, name)
	if err != nil {
		return nil, false, err
	}
	if t != nil {
		return t, false, nil
	}

	t = defaultResource()
	O(t).SetNamespace(name.Namespace)
	O(t).SetName(name.Name)
	if owner != nil {
		err = controllerutil.SetControllerReference(owner, O(t), f.scheme)
		if err != nil {
			f.logger.Error(
				err,
				"Failed to set controller reference",
				"SmbShare.Namespace", owner.GetNamespace(),
				"SmbShare.Name", owner.GetName(),
				"Object.Namespace", O(t).GetNamespace(),
				"Object.Name", O(t).GetName())
			return t, false, err
		}
	}
	f.logger.Info(
		"Creating a new StatefulSet",
		"Object.Namespace", O(t).GetNamespace(),
		"Object.Name", O(t).GetName())
	err = f.client.Create(ctx, O(t))
	if err != nil {
		f.logger.Error(
			err,
			"Failed to create new StatefulSet",
			"Object.Namespace", O(t).GetNamespace(),
			"Object.Name", O(t).GetName())
		return t, false, err
	}
	return t, true, err
}

func (m *SmbShareManager) play1(
	ctx context.Context,
	smbShare *sambaoperatorv1alpha1.SmbShare,
	ns string) (*corev1.PersistentVolumeClaim, bool, error) {
	// ---
	client := clientFromMgr[corev1.PersistentVolumeClaim](m)
	pvc, err := client.Get(ctx, types.NamespacedName{
		Namespace: "foo",
		Name:      "bar",
	})

	return pvc, false, err
}
