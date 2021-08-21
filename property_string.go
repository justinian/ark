package ark

import (
	"fmt"
)

type stringProperty string

func (p stringProperty) Type() PropertyType { return StringPropertyType }
func (p stringProperty) String() string     { return string(p) }

func readStringProperty(dataSize int, vr valueReader) (Property, error) {
	value, err := vr.readString()
	if err != nil {
		return nil, fmt.Errorf("Reading string value: %w", err)
	}

	return stringProperty(value), nil
}

func init() {
	addPropertyType("StrProperty", 0, readStringProperty, nil)
}
