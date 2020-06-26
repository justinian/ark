package main

import (
	"encoding/binary"
	"fmt"
)

type floatProperty struct {
	name  Name
	index int
	value float64
}

func (p *floatProperty) Type() PropertyType { return FloatProperty }
func (p *floatProperty) Name() Name         { return p.name }
func (p *floatProperty) Index() int         { return p.index }

func (p *floatProperty) String() string {
	return fmt.Sprintf(" FloatProperty %-28s [%d] = %f", p.name, p.index, p.value)
}

func readFloatProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	var err error
	var value float64

	switch dataSize {
	case 4:
		var tmp float32
		err = binary.Read(a, binary.LittleEndian, &tmp)
		value = float64(tmp)
	case 8:
		err = binary.Read(a, binary.LittleEndian, &value)
	}

	if err != nil {
		return nil, fmt.Errorf("Reading float value: %w", err)
	}

	return &floatProperty{
		name:  name,
		index: index,
		value: value,
	}, nil
}
