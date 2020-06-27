package main

import (
	"fmt"
)

type enumProperty struct {
	enum  Name
	value Name
}

func (p *enumProperty) Type() PropertyType { return EnumProperty }

func (p *enumProperty) String() string {
	return fmt.Sprintf("EnumProperty(%s : %s)", p.enum, p.value)
}

func readEnumProperty(dataSize int, a *Archive) (Property, error) {
	enum, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum name: %w", err)
	}

	if enum.IsNone() {
		return readIntProperty(false, 1, a)
	}

	value, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum value: %w", err)
	}

	return &enumProperty{
		enum:  enum,
		value: value,
	}, nil
}
