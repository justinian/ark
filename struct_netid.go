package ark

import "fmt"

type netIdProperty struct {
	unknown int
	NetId   string `json:"netId"`
}

func (p *netIdProperty) Type() PropertyType { return StructNetIdPropertyType }

func (p *netIdProperty) String() string {
	return fmt.Sprintf("StructNetIdProperty()")
}

func readNetIdStruct(dataSize int, vr valueReader) (Property, error) {
	unk, err := vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Error reading netid unknown int:\n%w", err)
	}

	netid, err := vr.readString()
	if err != nil {
		return nil, fmt.Errorf("Error reading netid:\n%w", err)
	}

	return &netIdProperty{unk, netid}, nil
}

func init() {
	//addStructType("UniqueNetIdRepl", readNetIdStruct)
}
