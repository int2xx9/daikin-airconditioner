package daikin

import (
	"net"
	"testing"
)

func TestRoomTemperatureSignedByte(t *testing.T) {
	resp := QueryResponse{
		Address: net.UDPAddr{},
		data: map[byte][]byte{
			// 0xfb is -5 in two's complement.
			EpcRoomTemperature: {0xfb},
		},
	}

	got, err := resp.RoomTemperature()
	if err != nil {
		t.Fatalf("RoomTemperature() error = %v", err)
	}
	if got != -5 {
		t.Fatalf("RoomTemperature() = %d, want %d", got, -5)
	}
}

func TestOutdoorTemperatureSignedByte(t *testing.T) {
	resp := QueryResponse{
		Address: net.UDPAddr{},
		data: map[byte][]byte{
			// 0xff is -1 in two's complement.
			EpcOutdoorTemperature: {0xff},
		},
	}

	got, err := resp.OutdoorTemperature()
	if err != nil {
		t.Fatalf("OutdoorTemperature() error = %v", err)
	}
	if got != -1 {
		t.Fatalf("OutdoorTemperature() = %d, want %d", got, -1)
	}
}
