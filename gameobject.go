package ark

import (
	"fmt"

	"github.com/google/uuid"
)

type GameObject struct {
	UUID       uuid.UUID   `json:"uuid"`
	ClassName  Name        `json:"className"`
	Names      []Name      `json:"names"`
	IsItem     bool        `json:"isItem"`
	IsCryopod  bool        `json:"isCryopod"`
	Location   Location    `json:"location"`
	Properties PropertyMap `json:"properties"`

	Parent           *GameObject `json:"parent"`
	FromDataFile     bool        `json:"fromDataFile"`
	DataFileIndex    int         `json:"dataFileIndex"`
	propertiesOffset int
}

func (o *GameObject) String() string {
	item := ""
	if o.IsItem {
		item = "*"
	}
	return fmt.Sprintf("Object %-12s: %s%s", o.ClassName, o.Names, item)
}

func readGameObject(vr valueReader) (*GameObject, error) {
	uuidBytes := make([]byte, 16)
	n, err := vr.Read(uuidBytes)
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

	obj.ClassName, err = vr.readName()
	if err != nil {
		return nil, fmt.Errorf("Reading class name: %w", err)
	}

	obj.IsItem, err = vr.readBool()
	if err != nil {
		return nil, fmt.Errorf("Reading item flag: %w", err)
	}

	nameCount, err := vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading name count: %w", err)
	}

	obj.Names = make([]Name, nameCount)
	for i := range obj.Names {
		obj.Names[i], err = vr.readName()
		if err != nil {
			return nil, fmt.Errorf("Reading name: %w", err)
		}
	}

	obj.FromDataFile, err = vr.readBool()
	if err != nil {
		return nil, fmt.Errorf("Reading data file flag: %w", err)
	}

	obj.DataFileIndex, err = vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading data file index: %w", err)
	}

	hasLocationData, err := vr.readBool()
	if err != nil {
		return nil, fmt.Errorf("Reading location flag: %w", err)
	}

	if hasLocationData {
		obj.Location, err = readLocation(vr)
		if err != nil {
			return nil, fmt.Errorf("Reading location: %w", err)
		}
	}

	obj.propertiesOffset, err = vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading properties offset: %w", err)
	}

	_, err = vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading reserved field: %w", err)
	}

	return &obj, nil
}
