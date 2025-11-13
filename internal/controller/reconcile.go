package controller

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func (c *Controller) reconcile(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "could not parse key", "key", key)
		return nil // do not retry malformed keys
	}

	pod, err := c.lister.Pods(namespace).Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.InfoS("pod not found, must have been deleted", "pod", key)
			return nil
		}

		return fmt.Errorf("error fetching pod: %w", err)
	}

	if pod.DeletionTimestamp != nil {
		klog.InfoS("pod is being deleted, skipping", "pod", pod.Name)
		return nil
	}

	patch, needsPatch := computeLabelPatch(pod.Labels, c.labels)
	if !needsPatch {
		return nil
	}

	klog.InfoS("labels not in desired state - reconciling", "pod", pod.Name, "labels", patch)

	patchBody := map[string]any{
		"metadata": map[string]any{
			"labels": patch,
		},
	}

	patchJSON, err := json.Marshal(patchBody)
	if err != nil {
		return fmt.Errorf("error marshalling label patch: %w", err)
	}

	_, err = c.client.ClientSet.CoreV1().Pods(namespace).Patch(ctx, name, types.MergePatchType, patchJSON, v1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("patch failed: %w", err)
	}

	return nil
}

func computeLabelPatch(existing, desired map[string]string) (map[string]string, bool) {
	patch := make(map[string]string)

	for k, v := range desired {
		if existingValue, ok := existing[k]; !ok || existingValue != v {
			patch[k] = v
		}
	}

	if len(patch) == 0 {
		return nil, false
	}

	return patch, true
}
