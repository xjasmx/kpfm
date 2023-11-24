package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xjasmx/kpfm/pkg/config"
	"github.com/xjasmx/kpfm/pkg/utils"
)

var createServiceCmd = &cobra.Command{
	Use:   "service [alias] --port [portMapping] --namespace [namespace] --service [serviceName]",
	Short: "Create a new configuration for a service or a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		portMapping, _ := cmd.Flags().GetString("port")
		namespace, _ := cmd.Flags().GetString("namespace")
		serviceName, _ := cmd.Flags().GetString("service")

		service := config.NewService(alias, namespace, serviceName, portMapping)
		if err := service.Save(); err != nil {
			return fmt.Errorf("error creating service configuration: %w", err)
		}
		fmt.Println("Service configuration created successfully")
		return nil
	},
}

func init() {
	createServiceCmd.Flags().StringP("namespace", "n", "", "Specify the Kubernetes namespace for the service")
	createServiceCmd.Flags().StringP("service", "s", "", "Specify the name of the Kubernetes service for the session")
	createServiceCmd.Flags().StringP("port", "p", "", "Specify the port mapping of the Kubernetes service for the session")

	utils.MustMarkFlagRequired(createServiceCmd, "namespace")
	utils.MustMarkFlagRequired(createServiceCmd, "service")
	utils.MustMarkFlagRequired(createServiceCmd, "port")
}
