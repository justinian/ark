package main

import "fmt"

const (
	ObjectTypeId   int = 0
	ObjectTypePath int = 1
)

type objectProperty struct {
	Id   int  `json:"id"`
	Path Name `json:"path"`
}

func (p *objectProperty) Type() PropertyType { return ObjectProperty }

func (p *objectProperty) String() string {
	var ref string
	if p.Id != 0 {
		ref = fmt.Sprintf("int:%d", p.Id)
	} else {
		ref = p.Path.Name
	}

	return fmt.Sprintf("ObjectProperty(%s)", ref)
}

func readObjectProperty(dataSize int, vr valueReader) (Property, error) {
	objType, err := vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading object type: %w", err)
	}

	var objId int
	var objPath Name

	switch objType {
	case ObjectTypeId:
		objId, err = vr.readInt()
	case ObjectTypePath:
		objPath, err = vr.readName()
	default:
		return nil, fmt.Errorf("Unsupported object reference typ %d", objType)
	}

	return &objectProperty{
		Id:   objId,
		Path: objPath,
	}, nil
}

func init() {
	addPropertyType("ObjectProperty", 0, readObjectProperty, nil)
}
