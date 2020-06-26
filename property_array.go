package main

import "fmt"

type arrayProperty struct {
	name      Name
	index     int
	arrayType Name
	data      []byte
}

func (p *arrayProperty) Type() PropertyType { return ArrayProperty }
func (p *arrayProperty) Name() Name         { return p.name }
func (p *arrayProperty) Index() int         { return p.index }

func (p *arrayProperty) String() string {
	return fmt.Sprintf(" ArrayProperty %s [%d] = %s", p.name, p.index, p.arrayType)
}

func readArrayProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
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
		name:      name,
		index:     index,
		arrayType: arrayType,
		data:      data,
	}, nil
}
