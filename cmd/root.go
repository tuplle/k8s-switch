package cmd

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tuplle/k8s-switch/internal"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-switch",
	Short: "Quickly switch between Kubernetes configurations",
	Long: `k8s-switch is a CLI tool that helps you manage multiple kubeconfig files.
It lists files from your ~/.kube/conf.d directory, allows you to select one
via an interactive prompt, and automatically updates your main ~/.kube/config.

Optionally, it can launch k9s directly using the selected configuration.
Optionally, it can open logs in default browser with kubetail.`,
	Run: func(cmd *cobra.Command, args []string) {
		k9s, _ := cmd.Flags().GetBool("k9s")
		k9sOnly, _ := cmd.Flags().GetBool("k9s-only")
		logs, _ := cmd.Flags().GetBool("logs")
		configDir, _ := cmd.Flags().GetString("dir")
		homedir, _ := os.UserHomeDir()
		if configDir == "" {
			configDir = filepath.Join(homedir, ".kube", "config.d")
		}

		configs, err := internal.GetFilesFromDir(configDir)
		if err != nil {
			panic(err)
		}

		filteredConfigs := make(map[string]string)
		for name, path := range configs {
			if filepath.Ext(name) == ".yaml" {
				filteredConfigs[name] = path
			}
		}
		configs = filteredConfigs

		prompt := promptui.Select{
			Label: "Select config",
			Items: slices.Sorted(maps.Keys(configs)),
			Size:  10,
		}
		_, result, err := prompt.Run()
		if err != nil {
			if errors.Is(err, promptui.ErrInterrupt) {
				fmt.Println("Selection cancelled")
				return
			}
			panic(err)
		}

		if !k9sOnly {
			err = internal.CopyFile(configs[result], filepath.Join(homedir, ".kube", "config"))
			if err != nil {
				panic(err)
			}
		}

		if k9s || k9sOnly {
			err = runK9s(configs[result])
			if err != nil {
				panic(err)
			}
		}

		if logs {
			err = runKubetail(configs[result])
			if err != nil {
				panic(err)
			}
		}
	},
}

func Execute() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("BEY BEY")
		os.Exit(0)
	}()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolP("k9s", "9", false, "launch k9s with selected config")
	rootCmd.Flags().Bool("k9s-only", false, "launch k9s with selected config only")
	rootCmd.Flags().Bool("logs", false, "open logs in default browser with kubetail")
	rootCmd.Flags().StringP("dir", "d", "", "path to config.d directory")
}

func runK9s(kubeconfigPath string) error {
	return internal.RunAnotherTUI("k9s", []string{"--kubeconfig", kubeconfigPath})
}

func runKubetail(s string) error {
	return internal.RunAnotherTUI("kubetail", []string{"serve", "-c", s})
}
