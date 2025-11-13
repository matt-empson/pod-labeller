package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Config struct {
	Namespace  string
	Labels     map[string]string
	Kubeconfig string
	LogLevel   string
}

func NewConfigFromFlags(args []string) (*Config, error) {
	fs := pflag.NewFlagSet("pod-labeller", pflag.ExitOnError)

	namespace := fs.String("namespace", "default", "Namespace for controller to watch")
	labels := fs.StringSlice("labels", []string{"example.com/managed-by=pod-labeller"}, "Labels to add to pods")
	kubeconfig := fs.String("kubeconfig", "", "Path to kubeconfig (optional)")
	logLevel := fs.String("log-level", "info", "Log level (debug, info, warn, error)")

	fs.Parse(args)

	*logLevel = strings.ToLower(*logLevel)
	allowedLogLevels := sets.NewString("debug", "info", "warn", "error")
	if !allowedLogLevels.Has(*logLevel) {
		return nil, fmt.Errorf("log-level %q: must be one of %v", *logLevel, allowedLogLevels.List())
	}

	if namespace == nil || strings.TrimSpace(*namespace) == "" {
		return nil, fmt.Errorf("namespace cannot be empty; please specify a valid namespace for the controller to watch")
	}

	labelMap := make(map[string]string)

	for _, v := range *labels {
		labelData := strings.SplitN(v, "=", 2)
		if len(labelData) != 2 {
			return nil, fmt.Errorf("invalid label %q: must be key=value", v)
		}

		labelKey := strings.TrimSpace(labelData[0])
		labelValue := strings.TrimSpace(labelData[1])

		if _, ok := labelMap[labelKey]; ok {
			return nil, fmt.Errorf("duplicate label key %q", labelKey)
		}

		labelMap[labelKey] = labelValue
	}

	conf := Config{
		Namespace:  *namespace,
		Labels:     labelMap,
		Kubeconfig: *kubeconfig,
		LogLevel:   *logLevel,
	}

	return &conf, nil
}
