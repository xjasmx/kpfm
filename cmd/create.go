package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new configuration",
}

func init() {
	createCmd.AddCommand(createGroupCmd, createServiceCmd)
}
