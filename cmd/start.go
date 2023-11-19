package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xjasmx/kpfm/pkg/config"
	"github.com/xjasmx/kpfm/pkg/kubernetes"
	offKubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a port forwarding session",
}

func startPortForwarding(ctx context.Context, clientset *offKubernetes.Clientset, restConfig *rest.Config, alias string) {
	serviceConfig, err := config.LoadService(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading service config: %v\n", err)
		return
	}
	pf, err := kubernetes.NewPortForwarder(clientset,
		restConfig,
		serviceConfig.Config.Namespace,
		serviceConfig.Config.ServiceName,
		serviceConfig.Config.PortMapping)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up port forwarding: %v\n", err)
		return
	}
	if err := pf.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error during port forwarding: %v\n", err)
	}
}

func init() {
	startCmd.AddCommand(startGroupCmd, startServiceCmd)
}
