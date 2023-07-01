package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justinian/ark"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <savefile>...", os.Args[0])
		os.Exit(1)
	}

	for _, filename := range os.Args[1:] {
		err := dumpSave(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error from %s: %v", filename, err)
			os.Exit(1)
		}
	}
}

func dumpSave(filename string) error {
	archive, err := ark.OpenArchive(filename)
	if err != nil {
		return fmt.Errorf("Opening archive: %w", err)
	}

	save, err := ark.ReadSaveGame(archive)
	if err != nil {
		return fmt.Errorf("Reading archive: %w", err)
	}

	outf, err := os.Create(filename + ".json")
	if err != nil {
		return fmt.Errorf("Creating output file: %w", err)
	}
	defer outf.Close()

	data, err := json.Marshal(save)
	if err != nil {
		return fmt.Errorf("Marshalling save data: %w", err)
	}

	_, err = outf.Write(data)
	if err != nil {
		return fmt.Errorf("Writing output file: %w", err)
	}

	return nil
}
