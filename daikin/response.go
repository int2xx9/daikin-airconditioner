package daikin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

var (
	ErrNoResponsesForEpc = errors.New("no responses for epc")
	ErrUnexpectedValue   = errors.New("unexpected value")
	ErrUnsupportedValue  = errors.New("unsupported value")
)

type OperationMode byte

const (
	OperationModeAuto             OperationMode = 0x41
	OperationModeCooling          OperationMode = 0x42
	OperationModeHeating          OperationMode = 0x43
	OperationModeDehumidification OperationMode = 0x44
	OperationModeVentilating      OperationMode = 0x45
	OperationModeOther            OperationMode = 0x40
)

type QueryResponse struct {
	Address net.UDPAddr
	data    map[byte][]byte
}

func (q QueryResponse) OperationStatus() (bool, error) {
	data, ok := q.data[EpcOperationStatus]
	if !ok {
		return false, ErrNoResponsesForEpc
	}

	switch data[0] {
	case 0x30:
		return true, nil
	case 0x31:
		return false, nil
	default:
		return false, ErrUnexpectedValue
	}
}

func (q QueryResponse) IdentificationNumber() ([]byte, error) {
	data, ok := q.data[EpcIdentificationNumber]
	if !ok {
		return nil, ErrNoResponsesForEpc
	}

	if data[0] != 0xfe {
		return nil, ErrUnsupportedValue
	}

	ret := make([]byte, 16)
	copy(ret, data[1:])
	return ret, nil
}

func (q QueryResponse) InstantaneousPowerConsumption() (int, error) {
	data, ok := q.data[EpcInstantaneousPowerConsumption]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	return int(data[0])<<8 | int(data[1]), nil
}

func (q QueryResponse) CumulativePowerConsumption() (int, error) {
	data, ok := q.data[EpcCumulativePowerConsumption]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	var value uint32
	if err := binary.Read(bytes.NewReader(data[0:4]), binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return int(value), nil
}

func (q QueryResponse) FaultStatus() (bool, error) {
	data, ok := q.data[EpcFaultStatus]
	if !ok {
		return false, ErrNoResponsesForEpc
	}

	switch data[0] {
	case 0x41:
		return true, nil
	case 0x42:
		return false, nil
	default:
		return false, ErrUnexpectedValue
	}
}

func (q QueryResponse) AirflowRate() (int, bool, error) {
	data, ok := q.data[EpcAirflowRate]
	if !ok {
		return 0, false, ErrNoResponsesForEpc
	}

	if data[0] == 0x41 {
		return 0, true, nil
	}

	return int(data[0] - 0x30), false, nil
}

func (q QueryResponse) OperationMode() (OperationMode, error) {
	data, ok := q.data[EpcOperationMode]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	if data[0] < 0x40 || data[0] > 0x45 {
		return 0, ErrUnsupportedValue
	}

	return OperationMode(data[0]), nil
}

func (q QueryResponse) TemperatureSetting() (int, error) {
	data, ok := q.data[EpcTemperatureSetting]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	return int(data[0]), nil
}

func (q QueryResponse) HumiditySetting() (int, error) {
	data, ok := q.data[EpcHumiditySetting]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	return int(data[0]), nil
}

func (q QueryResponse) DehumidifierStatus() (bool, error) {
	data, ok := q.data[EpcDehumidifyingSetting]
	if !ok {
		return false, ErrNoResponsesForEpc
	}

	switch data[0] {
	case 0x41:
		return true, nil
	case 0x42:
		return false, nil
	}

	return false, ErrUnexpectedValue
}

func (q QueryResponse) RoomTemperature() (int, error) {
	data, ok := q.data[EpcRoomTemperature]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	return int(data[0]), nil
}

func (q QueryResponse) RoomHumidity() (int, error) {
	data, ok := q.data[EpcRoomHumidity]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	return int(data[0]), nil
}

func (q QueryResponse) OutdoorTemperature() (int, error) {
	data, ok := q.data[EpcOutdoorTemperature]
	if !ok {
		return 0, ErrNoResponsesForEpc
	}

	return int(data[0]), nil
}
