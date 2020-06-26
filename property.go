package main

import (
	"fmt"
)

type PropertyType uint

const (
	IntProperty PropertyType = iota
	FloatProperty
	BoolProperty
	EnumProperty
	ArrayProperty
	StructProperty
	ObjectProperty
)

type Property interface {
	fmt.Stringer
	Name() Name
	Type() PropertyType
	Index() int
}

func readProperty(a *Archive) (Property, error) {
	name, err := a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading property name: %w", err)
	}

	if name.IsNone() {
		return nil, nil
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
		return readIntProperty(name, true, dataSize, index, a)

	case "UIntProperty":
		fallthrough
	case "UInt8Property":
		fallthrough
	case "UInt16Property":
		fallthrough
	case "UInt32Property":
		fallthrough
	case "UInt64Property":
		return readIntProperty(name, false, dataSize, index, a)

	case "FloatProperty":
		fallthrough
	case "DoubleProperty":
		return readFloatProperty(name, dataSize, index, a)

	case "BoolProperty":
		return readBoolProperty(name, dataSize, index, a)

	case "ByteProperty":
		return readEnumProperty(name, dataSize, index, a)

	case "ArrayProperty":
		return readArrayProperty(name, dataSize, index, a)

	case "StructProperty":
		return readStructProperty(name, dataSize, index, a)

	case "ObjectProperty":
		return readObjectProperty(name, dataSize, index, a)

	default:
		return nil, fmt.Errorf("Unknown property type %s", propertyType)
	}
}
