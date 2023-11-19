package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xjasmx/kpfm/pkg/config"
	"strings"
)

var createGroupCmd = &cobra.Command{
	Use:   "group [alias] --services [service-alias-list]",
	Short: "Create a new configuration for a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		services, _ := cmd.Flags().GetString("services")
		serviceList := strings.Split(services, ",")
		for i := range serviceList {
			serviceList[i] = strings.TrimSpace(serviceList[i])
		}

		group := config.NewGroup(alias, serviceList)
		if err := group.Save(); err != nil {
			return fmt.Errorf("error creating group configuration: %w", err)
		}

		fmt.Println("Group configuration created successfully")
		return nil
	},
}

func init() {
	createGroupCmd.Flags().StringP("services", "s", "", "Define or reference a list of service aliases when creating a group")
	createGroupCmd.MarkFlagRequired("services")
}
