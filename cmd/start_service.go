package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/xjasmx/kpfm/pkg/kubernetes"
)

var startServiceCmd = &cobra.Command{
	Use:   "service [alias]",
	Short: "Start port forwarding for a service",
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

		go startPortForwarding(ctx, clientset, restConfig, alias)

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
