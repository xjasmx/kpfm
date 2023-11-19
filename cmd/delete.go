package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a saved configuration",
}

func init() {
	deleteCmd.AddCommand(deleteGroupCmd, deleteServiceCmd)
}
