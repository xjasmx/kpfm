package cmd

import (
	"fmt"
	"github.com/xjasmx/kpfm/pkg/config"
	"strings"

	"github.com/spf13/cobra"
)

var editGroupCmd = &cobra.Command{
	Use:   "group [alias] --add [service-alias-list] --remove [service-alias-list]",
	Short: "Edit an existing configuration for a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		addServices, _ := cmd.Flags().GetString("add")
		removeServices, _ := cmd.Flags().GetString("remove")

		group, err := config.LoadGroup(alias)
		if err != nil {
			return fmt.Errorf("error loading group configuration: %w", err)
		}

		if addServices != "" {
			addList := strings.Split(addServices, ",")
			group.Config.Services = append(group.Config.Services, addList...)
		}

		if removeServices != "" {
			removeList := strings.Split(removeServices, ",")
			for _, remove := range removeList {
				for i, service := range group.Config.Services {
					if service == remove {
						group.Config.Services = append(group.Config.Services[:i], group.Config.Services[i+1:]...)
						break
					}
				}
			}
		}

		if err := group.Save(); err != nil {
			return fmt.Errorf("error editing group configuration: %w", err)
		}
		fmt.Println("Group configuration updated successfully")
		return nil
	},
}

func init() {
	editGroupCmd.Flags().StringP("add", "a", "", "List of service aliases to add to the group")
	editGroupCmd.Flags().StringP("remove", "r", "", "List of service aliases to remove from the group")
}
