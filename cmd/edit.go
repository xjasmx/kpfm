package cmd

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing configuration",
}

func init() {
	editCmd.AddCommand(editGroupCmd, editServiceCmd)
}
