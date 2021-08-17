package main

import (
	"fmt"
)

type structReader func(dataSize int, vr valueReader) (Property, error)

var structReaders map[string]structReader

func addStructType(name string, reader structReader) {
	if structReaders == nil {
		structReaders = make(map[string]structReader)
	}
	structReaders[name] = reader
}

func readStructProperty(dataSize int, vr valueReader) (Property, error) {
	kind, err := vr.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading struct type:\n%w", err)
	}

	if dataSize == 0 {
		panic(kind.Name)
	}

	sub, err := vr.subReader(dataSize)
	if err != nil {
		return nil, fmt.Errorf("Creating struct valueReader:\n%w", err)
	}

	reader, ok := structReaders[kind.Name]
	if !ok {
		prop, err := readUnknownStruct(kind, dataSize, sub)
		if err != nil {
			return nil, fmt.Errorf("Reading unknown struct type:%s size:%d\n%w", kind, dataSize, err)
		}
		return prop, nil
	}

	prop, err := reader(dataSize, sub)
	if err != nil {
		return nil, fmt.Errorf("Reading structProperty type:%s size:%d\n%w", kind, dataSize, err)
	}
	return prop, nil
}

func readStructArray(count, dataSize int, vr valueReader) ([]Property, error) {
	var err error

	itemSize := dataSize / count

	var kind string
	switch itemSize {
	case 4:
		kind = "Color"
	case 12:
		kind = "Vector"
	case 16:
		kind = "LinearColor"
	default:
		kind = "PropertyList"
	}

	reader, ok := structReaders[kind]
	if !ok {
		return nil, fmt.Errorf("No defined struct reader for: %s", kind)
	}

	data := make([]Property, count)
	for i := range data {
		data[i], err = reader(itemSize, vr)
		if err != nil {
			return nil, fmt.Errorf("Reading structArray type:%s index:%d size:%d\n%w", kind, i, itemSize, err)
		}
	}

	return data, nil
}

func init() {
	addPropertyType("StructProperty", 0, readStructProperty, readStructArray)
}
