package cmd

import (
	"fmt"
	"github.com/xjasmx/kpfm/pkg/config"

	"github.com/spf13/cobra"
)

var editServiceCmd = &cobra.Command{
	Use:   "service [alias] --port [portMapping] --namespace [namespace] --service [serviceName]",
	Short: "Edit an existing configuration for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		newPortMapping, _ := cmd.Flags().GetString("port")
		newNamespace, _ := cmd.Flags().GetString("namespace")
		newServiceName, _ := cmd.Flags().GetString("service")

		service, err := config.LoadService(alias)
		if err != nil {
			return fmt.Errorf("error loading service configuration: %w", err)
		}

		if newNamespace != "" {
			service.Config.Namespace = newNamespace
		}
		if newServiceName != "" {
			service.Config.ServiceName = newServiceName
		}
		if newPortMapping != "" {
			service.Config.PortMapping = newPortMapping
		}

		if err := service.Save(); err != nil {
			return fmt.Errorf("error editing service configuration: %w", err)
		}
		fmt.Println("Service configuration updated successfully")
		return nil
	},
}

func init() {
	editServiceCmd.Flags().StringP("namespace", "n", "", "Specify the Kubernetes namespace for the service")
	editServiceCmd.Flags().StringP("service", "s", "", "Specify the name of the Kubernetes service for the session")
	editServiceCmd.Flags().StringP("port", "p", "", "Specify the port mapping of the Kubernetes service for the session")
}
