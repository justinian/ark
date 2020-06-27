package main

import (
	"fmt"
)

type intProperty struct {
	bytes  uint8
	signed bool
	value  int64
}

func (p *intProperty) Type() PropertyType { return IntProperty }

func (p *intProperty) String() string {
	pre := "I"
	if !p.signed {
		pre = "UI"
	}

	return fmt.Sprintf("%snt%2dProperty(%d)", pre, p.bytes*8, p.value)
}

func readIntProperty(signed bool, dataSize int, a *Archive) (Property, error) {
	value, err := a.readIntOfSize(dataSize)
	if err != nil {
		return nil, fmt.Errorf("Reading int value: %w", err)
	}

	return &intProperty{
		bytes:  uint8(dataSize),
		signed: signed,
		value:  value,
	}, nil
}
