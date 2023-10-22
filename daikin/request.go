package daikin

import (
	"errors"
	"time"

	"github.com/int2xx9/daikin-airconditioner/echonetlite"
)

const (
	ObjectAircon     = 0x013001
	ObjectController = 0x05ff01
)

const (
	EchonetLiteTimeout = 1 * time.Second
)

type QueryRequest struct {
	daikin *Daikin
	epcs   map[byte]any
}

var (
	ErrQueryFailed = errors.New("query failed")
)

func (r QueryRequest) Query() ([]QueryResponse, error) {
	frame := r.daikin.controller.CreateFrame()
	frame.Edata = echonetlite.SpecifiedMessage{
		Seoj:       ObjectController,
		Deoj:       ObjectAircon,
		Esv:        echonetlite.ServiceTypeGet,
		Properties: []echonetlite.Property{},
	}
	for epc := range r.epcs {
		frame.Edata.Properties = append(frame.Edata.Properties, echonetlite.Property{
			Epc: epc,
			Edt: []byte{},
		})
	}

	responses, err := r.daikin.controller.QueryBuilder().SetTimeout(EchonetLiteTimeout).Query(frame)
	if err != nil {
		return []QueryResponse{}, err
	}

	retResponses := []QueryResponse{}
	var lasterror error = nil
	for _, res := range responses {
		if res.Frame.Edata.Esv != echonetlite.ServiceTypeGetRes {
			lasterror = ErrQueryFailed
			continue
		}

		retResponse := QueryResponse{
			Address: res.Addr,
			data:    map[byte][]byte{},
		}
		retResponses = append(retResponses, retResponse)
		for _, prop := range res.Frame.Edata.Properties {
			retResponse.data[prop.Epc] = prop.Edt
		}
	}

	return retResponses, lasterror
}

func (r QueryRequest) AddEpc(epc byte) QueryRequest {
	r.epcs[epc] = true
	return r
}

func (r QueryRequest) OperationStatus() QueryRequest {
	return r.AddEpc(EpcOperationStatus)
}

func (r QueryRequest) IdentificationNumber() QueryRequest {
	return r.AddEpc(EpcIdentificationNumber)
}

func (r QueryRequest) InstantaneousPowerConsumption() QueryRequest {
	return r.AddEpc(EpcInstantaneousPowerConsumption)
}

func (r QueryRequest) CumulativePowerConsumption() QueryRequest {
	return r.AddEpc(EpcCumulativePowerConsumption)
}

func (r QueryRequest) FaultStatus() QueryRequest {
	return r.AddEpc(EpcFaultStatus)
}

func (r QueryRequest) AirflowRate() QueryRequest {
	return r.AddEpc(EpcAirflowRate)
}

func (r QueryRequest) OperationMode() QueryRequest {
	return r.AddEpc(EpcOperationMode)
}

func (r QueryRequest) TemperatureSetting() QueryRequest {
	return r.AddEpc(EpcTemperatureSetting)
}

func (r QueryRequest) HumiditySetting() QueryRequest {
	return r.AddEpc(EpcHumiditySetting)
}

func (r QueryRequest) RoomTemperature() QueryRequest {
	return r.AddEpc(EpcRoomTemperature)
}

func (r QueryRequest) RoomHumidity() QueryRequest {
	return r.AddEpc(EpcRoomHumidity)
}

func (r QueryRequest) OutdoorTemperature() QueryRequest {
	return r.AddEpc(EpcOutdoorTemperature)
}
