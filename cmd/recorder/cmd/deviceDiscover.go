package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/adapter"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/repository"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/service"
	"github.com/int2xx9/daikin-airconditioner/daikin"
	"github.com/int2xx9/daikin-airconditioner/echonetlite"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var deviceDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover devices",
	Run: func(cmd *cobra.Command, args []string) {
		controller := echonetlite.NewController()
		controller.Start()
		defer controller.Close()
		daikin := daikin.NewDaikin(&controller)
		daikinAdapter := adapter.NewDaikinAdapter(daikin, slog.Default())

		db, err := sql.Open("postgres", config.Config.Database.ConnectionString())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		service := service.NewDeviceService(config.Config, daikinAdapter, service.NewDeviceQueryService(db), repository.NewDeviceRepository(db), repository.NewRecordRepository(db))
		if err := service.Discover(); err != nil {
			panic(err)
		}
	},
}

func init() {
	deviceCmd.AddCommand(deviceDiscoverCmd)
}
