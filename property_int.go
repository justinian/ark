package main

import (
	"fmt"
)

type intProperty struct {
	name   Name
	index  int
	bytes  uint8
	signed bool
	value  int64
}

func (p *intProperty) Type() PropertyType { return IntProperty }
func (p *intProperty) Name() Name         { return p.name }
func (p *intProperty) Index() int         { return p.index }

func (p *intProperty) String() string {
	pre := " I"
	if !p.signed {
		pre = "UI"
	}

	return fmt.Sprintf("%snt%dProperty %-28s [%d] = %d", pre, p.bytes*8, p.name, p.index, p.value)
}

func readIntProperty(name Name, signed bool, dataSize, index int, a *Archive) (Property, error) {
	value, err := a.readIntOfSize(dataSize)
	if err != nil {
		return nil, fmt.Errorf("Reading int value: %w", err)
	}

	return &intProperty{
		name:   name,
		index:  index,
		bytes:  uint8(dataSize),
		signed: signed,
		value:  value,
	}, nil
}
