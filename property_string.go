package main

import (
	"fmt"
)

type stringProperty struct {
	value string
}

func (p *stringProperty) Type() PropertyType { return StringProperty }

func (p *stringProperty) String() string {
	return fmt.Sprintf("StringProperty(%s)", p.value)
}

func readStringProperty(dataSize int, a *Archive) (Property, error) {
	value, err := a.readString()
	if err != nil {
		return nil, fmt.Errorf("Reading string value: %w", err)
	}

	return &stringProperty{
		value: value,
	}, nil
}
