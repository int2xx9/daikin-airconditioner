package echonetlite

import "errors"

var (
	ErrUnexpectedEpc           = errors.New("unexpected epc")
	ErrPropertyCountMismatched = errors.New("the number of properties is mismatched")
)

func GetPropertyMap(p Property) ([]byte, error) {
	if p.Epc < 0x9b || p.Epc > 0x9f {
		return nil, ErrUnexpectedEpc
	}

	propCount := int(p.Edt[0])
	if propCount <= 16 {
		list := make([]byte, propCount)
		copy(list, p.Edt[1:])
		return list, nil
	}

	list := []byte{}
	for i := byte(0); i < 16; i++ {
		b := p.Edt[1+i]
		if b&0b10000000 != 0 {
			list = append(list, 0xf0|i)
		}
		if b&0b01000000 != 0 {
			list = append(list, 0xe0|i)
		}
		if b&0b00100000 != 0 {
			list = append(list, 0xd0|i)
		}
		if b&0b00010000 != 0 {
			list = append(list, 0xc0|i)
		}
		if b&0b00001000 != 0 {
			list = append(list, 0xb0|i)
		}
		if b&0b00000100 != 0 {
			list = append(list, 0xa0|i)
		}
		if b&0b00000010 != 0 {
			list = append(list, 0x90|i)
		}
		if b&0b00000001 != 0 {
			list = append(list, 0x80|i)
		}
	}

	if len(list) != propCount {
		return nil, ErrPropertyCountMismatched
	}
	return list, nil
}
