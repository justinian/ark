package main

import (
	"fmt"
)

type enumProperty struct {
	Enum  Name `json:"enum"`
	Value Name `json:"value"`
}

func (p *enumProperty) Type() PropertyType { return EnumProperty }

func (p *enumProperty) String() string {
	return fmt.Sprintf("EnumProperty(%s : %s)", p.Enum, p.Value)
}

func readEnumProperty(dataSize int, vr valueReader) (Property, error) {
	enum, err := vr.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum name: %w", err)
	}

	if enum.IsNone() {
		return readUIntProperty(1, vr)
	}

	value, err := vr.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum value: %w", err)
	}

	return &enumProperty{
		Enum:  enum,
		Value: value,
	}, nil
}

func init() {
	addPropertyType("ByteProperty", 1, readEnumProperty, nil)
}
