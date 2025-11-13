package controller

import (
	"context"
	"encoding/json"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func (c *Controller) reconcile(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "unexpected error")
		return err
	}

	pod, err := c.lister.Pods(namespace).Get(name)
	if err != nil {
		klog.ErrorS(err, "unexpected error")
		return err
	}

	for k := range c.labels {
		_, ok := pod.Labels[k]
		if ok {
			klog.InfoS("label exists", "label", k)
		}

		if !ok {
			klog.InfoS("label does not exist", "label", k)
			patchBody := map[string]any{
				"metadata": map[string]any{
					"labels": c.labels,
				},
			}

			patchBytes, _ := json.Marshal(patchBody)
			c.client.ClientSet.CoreV1().Pods(namespace).Patch(ctx, name, types.MergePatchType, patchBytes, v1.PatchOptions{})
		}
	}
	return nil
}
