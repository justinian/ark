package ark

import (
	"fmt"
)

func (p Name) Type() PropertyType { return NamePropertyType }

func readNameProperty(dataSize int, vr valueReader) (Property, error) {
	value, err := vr.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading name value: %w", err)
	}

	return value, nil
}

func init() {
	addPropertyType("NameProperty", 0, readNameProperty, nil)
}
