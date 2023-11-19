package cmd

import (
	"fmt"
	"github.com/xjasmx/kpfm/pkg/config"

	"github.com/spf13/cobra"
)

var deleteServiceCmd = &cobra.Command{
	Use:   "service [alias]",
	Short: "Remove a saved configuration for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		cm := config.NewConfigManager()
		if err := cm.DeleteServiceConfig(alias); err != nil {
			return fmt.Errorf("error deleting service configuration: %w", err)
		}
		fmt.Println("Service configuration deleted successfully")
		return nil
	},
}
