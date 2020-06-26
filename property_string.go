package main

import (
	"fmt"
)

type stringProperty struct {
	name  Name
	index int
	value string
}

func (p *stringProperty) Type() PropertyType { return StringProperty }
func (p *stringProperty) Name() Name         { return p.name }
func (p *stringProperty) Index() int         { return p.index }

func (p *stringProperty) String() string {
	return fmt.Sprintf("StringProperty %-28s [%d] = %s", p.name, p.index, p.value)
}

func readStringProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	value, err := a.readString()
	if err != nil {
		return nil, fmt.Errorf("Reading string value: %w", err)
	}

	return &stringProperty{
		name:  name,
		index: index,
		value: value,
	}, nil
}
