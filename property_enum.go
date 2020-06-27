package main

import (
	"fmt"
)

type enumProperty struct {
	name  Name
	index int
	enum  Name
	value Name
}

func (p *enumProperty) Type() PropertyType { return EnumProperty }
func (p *enumProperty) Name() Name         { return p.name }
func (p *enumProperty) Index() int         { return p.index }

func (p *enumProperty) String() string {
	return fmt.Sprintf("  EnumProperty %-14s %-14s [%d] = %s", p.enum, p.name, p.index, p.value)
}

func readEnumProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	enum, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum name: %w", err)
	}

	if enum.IsNone() {
		return readIntProperty(name, false, 1, index, a)
	}

	value, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading enum value: %w", err)
	}

	return &enumProperty{
		name:  name,
		index: index,
		enum:  enum,
		value: value,
	}, nil
}
