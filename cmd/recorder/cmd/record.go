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

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Start recornding",
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

		//slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
		service := service.NewRecordService(config.Config, slog.Default(), repository.NewDeviceRepository(db), repository.NewRecordRepository(db), daikinAdapter)
		if err := service.Record(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
