package controller

import (
	"context"
	"fmt"

	"github.com/matt-empson/pod-labeller/internal/config"
	"github.com/matt-empson/pod-labeller/internal/kube"
	coreinformers "k8s.io/client-go/informers/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type Controller struct {
	client    *kube.Client
	lister    corelisters.PodLister
	informer  cache.SharedIndexInformer
	queue     workqueue.TypedRateLimitingInterface[string]
	namespace string
	labels    map[string]string
}

func NewController(client *kube.Client, informer coreinformers.PodInformer, cfg *config.Config) *Controller {
	queue := workqueue.NewTypedRateLimitingQueueWithConfig(workqueue.DefaultTypedControllerRateLimiter[string](), workqueue.TypedRateLimitingQueueConfig[string]{Name: "pod-labeller"})

	c := &Controller{
		client:    client,
		lister:    informer.Lister(),
		informer:  informer.Informer(),
		queue:     queue,
		namespace: cfg.Namespace,
		labels:    cfg.Labels,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj any) { c.enqueuePod(obj) },
		UpdateFunc: func(oldObj, newObj any) { c.enqueuePod(newObj) },
	})

	return c
}

func (c *Controller) Run(ctx context.Context) error {
	defer c.queue.ShutDown()

	if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
		return fmt.Errorf("failed to sync informer cache")
	}

	klog.Info("controller started, watching pods")

	for range 1 {
		go c.runWorker(ctx)
	}

	<-ctx.Done()
	klog.Info("controller shutting down")

	return nil
}
