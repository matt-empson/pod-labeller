package kube

import (
	"errors"
	"testing"

	"github.com/matt-empson/pod-labeller/internal/config"
	"k8s.io/client-go/rest"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name                 string
		inClusterConfig      func() (*rest.Config, error)
		buildConfigFromFlags func(masterUrl, kubeconfigPath string) (*rest.Config, error)
		expectErr            bool
		expectNilClient      bool
	}{
		{
			name: "in-cluster success",
			inClusterConfig: func() (*rest.Config, error) {
				return &rest.Config{}, nil
			},
			buildConfigFromFlags: func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
				return nil, nil
			},
			expectErr:       false,
			expectNilClient: false,
		},
		{
			name: "kubeconfig success",
			inClusterConfig: func() (*rest.Config, error) {
				return nil, errors.New("in-cluster config failed, falling back to kubeconfig")
			},
			buildConfigFromFlags: func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
				return &rest.Config{}, nil
			},
			expectErr:       false,
			expectNilClient: false,
		},
		{
			name: "in-cluster and kubeconfig failure",
			inClusterConfig: func() (*rest.Config, error) {
				return nil, errors.New("in-cluster config failed, falling back to kubeconfig")
			},
			buildConfigFromFlags: func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
				return &rest.Config{}, errors.New("kubeconfig not found")
			},
			expectErr:       true,
			expectNilClient: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf := &config.Config{Kubeconfig: "mock"}

			clientBuilder := &ClientBuilder{
				InClusterConfig:      test.inClusterConfig,
				BuildConfigFromFlags: test.buildConfigFromFlags,
			}

			client, err := clientBuilder.NewClient(conf)

			if test.expectErr && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !test.expectErr && err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			if test.expectNilClient && client != nil {
				t.Fatalf("expected nil-client but got %#v", client)
			}

			if !test.expectNilClient && client == nil {
				t.Fatalf("expected client but got nil")
			}

		})
	}
}
