package echonetlite_test

import (
	"reflect"
	"testing"

	"github.com/int2xx9/daikin-airconditioner/echonetlite"
)

func TestFrame(t *testing.T) {
	t.Run("Serialize", func(t *testing.T) {
		actual, _ := echonetlite.Frame{
			Ehd1: 0x10,
			Ehd2: 0x81,
			Tid:  0x1234,
			Edata: echonetlite.SpecifiedMessage{
				Seoj: 0x123456,
				Deoj: 0x789abc,
				Esv:  0x60,
			},
		}.Serialize()
		expect := []byte{
			0x10,
			0x81,
			0x12, 0x34,
			0x12, 0x34, 0x56,
			0x78, 0x9a, 0xbc,
			0x60,
			0x00,
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Serialize failure")
		}
	})
}

func TestDeserializeFrame(t *testing.T) {
	actual, err := echonetlite.DeserializeFrame([]byte{
		0x10,
		0x81,
		0x12, 0x34,
		0x12, 0x34, 0x56,
		0x78, 0x9a, 0xbc,
		0x60,
		0x00,
	})
	expect := echonetlite.Frame{
		Ehd1: 0x10,
		Ehd2: 0x81,
		Tid:  0x1234,
		Edata: echonetlite.SpecifiedMessage{
			Seoj:       0x123456,
			Deoj:       0x789abc,
			Esv:        0x60,
			Properties: []echonetlite.Property{},
		},
	}
	if err != nil {
		t.Errorf("DeserializeFrame failure")
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("DeserializeFrame failure")
	}
}

func TestSpecifiedMessage(t *testing.T) {
	t.Run("Serialize", func(t *testing.T) {
		actual, _ := echonetlite.SpecifiedMessage{
			Seoj: 0x123456,
			Deoj: 0x789abc,
			Esv:  0x60,
			Properties: []echonetlite.Property{
				{
					Epc: 0x10,
					Edt: []byte{0x01, 0x02},
				},
			},
		}.Serialize()
		expect := []byte{
			0x12, 0x34, 0x56,
			0x78, 0x9a, 0xbc,
			0x60,
			0x01,
			0x10,
			0x02,
			0x01, 0x02,
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Serialize failure")
		}
	})
}

func TestDeserializeSpecifiedMessage(t *testing.T) {
	actual, err := echonetlite.DeserializeSpecifiedMessage([]byte{
		0x12, 0x34, 0x56,
		0x78, 0x9a, 0xbc,
		0x60,
		0x01,
		0x01,
		0x02,
		0x02, 0x03,
	})
	expect := echonetlite.SpecifiedMessage{
		Seoj: 0x123456,
		Deoj: 0x789abc,
		Esv:  0x60,
		Properties: []echonetlite.Property{
			{
				Epc: 0x01,
				Edt: []byte{0x02, 0x03},
			},
		},
	}
	if err != nil {
		t.Errorf("DeserializeSepecifiedMessage failure")
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("DeserializeSepecifiedMessage failure")
	}
}

func TestDeserializeProperties(t *testing.T) {
	actual, err := echonetlite.DeserializeProperties([]byte{
		0x01,
		0x02,
		0x02, 0x03,
		0x04,
		0x03,
		0x05, 0x06, 0x07,
	})
	expect := []echonetlite.Property{
		{
			Epc: 0x01,
			Edt: []byte{0x02, 0x03},
		},
		{
			Epc: 0x04,
			Edt: []byte{0x05, 0x06, 0x07},
		},
	}
	if err != nil {
		t.Errorf("DeserializeProperties failure")
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("DeserializeProperties failure")
	}
}

func TestProperty(t *testing.T) {
	t.Run("Serialize", func(t *testing.T) {
		actual, _ := echonetlite.Property{
			Epc: 0x10,
			Edt: []byte{0x01, 0x02},
		}.Serialize()
		expect := []byte{
			0x10,
			0x02,
			0x01, 0x02,
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Serialize failure")
		}
	})
}

func TestDeserializeProperty(t *testing.T) {
	actual, err := echonetlite.DeserializeProperty([]byte{
		0x01,
		0x02,
		0x02, 0x03,
	})
	expect := echonetlite.Property{
		Epc: 0x01,
		Edt: []byte{0x02, 0x03},
	}
	if err != nil {
		t.Errorf("DeserializeProperty failure")
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("DeserializeProperty failure")
	}
}
