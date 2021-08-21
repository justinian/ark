package ark

import (
	"encoding/json"
	"fmt"
)

type unknownStructProperty struct {
	kind  Name
	value []byte
}

func (p *unknownStructProperty) Type() PropertyType { return UnknownPropertyType }

func (p *unknownStructProperty) MarshalJSON() ([]byte, error) {
	val := struct {
		UnknownType string `json:"unknownType"`
		Length      int    `json:"length"`
	}{p.kind.String(), len(p.value)}
	return json.Marshal(val)
}

func (p *unknownStructProperty) String() string {
	return fmt.Sprintf("UnknownProperty(%s: %d bytes)", p.kind, len(p.value))
}

func readUnknownStruct(kind Name, dataSize int, vr valueReader) (Property, error) {
	value := make([]byte, dataSize)
	n, err := vr.Read(value)
	if err != nil {
		return nil, fmt.Errorf("Reading unknown value:\n%w", err)
	} else if n != dataSize {
		return nil, fmt.Errorf("Reading unknown value: short read %d/%d", n, dataSize)
	}

	return &unknownStructProperty{kind, value}, nil
}

func dumpBytes(data []byte) {
	const line_width = 16

	fmt.Println("========================================================")

	dataSize := len(data)
	var offset int = 0

	for offset < dataSize {
		fmt.Printf("%04x:", offset)
		var limit int = offset + line_width
		for offset < limit && offset < dataSize {
			if offset%4 == 0 && offset%line_width != 0 {
				fmt.Printf(" ")
			}
			fmt.Printf(" %02x", data[offset])
			offset++
		}
		fmt.Println("")
	}

	fmt.Println("========================================================")
	fmt.Println("")
}

func (p *unknownStructProperty) dump() {

	fmt.Printf("Unknown property: [%4x bytes] %s\n", len(p.value), p.kind)
	dumpBytes(p.value)
}
