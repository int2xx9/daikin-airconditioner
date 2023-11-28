package cmd

import (
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device subcommand",
}

func init() {
	rootCmd.AddCommand(deviceCmd)
}
