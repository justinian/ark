package ark

import "fmt"

type arrayProperty struct {
	ArrayType  string     `json:"arrayType"`
	Properties []Property `json:"properties"`
}

func (p *arrayProperty) Type() PropertyType { return ArrayPropertyType }

func (p *arrayProperty) String() string {
	return fmt.Sprintf("ArrayProperty(%s[%d])", p.ArrayType, len(p.Properties))
}

type byteArrayProperty []byte

func (p byteArrayProperty) Type() PropertyType { return ByteArrayPropertyType }

func (p byteArrayProperty) String() string {
	return fmt.Sprintf("ByteArrayProperty[%d]", len(p))
}

func readArrayProperty(dataSize int, vr valueReader) (Property, error) {
	arrayType, err := vr.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading array type:\n%w", err)
	}

	return readArrayPropertyOfType(arrayType.Name, dataSize, vr)
}

func readArrayPropertyOfType(typeName string, dataSize int, vr valueReader) (Property, error) {
	count, err := vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading arrayProperty type:%s array length:\n%w", typeName, err)
	}

	dataSize -= 4 // Compensate for reading the count

	if typeName == "ByteProperty" {
		p := make(byteArrayProperty, dataSize)
		_, err := vr.Read(p)
		if err != nil {
			return nil, fmt.Errorf("Reading byteArrayProperty length:%d:\n%w", dataSize, err)
		}
		return p, nil
	}

	props, err := readPropertyArray(typeName, count, dataSize, vr)
	if err != nil {
		return nil, fmt.Errorf("Reading arrayProperty type:%s count:%d size:%d\n%w",
			typeName, count, dataSize, err)
	}

	return &arrayProperty{
		ArrayType:  typeName,
		Properties: props,
	}, nil
}

func init() {
	addPropertyType("ArrayProperty", 0, readArrayProperty, nil)
}
