package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/xjasmx/kpfm/pkg/config"
	"github.com/xjasmx/kpfm/pkg/kubernetes"
)

var startGroupCmd = &cobra.Command{
	Use:   "group [alias]",
	Short: "Start port forwarding for a group services",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		clientset, restConfig, err := kubernetes.ConnectKubernetes()
		if err != nil {
			return fmt.Errorf("error connecting to Kubernetes: %w", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		groupConfig, err := config.LoadGroup(alias)
		if err != nil {
			return fmt.Errorf("error reading group config: %w", err)
		}
		for _, serviceAlias := range groupConfig.Config.Services {
			go startPortForwarding(ctx, clientset, restConfig, serviceAlias)
		}

		select {
		case <-sigs:
			fmt.Println("Received an interrupt, stopping services...")
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping services...")
		}

		time.Sleep(1 * time.Second)
		return nil
	},
}
