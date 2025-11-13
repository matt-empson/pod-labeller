package controller

import (
	"context"

	"k8s.io/klog/v2"
)

func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	obj, shutdown := c.queue.Get()
	defer c.queue.Done(obj)

	if shutdown {
		klog.Info("shutting down processor")
		return false
	}

	err := c.reconcile(ctx, obj)
	if err != nil {
		klog.ErrorS(err, "failed to reconcile", "obj", obj)
		c.queue.AddRateLimited(obj)
		return true
	}

	c.queue.Forget(obj)
	return true
}

func (c *Controller) runWorker(ctx context.Context) {
	for {
		if !c.processNextWorkItem(ctx) {
			return
		}
	}
}
