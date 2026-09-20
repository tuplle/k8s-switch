package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func writeKubeconfigFixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write fixture file: %v", err)
	}
	return path
}

func TestReadKubeconfigSummary(t *testing.T) {
	path := writeKubeconfigFixture(t, `
apiVersion: v1
kind: Config
clusters:
  - name: my-cluster
    cluster:
      server: https://cluster.example.com:6443
contexts:
  - name: my-context
    context:
      cluster: my-cluster
      user: my-user
current-context: my-context
users:
  - name: my-user
    user: {}
`)

	summary, err := ReadKubeconfigSummary(path)
	if err != nil {
		t.Fatalf("ReadKubeconfigSummary returned error: %v", err)
	}
	if summary.ContextName != "my-context" {
		t.Errorf("expected context name %q, got %q", "my-context", summary.ContextName)
	}
	if summary.ServerAddress != "https://cluster.example.com:6443" {
		t.Errorf("expected server %q, got %q", "https://cluster.example.com:6443", summary.ServerAddress)
	}
}

func TestReadKubeconfigSummaryMissingFile(t *testing.T) {
	_, err := ReadKubeconfigSummary(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadKubeconfigSummaryMalformedYAML(t *testing.T) {
	path := writeKubeconfigFixture(t, "clusters: [this is not valid: yaml")

	_, err := ReadKubeconfigSummary(path)
	if err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}

func TestReadKubeconfigSummaryEmptySections(t *testing.T) {
	path := writeKubeconfigFixture(t, `
apiVersion: v1
kind: Config
`)

	summary, err := ReadKubeconfigSummary(path)
	if err != nil {
		t.Fatalf("ReadKubeconfigSummary returned error: %v", err)
	}
	if summary.ContextName != "" {
		t.Errorf("expected empty context name, got %q", summary.ContextName)
	}
	if summary.ServerAddress != "" {
		t.Errorf("expected empty server address, got %q", summary.ServerAddress)
	}
}
