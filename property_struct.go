package main

import "fmt"

type structProperty struct {
	name       Name
	index      int
	structType Name
	data       []byte
}

func (p *structProperty) Type() PropertyType { return StructProperty }
func (p *structProperty) Name() Name         { return p.name }
func (p *structProperty) Index() int         { return p.index }

func (p *structProperty) String() string {
	return fmt.Sprintf("StructProperty %s [%d] = %s", p.name, p.index, p.structType)
}

func readStructProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	structType, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading struct type: %w", err)
	}

	data := make([]byte, dataSize)
	n, err := a.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Reading property data: %w", err)
	} else if n != dataSize {
		return nil, fmt.Errorf("Property data short read: %d/%d", n, dataSize)
	}

	return &structProperty{
		name:       name,
		index:      index,
		structType: structType,
		data:       data,
	}, nil
}
