package controller

import (
	"reflect"
	"testing"
)

func TestComputeLabelPatch(t *testing.T) {
	desiredLabels := map[string]string{
		"example.com/managed-by": "pod-labeller",
		"test-label":             "test-value",
	}

	tests := []struct {
		name               string
		existing           map[string]string
		desired            map[string]string
		expectedPatch      map[string]string
		expectedNeedsPatch bool
	}{
		{
			name:               "in-sync",
			existing:           desiredLabels,
			desired:            desiredLabels,
			expectedPatch:      nil,
			expectedNeedsPatch: false,
		},
		{
			name:               "missing_label",
			existing:           map[string]string{"example.com/managed-by": "pod-labeller"},
			desired:            desiredLabels,
			expectedPatch:      map[string]string{"test-label": "test-value"},
			expectedNeedsPatch: true,
		},
		{
			name: "incorrect_value",
			existing: map[string]string{
				"example.com/managed-by": "pod-labeller",
				"test-label":             "incorrect-value",
			},
			desired:            desiredLabels,
			expectedPatch:      map[string]string{"test-label": "test-value"},
			expectedNeedsPatch: true,
		},
		{
			name:               "nil_labels",
			existing:           nil,
			desired:            desiredLabels,
			expectedPatch:      desiredLabels,
			expectedNeedsPatch: true,
		},
		{
			name: "multi_incorrect_value",
			existing: map[string]string{
				"example.com/managed-by": "incorrect-value",
				"test-label":             "incorrect-value",
			},
			desired:            desiredLabels,
			expectedPatch:      desiredLabels,
			expectedNeedsPatch: true,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			patch, needsPatch := computeLabelPatch(v.existing, v.desired)

			if needsPatch != v.expectedNeedsPatch {
				t.Fatalf("\ngot: %v\n\nwant: %v", needsPatch, v.expectedNeedsPatch)
			}

			if !reflect.DeepEqual(patch, v.expectedPatch) {
				t.Errorf("\ngot: %v\n\nwant: %v", patch, v.expectedPatch)
			}
		})
	}
}
