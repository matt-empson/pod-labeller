package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewConfigFromFlags(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantCfg   *Config
		wantErr   bool
		errSubstr string
	}{
		{
			name: "no_flags",
			wantCfg: &Config{
				Namespace:  "default",
				Labels:     map[string]string{"example.com/managed-by": "pod-labeller"},
				Kubeconfig: "",
				LogLevel:   "info",
			},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name: "custom_namespace_kubeconfig",
			args: []string{
				"--namespace=production",
				"--kubeconfig=/tmp/kubeconfig",
			},
			wantCfg: &Config{
				Namespace:  "production",
				Labels:     map[string]string{"example.com/managed-by": "pod-labeller"},
				Kubeconfig: "/tmp/kubeconfig",
				LogLevel:   "info",
			},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name: "single_label",
			args: []string{
				"--labels=team=test",
			},
			wantCfg: &Config{
				Namespace:  "default",
				Labels:     map[string]string{"team": "test"},
				Kubeconfig: "",
				LogLevel:   "info",
			},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name: "multi_label",
			args: []string{
				"--labels=team=test,env=test",
			},
			wantCfg: &Config{
				Namespace: "default",
				Labels: map[string]string{
					"team": "test",
					"env":  "test",
				},
				Kubeconfig: "",
				LogLevel:   "info",
			},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name: "invalid_label",
			args: []string{
				"--labels=invalid",
			},
			wantCfg: &Config{
				Namespace:  "default",
				Labels:     map[string]string{},
				Kubeconfig: "",
				LogLevel:   "info",
			},
			wantErr:   true,
			errSubstr: "invalid label \"invalid\": must be key=value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := NewConfigFromFlags(test.args)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}

				if !strings.Contains(err.Error(), test.errSubstr) {
					t.Errorf("got: %q\n\nwant: %q", err, test.errSubstr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error %q", err)
			}

			if !reflect.DeepEqual(cfg, test.wantCfg) {
				t.Errorf("got: %+v\n\nwant: %+v", cfg, test.wantCfg)
			}

		})
	}
}
