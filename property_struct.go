package main

import "fmt"

type structProperty struct {
	structType Name
	data       []byte
}

func (p *structProperty) Type() PropertyType { return StructProperty }

func (p *structProperty) String() string {
	return fmt.Sprintf("StructProperty(%s)", p.structType)
}

func readStructProperty(dataSize int, a *Archive) (Property, error) {
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
		structType: structType,
		data:       data,
	}, nil
}
