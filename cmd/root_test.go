package cmd

import "testing"

func TestRootCmdVersion(t *testing.T) {
	if rootCmd.Version != Version {
		t.Errorf("expected rootCmd.Version to be %q, got %q", Version, rootCmd.Version)
	}
}

func TestRootCmdFlags(t *testing.T) {
	for _, name := range []string{"k9s", "k9s-only", "logs", "dir"} {
		if rootCmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
	if rootCmd.PersistentFlags().Lookup("verbose") == nil {
		t.Error("expected persistent flag \"verbose\" to be registered")
	}
}
