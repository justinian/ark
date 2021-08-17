package main

import "fmt"

type colorProperty struct {
	B uint8
	G uint8
	R uint8
	A uint8
}

func (p *colorProperty) Type() PropertyType { return StructColorProperty }

func (p *colorProperty) String() string {
	return fmt.Sprintf("StructColorProperty()")
}

func readColorStruct(dataSize int, vr valueReader) (Property, error) {
	b, err := vr.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Error reading color b: %w", err)
	}

	g, err := vr.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Error reading color g: %w", err)
	}

	r, err := vr.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Error reading color r: %w", err)
	}

	a, err := vr.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Error reading color a: %w", err)
	}

	return &colorProperty{uint8(b), uint8(g), uint8(r), uint8(a)}, nil
}

type linearColorProperty struct {
	R float32
	G float32
	B float32
	A float32
}

func (p *linearColorProperty) Type() PropertyType { return StructLinearColorProperty }

func (p *linearColorProperty) String() string {
	return fmt.Sprintf("StructLinearColorProperty()")
}

func readLinearColorStruct(dataSize int, vr valueReader) (Property, error) {
	r, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading linear color r: %w", err)
	}

	g, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading linear color g: %w", err)
	}

	b, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading linear color b: %w", err)
	}

	a, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading linear color a: %w", err)
	}

	return &linearColorProperty{r, g, b, a}, nil
}

func init() {
	addStructType("Color", readColorStruct)
	addStructType("LinearColor", readLinearColorStruct)
}
