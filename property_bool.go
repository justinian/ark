package ark

import (
	"encoding/json"
	"fmt"
)

type BoolProperty struct {
	Value bool
}

func (p *BoolProperty) Type() PropertyType { return BoolPropertyType }

func (p *BoolProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *BoolProperty) String() string {
	return fmt.Sprintf("BoolProperty(%v)", p.Value)
}

func readBoolProperty(dataSize int, vr valueReader) (Property, error) {
	value, err := vr.readIntOfSize(1)
	if err != nil {
		return nil, fmt.Errorf("Reading bool value: %w", err)
	}

	return &BoolProperty{
		Value: value != 0,
	}, nil
}

func init() {
	addPropertyType("BoolProperty", 1, readBoolProperty, nil)
}
