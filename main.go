package main

import (
	"github.com/xjasmx/kpfm/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
