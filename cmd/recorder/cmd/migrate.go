package cmd

import (
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/service"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate a database",
	Run: func(cmd *cobra.Command, args []string) {
		service := service.NewPersistentService(config.Config)
		if err := service.Migrate(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
