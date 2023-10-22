package daikin

import (
	"github.com/int2xx9/daikin-airconditioner/echonetlite"
)

const (
	EpcOperationStatus               byte = 0x80
	EpcIdentificationNumber          byte = 0x83
	EpcInstantaneousPowerConsumption byte = 0x84
	EpcCumulativePowerConsumption    byte = 0x85
	EpcFaultStatus                   byte = 0x88
	EpcAirflowRate                   byte = 0xa0
	EpcOperationMode                 byte = 0xb0
	EpcTemperatureSetting            byte = 0xb3
	EpcHumiditySetting               byte = 0xb4
	EpcRoomTemperature               byte = 0xbb
	EpcRoomHumidity                  byte = 0xba
	EpcOutdoorTemperature            byte = 0xbe
)

type Daikin struct {
	controller *echonetlite.Controller
}

func NewDaikin(c *echonetlite.Controller) Daikin {
	return Daikin{
		controller: c,
	}
}

func (d *Daikin) Request() QueryRequest {
	return QueryRequest{
		daikin: d,
		epcs:   map[byte]any{},
	}
}
