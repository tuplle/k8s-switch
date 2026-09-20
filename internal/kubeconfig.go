package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

// KubeconfigSummary holds the display-relevant fields read from a kubeconfig file by
// ReadKubeconfigSummary.
type KubeconfigSummary struct {
	ContextName   string
	ServerAddress string
}

type kubeconfigFile struct {
	Clusters []struct {
		Cluster struct {
			Server string `yaml:"server"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Contexts []struct {
		Name string `yaml:"name"`
	} `yaml:"contexts"`
}

// ReadKubeconfigSummary reads path as a kubeconfig YAML file and returns the name of
// its context and the server address of its cluster, assuming exactly one of each
// (only the first entry of "clusters" and "contexts" is read; nothing is merged and
// "current-context" is not consulted). Either field is empty if that section is
// missing or empty in the file.
func ReadKubeconfigSummary(path string) (KubeconfigSummary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return KubeconfigSummary{}, err
	}

	var cfg kubeconfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return KubeconfigSummary{}, err
	}

	var summary KubeconfigSummary
	if len(cfg.Contexts) > 0 {
		summary.ContextName = cfg.Contexts[0].Name
	}
	if len(cfg.Clusters) > 0 {
		summary.ServerAddress = cfg.Clusters[0].Cluster.Server
	}
	return summary, nil
}
