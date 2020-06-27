package main

import "fmt"

type arrayProperty struct {
	arrayType Name
	data      []byte
}

func (p *arrayProperty) Type() PropertyType { return ArrayProperty }

func (p *arrayProperty) String() string {
	return fmt.Sprintf("ArrayProperty(%s)", p.arrayType)
}

func readArrayProperty(dataSize int, a *Archive) (Property, error) {
	arrayType, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading array type: %w", err)
	}

	data := make([]byte, dataSize)
	n, err := a.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Reading property data: %w", err)
	} else if n != dataSize {
		return nil, fmt.Errorf("Property data short read: %d/%d", n, dataSize)
	}

	return &arrayProperty{
		arrayType: arrayType,
		data:      data,
	}, nil
}
