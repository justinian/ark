package main

import (
	"encoding/json"
	"fmt"
)

type unknownStructProperty struct {
	kind  Name
	value []byte
}

func (p *unknownStructProperty) Type() PropertyType { return UnknownProperty }

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

	const line_width = 16

	fmt.Printf("Unknown property: [%4x bytes] %s\n", dataSize, kind)
	fmt.Println("========================================================")
	var offset int = 0
	for offset < dataSize {
		fmt.Printf("%04x:", offset)
		var limit int = offset + line_width
		for offset < limit && offset < dataSize {
			if offset%4 == 0 && offset%line_width != 0 {
				fmt.Printf(" ")
			}
			fmt.Printf(" %02x", value[offset])
			offset++
		}
		fmt.Println("")
	}
	fmt.Println("========================================================")
	fmt.Println("")

	return &unknownStructProperty{kind, value}, nil
}
