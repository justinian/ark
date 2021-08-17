package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

type floatProperty struct {
	value float64
}

func (p *floatProperty) Type() PropertyType { return FloatProperty }

func (p *floatProperty) MarshalJSON() ([]byte, error) {
	if math.IsNaN(p.value) {
		return json.Marshal(nil)
	}
	return json.Marshal(p.value)
}

func (p *floatProperty) String() string {
	return fmt.Sprintf("FloatProperty(%f)", p.value)
}

func readFloatProperty(dataSize int, vr valueReader) (Property, error) {
	var err error
	var value float64

	switch dataSize {
	case 4:
		var tmp float32
		err = binary.Read(vr, binary.LittleEndian, &tmp)
		value = float64(tmp)
	case 8:
		err = binary.Read(vr, binary.LittleEndian, &value)
	}

	if err != nil {
		return nil, fmt.Errorf("Reading float value: %w", err)
	}

	return &floatProperty{
		value: value,
	}, nil
}

type structDoublesProperty []float64

func (p structDoublesProperty) Type() PropertyType { return StructDoublesProperty }

func (p structDoublesProperty) String() string {
	return fmt.Sprintf("StructDoublesProperty[%v]", len(p))
}

func readStructDoubles(vr valueReader) (Property, error) {
	count, err := vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading doubles array length:\n%w", err)
	}

	p := make(structDoublesProperty, count)
	for i := range p {
		err = binary.Read(vr, binary.LittleEndian, &p[i])
		if err != nil {
			return nil, fmt.Errorf("Reading doubles array value %d/%d:\n%w", i, count, err)
		}
	}

	return p, nil
}

func init() {
	addPropertyType("FloatProperty", 4, readFloatProperty, nil)
	addPropertyType("DoubleProperty", 8, readFloatProperty, nil)
	//addStructType("CustomItemDoubles", readStructDoubles)
}
