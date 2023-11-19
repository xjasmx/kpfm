package cmd

import (
	"fmt"
	"github.com/xjasmx/kpfm/pkg/config"

	"github.com/spf13/cobra"
)

var deleteGroupCmd = &cobra.Command{
	Use:   "group [alias]",
	Short: "Remove a saved configuration for a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		cm := config.NewConfigManager()
		if err := cm.DeleteGroupConfig(alias); err != nil {
			return fmt.Errorf("error deleting group configuration: %w", err)
		}
		fmt.Println("Group configuration deleted successfully")
		return nil
	},
}
