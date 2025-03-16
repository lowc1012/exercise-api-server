package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Code exercise",
	Long:  "HTTP API server exercise",
	RunE:  appRun,
}

func appRun(cmd *cobra.Command, args []string) error {
	// TODO
	return nil
}
