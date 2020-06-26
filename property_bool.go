package main

import (
	"fmt"
)

type boolProperty struct {
	name  Name
	index int
	value bool
}

func (p *boolProperty) Type() PropertyType { return BoolProperty }
func (p *boolProperty) Name() Name         { return p.name }
func (p *boolProperty) Index() int         { return p.index }

func (p *boolProperty) String() string {
	return fmt.Sprintf("  BoolProperty %-28s [%d] = %v", p.name, p.index, p.value)
}

func readBoolProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	value, err := a.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Reading bool value: %w", err)
	}

	return &boolProperty{
		name:  name,
		index: index,
		value: value != 0,
	}, nil
}
