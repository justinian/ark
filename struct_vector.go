package ark

import "fmt"

type vector2DProperty struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (p *vector2DProperty) Type() PropertyType { return StructVector2DPropertyType }

func (p *vector2DProperty) String() string {
	return fmt.Sprintf("StructVector2DProperty()")
}

func readVector2DStruct(dataSize int, vr valueReader) (Property, error) {
	x, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading vector x: %w", err)
	}

	y, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading vector y: %w", err)
	}

	return &vector2DProperty{x, y}, nil
}

type vectorProperty struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func (p *vectorProperty) Type() PropertyType { return StructVectorPropertyType }

func (p *vectorProperty) String() string {
	return fmt.Sprintf("StructVectorProperty()")
}

func readVectorStruct(dataSize int, vr valueReader) (Property, error) {
	x, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading vector x: %w", err)
	}

	y, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading vector y: %w", err)
	}

	z, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading vector z: %w", err)
	}

	return &vectorProperty{x, y, z}, nil
}

type quatProperty struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
	W float32 `json:"w"`
}

func (p *quatProperty) Type() PropertyType { return StructQuatPropertyType }

func (p *quatProperty) String() string {
	return fmt.Sprintf("StructQuatProperty()")
}

func readQuatStruct(dataSize int, vr valueReader) (Property, error) {
	x, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading quat x: %w", err)
	}

	y, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading quat y: %w", err)
	}

	z, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading quat z: %w", err)
	}

	w, err := vr.readFloat()
	if err != nil {
		return nil, fmt.Errorf("Error reading quat w: %w", err)
	}

	return &quatProperty{x, y, z, w}, nil
}

func init() {
	addStructType("Vector2D", readVector2DStruct)
	addStructType("Vector", readVectorStruct)
	addStructType("Rotator", readVectorStruct)
	addStructType("Quat", readQuatStruct)
}
