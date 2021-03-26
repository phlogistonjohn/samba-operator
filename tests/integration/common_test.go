// +build integration

package integration

import (
	"github.com/spf13/viper"
)

type testConfig struct {
	Namespace     string `mapstructure:"namespace"`
	FilesDir      string `mapstructure:"files_dir"`
	ConfigDir     string `mapstructure:"config_dir"`
	KustomizeCmd  string `mapstructure:"kustomize"`
	ExpectedImage string `mapstructure:"expected_img"`
}

var config testConfig

func init() {
	must := func(e error) {
		if e != nil {
			panic(e)
		}
	}

	v := viper.New()
	v.SetDefault("namespace", "samba-operator-system")
	v.SetDefault("files_dir", "../files")
	v.SetDefault("config_dir", "../../config")
	v.SetDefault("kustomize", "kustomize")
	v.SetDefault("expected_img", "quay.io/samba.org/samba-operator:latest")

	v.SetEnvPrefix("SMBOP_TEST")
	must(v.BindEnv("namespace"))
	must(v.BindEnv("files_dir"))
	must(v.BindEnv("config_dir"))
	must(v.BindEnv("expected_img"))
	// kustomize is special so that the env var can be shared among
	// more programs than just our test suite.
	// in future versions of viper we can specify multiple env vars
	// as "aliases"
	must(v.BindEnv("kustomize", "KUSTOMIZE"))

	must(v.Unmarshal(&config))
}
