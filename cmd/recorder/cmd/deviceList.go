package cmd

import (
	"database/sql"

	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/repository"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/service"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show a list of devices",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("postgres", config.Config.Database.ConnectionString())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		service := service.NewDeviceService(config.Config, nil, service.NewDeviceQueryService(db), repository.NewDeviceRepository(db), repository.NewRecordRepository(db))
		if err := service.List(); err != nil {
			panic(err)
		}
	},
}

func init() {
	deviceCmd.AddCommand(deviceListCmd)
}
