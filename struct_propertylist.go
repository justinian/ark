package main

import (
	"fmt"
)

func (p PropertyMap) Type() PropertyType { return StructPropertyListProperty }

func (p PropertyMap) String() string {
	return fmt.Sprintf("StructPropertyListProperty(%d entries)", len(p))
}

func readPropertyListStruct(dataSize int, vr valueReader) (Property, error) {
	propMap, err := readPropertyMap(vr)
	if err != nil {
		return nil, err
		//return nil, fmt.Errorf("Reading PropertyList struct [%d bytes]:\n%w", dataSize, err)
	}

	return propMap, nil
}

func init() {
	addStructType("PropertyList", readPropertyListStruct)
	addStructType("CustomItemByteArrays", readPropertyListStruct)
}
