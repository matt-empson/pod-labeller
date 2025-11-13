package kube

import (
	"context"
	"time"

	"github.com/matt-empson/pod-labeller/internal/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/klog/v2"
)

type Client struct {
	ClientSet  *kubernetes.Clientset
	RestConfig *rest.Config
}

type ClientBuilder struct {
	InClusterConfig      func() (*rest.Config, error)
	BuildConfigFromFlags func(masterUrl, kubeconfigPath string) (*rest.Config, error)
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		InClusterConfig:      rest.InClusterConfig,
		BuildConfigFromFlags: clientcmd.BuildConfigFromFlags,
	}
}

// NewClient creates a Kubernetes API client using in-cluster or kubeconfig auth.
func (c *ClientBuilder) NewClient(cfg *config.Config) (*Client, error) {
	klog.InfoS("attempting in-cluster config")

	kubeConfig, err := c.InClusterConfig()
	if err != nil {
		klog.InfoS("in-cluster config failed, falling back to kubeconfig", "path", cfg.Kubeconfig)
		kubeConfig, err = c.BuildConfigFromFlags("", cfg.Kubeconfig)
		if err != nil {
			klog.ErrorS(err, "failed to load kubeconfig")
			return nil, err
		}
	}

	kubeConfig.QPS = 100
	kubeConfig.Burst = 200

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		klog.ErrorS(err, "failed to create Kubernetes client")
		return nil, err
	}

	klog.InfoS("successfully created Kubernetes client")
	return &Client{ClientSet: client, RestConfig: kubeConfig}, nil
}

func (c *Client) CheckConnection(ctx context.Context, namespace string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.ClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		klog.ErrorS(err, "connection failed")
		return err
	}

	klog.InfoS("connection successful")
	return nil
}
