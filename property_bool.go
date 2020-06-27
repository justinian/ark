package main

import (
	"fmt"
)

type boolProperty struct {
	value bool
}

func (p *boolProperty) Type() PropertyType { return BoolProperty }

func (p *boolProperty) String() string {
	return fmt.Sprintf("BoolProperty(%v)", p.value)
}

func readBoolProperty(dataSize int, a *Archive) (Property, error) {
	value, err := a.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Reading bool value: %w", err)
	}

	return &boolProperty{
		value: value != 0,
	}, nil
}
