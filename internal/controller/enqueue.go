package controller

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func (c *Controller) enqueuePod(obj any) {
	objRef, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.ErrorS(err, "error getting namespace/key for pod")
		return
	}

	namespace, _, _ := cache.SplitMetaNamespaceKey(objRef)

	if namespace == c.namespace {
		klog.InfoS("enqueued pod", "pod", objRef)
		c.queue.Add(objRef)
	}
}
