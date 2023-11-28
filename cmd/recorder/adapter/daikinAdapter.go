package adapter

import (
	"encoding/hex"
	"log/slog"
	"net"

	"github.com/int2xx9/daikin-airconditioner/daikin"
)

type DaikinAdapter struct {
	daikin daikin.Daikin
	logger *slog.Logger
}

type DaikinDevice struct {
	Address net.UDPAddr
	ID      string
}

type DaikinStatus struct {
	IdentificationNumber          []byte
	OperationStatus               *bool
	InstantaneousPowerConsumption *int
	CumulativePowerConsumption    *int
	FaultStatus                   *bool
	AirflowRateAuto               *bool
	AirflowRate                   *int
	OperationMode                 *OperationMode
	TemperatureSetting            *int
	HumiditySetting               *int
	RoomTemperature               *int
	RoomHumidity                  *int
	OutdoorTemperature            *int
}

type OperationMode string

const (
	OperationModeAuto             OperationMode = "auto"
	OperationModeCooling          OperationMode = "cooling"
	OperationModeHeating          OperationMode = "heating"
	OperationModeDehumidification OperationMode = "dehumidification"
	OperationModeVentilating      OperationMode = "ventilating"
	OperationModeOther            OperationMode = "other"
)

func NewDaikinAdapter(daikin daikin.Daikin, logger *slog.Logger) *DaikinAdapter {
	return &DaikinAdapter{
		daikin: daikin,
		logger: logger,
	}
}

func (d DaikinAdapter) Discover() ([]DaikinDevice, error) {
	resps, err := d.daikin.Request().IdentificationNumber().Query()
	if err != nil {
		return nil, err
	}

	devices := []DaikinDevice{}
	for _, resp := range resps {
		id, err := resp.IdentificationNumber()
		if err != nil {
			return nil, err
		}
		devices = append(devices, DaikinDevice{
			ID:      hex.EncodeToString(id),
			Address: resp.Address,
		})
	}

	return devices, nil
}

func (d DaikinAdapter) GetStatuses() ([]DaikinStatus, error) {
	resps, err := d.daikin.Request().
		IdentificationNumber().
		OperationStatus().
		InstantaneousPowerConsumption().
		CumulativePowerConsumption().
		FaultStatus().
		AirflowRate().
		OperationMode().
		TemperatureSetting().
		HumiditySetting().
		RoomTemperature().
		RoomHumidity().
		OutdoorTemperature().
		Query()
	if err != nil {
		return nil, err
	}

	ret := []DaikinStatus{}
	for _, resp := range resps {
		status := DaikinStatus{}

		id, err := resp.IdentificationNumber()
		if err != nil {
			d.logger.Error("Query failed: IdentificationNumber", "error", err)
			continue
		}
		status.IdentificationNumber = id

		if value, err := resp.OperationStatus(); err != nil {
			d.logger.Error("Query failed: OperationStatus", "error", err)
		} else {
			status.OperationStatus = &value
		}

		if value, err := resp.InstantaneousPowerConsumption(); err != nil {
			d.logger.Error("Query failed: InstantaneousPowerConsumption", "error", err)
		} else {
			status.InstantaneousPowerConsumption = &value
		}

		if value, err := resp.CumulativePowerConsumption(); err != nil {
			d.logger.Error("Query failed: CumulativePowerConsumption", "error", err)
		} else {
			status.CumulativePowerConsumption = &value
		}

		if value, err := resp.FaultStatus(); err != nil {
			d.logger.Error("Query failed: FaultStatus", "error", err)
		} else {
			status.FaultStatus = &value
		}

		if value, auto, err := resp.AirflowRate(); err != nil {
			d.logger.Error("Query failed: AirflowRate", "error", err)
		} else {
			status.AirflowRateAuto = &auto
			status.AirflowRate = &value
		}

		if value, err := resp.OperationMode(); err != nil {
			d.logger.Error("Query failed: OperationMode", "error", err)
		} else {
			var mode OperationMode
			switch value {
			case daikin.OperationModeAuto:
				mode = OperationModeAuto
			case daikin.OperationModeCooling:
				mode = OperationModeCooling
			case daikin.OperationModeHeating:
				mode = OperationModeHeating
			case daikin.OperationModeDehumidification:
				mode = OperationModeDehumidification
			case daikin.OperationModeVentilating:
				mode = OperationModeVentilating
			case daikin.OperationModeOther:
				mode = OperationModeOther
			}
			status.OperationMode = &mode
		}

		if value, err := resp.TemperatureSetting(); err != nil {
			d.logger.Error("Query failed: TemperatureSetting", "error", err)
		} else {
			status.TemperatureSetting = &value
		}

		if value, err := resp.HumiditySetting(); err != nil {
			d.logger.Error("Query failed: HumiditySetting", "error", err)
		} else {
			status.HumiditySetting = &value
		}

		if value, err := resp.RoomTemperature(); err != nil {
			d.logger.Error("Query failed: RoomTemperature", "error", err)
		} else {
			status.RoomTemperature = &value
		}

		if value, err := resp.RoomHumidity(); err != nil {
			d.logger.Error("Query failed: RoomHumidity", "error", err)
		} else {
			status.RoomHumidity = &value
		}

		if value, err := resp.OutdoorTemperature(); err != nil {
			d.logger.Error("Query failed: OutdoorTemperature", "error", err)
		} else {
			status.OutdoorTemperature = &value
		}

		ret = append(ret, status)
	}

	return ret, nil
}
