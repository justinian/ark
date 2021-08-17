package main

import (
	"encoding/json"
	"fmt"
)

type intProperty struct {
	bytes  uint8
	signed bool
	value  int64
}

func (p *intProperty) Type() PropertyType { return IntProperty }

func (p *intProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

func (p *intProperty) String() string {
	pre := "I"
	if !p.signed {
		pre = "UI"
	}

	return fmt.Sprintf("%snt%2dProperty(%d)", pre, p.bytes*8, p.value)
}

func readUIntProperty(dataSize int, vr valueReader) (Property, error) {
	return readIntPropertyBase(false, dataSize, vr)
}

func readIntProperty(dataSize int, vr valueReader) (Property, error) {
	return readIntPropertyBase(true, dataSize, vr)
}

func readIntPropertyBase(signed bool, dataSize int, vr valueReader) (Property, error) {
	value, err := vr.readIntOfSize(dataSize)
	if err != nil {
		return nil, fmt.Errorf("Reading int value: %w", err)
	}

	return &intProperty{
		bytes:  uint8(dataSize),
		signed: signed,
		value:  value,
	}, nil
}

func init() {
	addPropertyType("IntProperty", 4, readIntProperty, nil)
	addPropertyType("Int8Property", 1, readIntProperty, nil)
	addPropertyType("Int16Property", 2, readIntProperty, nil)
	addPropertyType("Int32Property", 4, readIntProperty, nil)
	addPropertyType("Int64Property", 8, readIntProperty, nil)

	addPropertyType("UIntProperty", 4, readUIntProperty, nil)
	addPropertyType("UInt8Property", 1, readUIntProperty, nil)
	addPropertyType("UInt16Property", 2, readUIntProperty, nil)
	addPropertyType("UInt32Property", 4, readUIntProperty, nil)
	addPropertyType("UInt64Property", 8, readUIntProperty, nil)
}
