package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <savefile>\n", os.Args[0])
		os.Exit(1)
	}

	archive, err := NewArchive(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open save file:\n%s\n", err)
		os.Exit(1)
	}

	save, err := ReadSaveGame(archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read save game:\n%+v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(save, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not serialize save game:\n%+v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(fmt.Sprintf("%s.json", os.Args[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create JSON file:\n%+v\n", err)
		os.Exit(1)
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write JSON file:\n%+v\n", err)
		os.Exit(1)
	}

	file.Close()
}
