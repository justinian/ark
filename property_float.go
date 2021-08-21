package ark

import (
	"encoding/binary"
	"fmt"
)

type FloatProperty float64

func (p FloatProperty) Type() PropertyType { return FloatPropertyType }

func (p FloatProperty) String() string {
	return fmt.Sprintf("FloatProperty(%f)", float64(p))
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

	return FloatProperty(value), nil
}

type structDoublesProperty []float64

func (p structDoublesProperty) Type() PropertyType { return StructDoublesPropertyType }

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
