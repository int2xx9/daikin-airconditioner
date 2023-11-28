package cmd

import (
	"database/sql"

	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/repository"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/service"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var deviceRenameCmd = &cobra.Command{
	Use:   "rename <id> <name>",
	Short: "Rename a device",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}
		id := args[0]
		name := args[1]

		db, err := sql.Open("postgres", config.Config.Database.ConnectionString())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		service := service.NewDeviceService(config.Config, nil, service.NewDeviceQueryService(db), repository.NewDeviceRepository(db), repository.NewRecordRepository(db))
		if err := service.Rename(id, name); err != nil {
			panic(err)
		}
	},
}

func init() {
	deviceCmd.AddCommand(deviceRenameCmd)
}