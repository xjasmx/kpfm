package utils

import (
	"github.com/spf13/cobra"
	"log"
)

// MustMarkFlagRequired ensures that a flag is marked as required, and exits if an error occurs.
func MustMarkFlagRequired(cmd *cobra.Command, flagName string) {
	if err := cmd.MarkFlagRequired(flagName); err != nil {
		log.Fatalf("Error marking '%s' flag as required: %v", flagName, err)
	}
}
