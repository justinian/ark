package main

import (
	"fmt"
)

type enumProperty struct {
	name      Name
	index     int
	enum      Name
	value     int64
	valueName Name
}

func (p *enumProperty) Type() PropertyType { return EnumProperty }
func (p *enumProperty) Name() Name         { return p.name }
func (p *enumProperty) Index() int         { return p.index }

func (p *enumProperty) String() string {
	return fmt.Sprintf("  EnumProperty %-14s %-14s [%d] = %s(%d)", p.enum, p.name, p.index, p.valueName, p.value)
}

func readEnumProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	enum, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum name: %w", err)
	}

	if enum.IsNone() {
		return readIntProperty(name, false, 1, index, a)
	}

	valueName, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum value name: %w", err)
	}

	value, err := a.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Reading int value: %w", err)
	}

	return &enumProperty{
		name:      name,
		index:     index,
		enum:      enum,
		value:     value,
		valueName: valueName,
	}, nil
}
