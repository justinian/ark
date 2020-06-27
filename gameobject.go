package main

import (
	"fmt"

	"github.com/google/uuid"
)

type GameObject struct {
	UUID       uuid.UUID
	ClassName  Name
	Names      []Name
	IsItem     bool
	Location   Location
	Properties PropertyMap

	FromDataFile     bool
	DataFileIndex    int
	propertiesOffset int
}

func (o *GameObject) String() string {
	item := ""
	if o.IsItem {
		item = "*"
	}
	return fmt.Sprintf("Object %-12s: %s%s", o.ClassName, o.Names, item)
}

func readGameObject(a *Archive) (*GameObject, error) {
	uuidBytes := make([]byte, 16)
	n, err := a.Read(uuidBytes)
	if err != nil {
		return nil, fmt.Errorf("Reading object UUID: %w", err)
	} else if n != 16 {
		return nil, fmt.Errorf("Read wrong number of UUID bytes")
	}

	var obj GameObject
	obj.UUID, err = uuid.FromBytes(uuidBytes)
	if err != nil {
		return nil, fmt.Errorf("Invalid object UUID: %w", err)
	}

	obj.ClassName, err = a.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading class name: %w", err)
	}

	obj.IsItem, err = a.readBool()
	if err != nil {
		return nil, fmt.Errorf("Reading item flag: %w", err)
	}

	nameCount, err := a.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading name count: %w", err)
	}

	obj.Names = make([]Name, nameCount)
	for i := range obj.Names {
		obj.Names[i], err = a.readName()
		if err != nil {
			return nil, fmt.Errorf("Reading name: %w", err)
		}
	}

	obj.FromDataFile, err = a.readBool()
	if err != nil {
		return nil, fmt.Errorf("Reading data file flag: %w", err)
	}

	obj.DataFileIndex, err = a.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading data file index: %w", err)
	}

	hasLocationData, err := a.readBool()
	if err != nil {
		return nil, fmt.Errorf("Reading location flag: %w", err)
	}

	if hasLocationData {
		obj.Location, err = readLocation(a)
		if err != nil {
			return nil, fmt.Errorf("Reading location: %w", err)
		}
	}

	obj.propertiesOffset, err = a.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading properties offset: %w", err)
	}

	_, err = a.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading reserved field: %w", err)
	}

	return &obj, nil
}
