package main

import (
	"fmt"
)

type PropertyType uint
type PropertyMap map[string]map[int]Property

const (
	IntProperty PropertyType = iota
	FloatProperty
	BoolProperty
	EnumProperty
	StringProperty
	NameProperty
	ArrayProperty
	StructProperty
	ObjectProperty
)

type Property interface {
	fmt.Stringer
	Type() PropertyType
}

func readPropertyMap(a *Archive) (PropertyMap, error) {
	properties := make(PropertyMap)

	for {
		name, err := a.readName()
		if err != nil {
			return nil, fmt.Errorf("Reading property name: %w", err)
		}

		if name.IsNone() {
			break
		}

		propertyType, err := a.readName()
		if err != nil {
			return nil, fmt.Errorf("Reading property type: %w", err)
		}

		dataSize, err := a.readInt()
		if err != nil {
			return nil, fmt.Errorf("Reading property size: %w", err)
		}

		index, err := a.readInt()
		if err != nil {
			return nil, fmt.Errorf("Reading property index: %w", err)
		}

		var p Property

		switch propertyType.Name {

		case "IntProperty":
			fallthrough
		case "Int8Property":
			fallthrough
		case "Int16Property":
			fallthrough
		case "Int32Property":
			fallthrough
		case "Int64Property":
			p, err = readIntProperty(true, dataSize, a)

		case "UIntProperty":
			fallthrough
		case "UInt8Property":
			fallthrough
		case "UInt16Property":
			fallthrough
		case "UInt32Property":
			fallthrough
		case "UInt64Property":
			p, err = readIntProperty(false, dataSize, a)

		case "FloatProperty":
			fallthrough
		case "DoubleProperty":
			p, err = readFloatProperty(dataSize, a)

		case "BoolProperty":
			p, err = readBoolProperty(dataSize, a)

		case "ByteProperty":
			p, err = readEnumProperty(dataSize, a)

		case "StrProperty":
			p, err = readStringProperty(dataSize, a)

		case "NameProperty":
			p, err = readNameProperty(dataSize, a)

		case "ArrayProperty":
			p, err = readArrayProperty(dataSize, a)

		case "StructProperty":
			p, err = readStructProperty(dataSize, a)

		case "ObjectProperty":
			p, err = readObjectProperty(dataSize, a)

		default:
			return nil, fmt.Errorf("Unknown property type %s", propertyType)
		}

		key := name.String()
		propMap, ok := properties[key]
		if !ok {
			propMap = make(map[int]Property)
		}
		propMap[index] = p
		properties[key] = propMap
	}

	return properties, nil
}
