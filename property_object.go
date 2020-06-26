package main

import "fmt"

const (
	ObjectTypeId   int = 0
	ObjectTypePath int = 1
)

type objectProperty struct {
	name  Name
	index int
	id    int
	path  Name
}

func (p *objectProperty) Type() PropertyType { return ObjectProperty }
func (p *objectProperty) Name() Name         { return p.name }
func (p *objectProperty) Index() int         { return p.index }

func (p *objectProperty) String() string {
	var ref string
	if p.id != 0 {
		ref = fmt.Sprintf("int:%d", p.id)
	} else {
		ref = p.path.Name
	}

	return fmt.Sprintf("ObjectProperty %s %s [%d]", p.name, ref, p.index)
}

func readObjectProperty(name Name, dataSize, index int, a *Archive) (Property, error) {
	if dataSize < 8 {
		return nil, fmt.Errorf("Out of date object property: size too small")
	}

	objType, err := a.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading object type: %w", err)
	}

	var objId int
	var objPath Name

	switch objType {
	case ObjectTypeId:
		objId, err = a.readInt()
	case ObjectTypePath:
		objPath, err = a.readName()
	default:
		return nil, fmt.Errorf("Unsupported object reference typ %d", objType)
	}

	return &objectProperty{
		name:  name,
		index: index,
		id:    objId,
		path:  objPath,
	}, nil
}
