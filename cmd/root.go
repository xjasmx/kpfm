package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kpfm",
	Short: "Kubectl Port Forward Manager (KPFM)",
	Long:  `KPFM is a CLI tool to simplify and manage kubectl port-forward commands for Kubernetes services.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kubectl Port Forward Manager (KPFM)")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd, editCmd, deleteCmd, startCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
