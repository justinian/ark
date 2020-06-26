package main

import (
	"fmt"
)

type nameProperty struct {
	name  Name
	index int
	value Name
}

func (p *nameProperty) Type() PropertyType { return NameProperty }
func (p *nameProperty) Name() Name         { return p.name }
func (p *nameProperty) Index() int         { return p.index }

func (p *nameProperty) String() string {
	return fmt.Sprintf("  NameProperty %-28s [%d] = %s", p.name, p.index, p.value)
}

func readNameProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	value, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading name value: %w", err)
	}

	return &nameProperty{
		name:  name,
		index: index,
		value: value,
	}, nil
}
