package echonetlite

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrUnsupportedFrame  = errors.New("unsupported frame")
	ErrElementsMismatch  = errors.New("the number of elements is mismatched")
	ErrTooManyProperties = errors.New("too many properties")
	ErrWrongLength       = errors.New("wrong length")
)

type ServiceType byte

const (
	ServiceTypeSetI   ServiceType = 0x60
	ServiceTypeSetC   ServiceType = 0x61
	ServiceTypeGet    ServiceType = 0x62
	ServiceTypeInfReq ServiceType = 0x63
	ServiceTypeSetGet ServiceType = 0x6e

	ServiceTypeSetReq    ServiceType = 0x71
	ServiceTypeGetRes    ServiceType = 0x72
	ServiceTypeInf       ServiceType = 0x73
	ServiceTypeInfC      ServiceType = 0x74
	ServiceTypeInfCRes   ServiceType = 0x7a
	ServiceTypeSetGetRes ServiceType = 0x7e

	ServiceTypeSetISna   ServiceType = 0x50
	ServiceTypeSetCSna   ServiceType = 0x51
	ServiceTypeGetSna    ServiceType = 0x52
	ServiceTypeInfSna    ServiceType = 0x53
	ServiceTypeSetGetSna ServiceType = 0x5e
)

var (
	ServiceTypeNames = map[ServiceType]string{
		ServiceTypeSetI:   "SetI",
		ServiceTypeSetC:   "SetC",
		ServiceTypeGet:    "Get",
		ServiceTypeInfReq: "InfReq",
		ServiceTypeSetGet: "SetGet",

		ServiceTypeSetReq:    "SetReq",
		ServiceTypeGetRes:    "GetRes",
		ServiceTypeInf:       "Inf",
		ServiceTypeInfC:      "InfC",
		ServiceTypeInfCRes:   "InfCRes",
		ServiceTypeSetGetRes: "SetGetRes",

		ServiceTypeSetISna:   "SetISna",
		ServiceTypeSetCSna:   "SetCSna",
		ServiceTypeGetSna:    "GetSna",
		ServiceTypeInfSna:    "InfSna",
		ServiceTypeSetGetSna: "SetGetSna",
	}
)

func (t ServiceType) String() string {
	if value, ok := ServiceTypeNames[t]; ok {
		return value
	}
	return "Unknown"
}

type Serializable interface {
	Serialize() []byte
}

type Frame struct {
	Ehd1  byte
	Ehd2  byte
	Tid   uint16
	Edata SpecifiedMessage
}

func (f Frame) Serialize() ([]byte, error) {
	bytes, err := f.Edata.Serialize()
	if err != nil {
		return nil, err
	}

	data := []byte{f.Ehd1, f.Ehd2, byte(f.Tid >> 8), byte(f.Tid & 0xff)}
	data = append(data, bytes...)
	return data, nil
}

func DeserializeFrame(data []byte) (Frame, error) {
	f := Frame{}
	f.Ehd1 = data[0]
	f.Ehd2 = data[1]
	f.Tid = (uint16(data[2]) << 8) | uint16(data[3])

	if f.Ehd1 != 0x10 || f.Ehd2 != 0x81 {
		return Frame{}, ErrUnsupportedFrame
	}

	edata, err := DeserializeSpecifiedMessage(data[4:])
	if err != nil {
		return Frame{}, err
	}
	f.Edata = edata

	return f, nil
}

type SpecifiedMessage struct {
	Seoj       uint32
	Deoj       uint32
	Esv        ServiceType
	Properties []Property
}

func DeserializeSpecifiedMessage(data []byte) (SpecifiedMessage, error) {
	m := SpecifiedMessage{}

	if err := binary.Read(bytes.NewReader(append([]byte{0}, data[0:3]...)), binary.BigEndian, &m.Seoj); err != nil {
		return SpecifiedMessage{}, err
	}

	if err := binary.Read(bytes.NewReader(append([]byte{0}, data[3:6]...)), binary.BigEndian, &m.Deoj); err != nil {
		return SpecifiedMessage{}, err
	}

	m.Esv = ServiceType(data[6])
	opc := int(data[7])
	properties, err := DeserializeProperties(data[8:])
	if err != nil {
		return SpecifiedMessage{}, err
	}
	if opc != len(properties) {
		return SpecifiedMessage{}, ErrElementsMismatch
	}

	m.Properties = properties
	return m, nil
}

func DeserializeProperties(data []byte) ([]Property, error) {
	p := []Property{}

	offset := 0
	for offset < len(data) {
		pdc := int(data[offset+1])
		prop, err := DeserializeProperty(data[offset : offset+pdc+2])
		if err != nil {
			return []Property{}, err
		}
		p = append(p, prop)
		offset += pdc + 2
	}

	return p, nil
}

func (m SpecifiedMessage) Serialize() ([]byte, error) {
	if len(m.Properties) > 255 {
		return nil, ErrTooManyProperties
	}

	data := []byte{
		// Seoj
		byte((m.Seoj >> 16) & 0xff),
		byte((m.Seoj >> 8) & 0xff),
		byte(m.Seoj & 0xff),
		// Deoj
		byte((m.Deoj >> 16) & 0xff),
		byte((m.Deoj >> 8) & 0xff),
		byte(m.Deoj & 0xff),
		// Esv
		byte(m.Esv),
		// Opc
		byte(len(m.Properties)),
	}
	for _, property := range m.Properties {
		propertyBytes, err := property.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, propertyBytes...)
	}
	return data, nil
}

type Property struct {
	Epc byte
	Edt []byte
}

func DeserializeProperty(data []byte) (Property, error) {
	p := Property{}

	p.Epc = data[0]
	pdc := int(data[1])
	if pdc+2 != len(data) {
		return Property{}, ErrWrongLength
	}
	p.Edt = make([]byte, pdc)
	copy(p.Edt, data[2:2+pdc])

	return p, nil
}

func (p Property) Serialize() ([]byte, error) {
	if len(p.Edt) > 255 {
		return nil, ErrWrongLength
	}
	data := []byte{p.Epc, byte(len(p.Edt))}
	data = append(data, p.Edt...)
	return data, nil
}
