package ark

import (
	"fmt"
)

type PropertyType uint
type PropertyMap map[string]map[int]Property

func (pm PropertyMap) Get(name string, index int) Property {
	if props, ok := pm[name]; ok {
		if p, ok := props[index]; ok {
			return p
		}
	}
	return nil
}

func (pm PropertyMap) GetTyped(name string, index int, pt PropertyType) Property {
	p := pm.Get(name, index)
	if p != nil && p.Type() == pt {
		return p
	}
	return nil
}

func (pm PropertyMap) GetString(name string, index int) string {
	stringProp := pm.Get(name, index)
	if stringProp != nil {
		return stringProp.String()
	}
	return ""
}

func (pm PropertyMap) GetInt(name string, index int) int64 {
	intProp := pm.GetTyped(name, index, IntPropertyType)
	if intProp != nil {
		return intProp.(*IntProperty).value
	}
	return 0
}

func (pm PropertyMap) GetFloat(name string, index int) float64 {
	floatProp := pm.GetTyped(name, index, FloatPropertyType)
	if floatProp != nil {
		return float64(floatProp.(FloatProperty))
	}
	return 0
}

const (
	UnknownPropertyType PropertyType = iota

	IntPropertyType
	FloatPropertyType
	BoolPropertyType
	EnumPropertyType
	StringPropertyType
	NamePropertyType
	ArrayPropertyType
	ObjectPropertyType
	ByteArrayPropertyType

	StructColorPropertyType
	StructLinearColorPropertyType
	StructVectorPropertyType
	StructVector2DPropertyType
	StructQuatPropertyType

	StructNetIdPropertyType
	StructDoublesPropertyType

	StructPropertyListPropertyType
)

type Property interface {
	fmt.Stringer
	Type() PropertyType
}

type propertyReader func(int, valueReader) (Property, error)
type propertyArrayReader func(int, int, valueReader) ([]Property, error)

type propertyType struct {
	reader      propertyReader
	arrayReader propertyArrayReader
	defaultSize int
}

var propertyTypes map[string]propertyType

func addPropertyType(name string, defaultSize int, reader propertyReader, arrayReader propertyArrayReader) {
	if propertyTypes == nil {
		propertyTypes = make(map[string]propertyType, 20)
	}
	propertyTypes[name] = propertyType{reader: reader, arrayReader: arrayReader, defaultSize: defaultSize}
}

func readPropertyMap(vr valueReader) (PropertyMap, error) {
	properties := make(PropertyMap)

	count := 0
	for {
		name, err := vr.readName()
		if err != nil {
			return nil, fmt.Errorf("Reading property name:\n%w", err)
		}

		if name.IsNone() {
			break
		}

		propertyType, err := vr.readName()
		if err != nil {
			return nil, fmt.Errorf("Reading property type:\n%w", err)
		}

		dataSize, err := vr.readInt()
		if err != nil {
			return nil, fmt.Errorf("Reading property size:\n%w", err)
		}

		index, err := vr.readInt()
		if err != nil {
			return nil, fmt.Errorf("Reading property index:\n%w", err)
		}

		p, err := readProperty(propertyType.Name, dataSize, vr)
		if err != nil {
			return nil, fmt.Errorf("Reading propertyMap item %d: name:%s type:%s bytes:%d:\n%w", count, name, propertyType, dataSize, err)
		}

		if err != nil {
			return nil, fmt.Errorf("Reading property:\n%w", err)
		} else if p == nil {
			// nil can be returned without error for unhandled or ignored properties
			continue
		}

		key := name.String()
		propMap, ok := properties[key]
		if !ok {
			propMap = make(map[int]Property)
		}
		propMap[index] = p
		properties[key] = propMap
		count++
	}

	return properties, nil
}

func readProperty(name string, dataSize int, vr valueReader) (Property, error) {
	if propType, ok := propertyTypes[name]; ok {
		return propType.reader(dataSize, vr)
	}

	return nil, fmt.Errorf("Unknown property type %s", name)
}

func readPropertyArray(name string, count, dataSize int, vr valueReader) ([]Property, error) {
	propType, ok := propertyTypes[name]
	if !ok {
		return nil, fmt.Errorf("Unknown array property type %s", name)
	}

	if count == 0 {
		return nil, nil
	}

	if propType.arrayReader != nil {
		return propType.arrayReader(count, dataSize, vr)
	}

	var err error
	result := make([]Property, count)
	for i := range result {
		itemSize := propType.defaultSize
		result[i], err = propType.reader(itemSize, vr)
		if err != nil {
			return nil, fmt.Errorf("Reading basic array:\n%w", err)
		}
	}

	return result, nil
}
