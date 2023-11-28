package service

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/adapter"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/repository"
	"github.com/int2xx9/daikin-airconditioner/daikin"
	"github.com/int2xx9/daikin-airconditioner/echonetlite"
)

type RecordService struct {
	config           config.Configuration
	logger           *slog.Logger
	deviceRepository *repository.DeviceRepository
	recordRepository *repository.RecordRepository
	daikinAdapter    *adapter.DaikinAdapter
}

func NewRecordService(config config.Configuration, logger *slog.Logger, deviceRepository *repository.DeviceRepository, recordRepository *repository.RecordRepository, daikinAdapter *adapter.DaikinAdapter) RecordService {
	return RecordService{
		config:           config,
		logger:           logger,
		deviceRepository: deviceRepository,
		recordRepository: recordRepository,
		daikinAdapter:    daikinAdapter,
	}
}

func (s RecordService) Record() error {
	controller := echonetlite.NewController()
	controller.Start()
	defer controller.Close()
	daikin := daikin.NewDaikin(&controller)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	ticker := time.NewTicker(s.config.Scrape.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.recordInternal(daikin); err != nil {
				s.logger.Error("RecordInternal failed", "error", err)
			}
		case sig := <-sigCh:
			fmt.Printf("Trap a signal: %s\n", sig)
			return nil
		}
	}
}

func (s RecordService) recordInternal(daikinClient daikin.Daikin) error {
	devices, err := s.deviceRepository.GetDevices()
	if err != nil {
		return err
	}

	now := time.Now()
	statuses, err := s.daikinAdapter.GetStatuses()
	if err != nil {
		return err
	}

	for _, status := range statuses {
		idstr := hex.EncodeToString(status.IdentificationNumber)

		isKnownDevice := false
		for _, device_id := range devices {
			if device_id == idstr {
				isKnownDevice = true
			}
		}
		if !isKnownDevice {
			s.logger.Debug("Unknown device_id. Ignore.", "device_id", idstr)
			continue
		}

		data := map[string]any{}

		if status.OperationStatus != nil {
			data["operation_status"] = *status.OperationStatus
		}

		if status.InstantaneousPowerConsumption != nil {
			data["instantaneous_power_consumption"] = *status.InstantaneousPowerConsumption
		}

		if status.CumulativePowerConsumption != nil {
			data["cumulative_power_consumption"] = *status.CumulativePowerConsumption
		}

		if status.FaultStatus != nil {
			data["fault_status"] = *status.FaultStatus
		}

		if status.AirflowRateAuto != nil {
			data["airflowrate_auto"] = *status.AirflowRateAuto
		}

		if status.AirflowRate != nil {
			data["airflowrate_setting"] = *status.AirflowRate
		}

		if status.OperationMode != nil {
			data["operation_mode"] = *status.OperationMode
		}

		if status.TemperatureSetting != nil {
			data["temperature_setting"] = *status.TemperatureSetting
		}

		if status.HumiditySetting != nil {
			data["humidity_setting"] = *status.HumiditySetting
		}

		if status.RoomTemperature != nil {
			data["room_temperature"] = *status.RoomTemperature
		}

		if status.RoomHumidity != nil {
			data["room_humidity"] = *status.RoomHumidity
		}

		if status.OutdoorTemperature != nil {
			data["outdoor_temperature"] = *status.OutdoorTemperature
		}

		affected, err := s.recordRepository.Add(idstr, now, data)
		if err != nil {
			s.logger.Error("Insert error", "error", err)
		} else {
			s.logger.Debug("Insert success", "id", idstr, "affected_rows", affected)
		}
	}

	return nil
}
