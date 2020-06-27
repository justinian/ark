package main

import (
	"fmt"
)

type nameProperty struct {
	value Name
}

func (p *nameProperty) Type() PropertyType { return NameProperty }

func (p *nameProperty) String() string {
	return fmt.Sprintf("NameProperty(%s)", p.value)
}

func readNameProperty(dataSize int, a *Archive) (Property, error) {
	value, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading name value: %w", err)
	}

	return &nameProperty{
		value: value,
	}, nil
}
