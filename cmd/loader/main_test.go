package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/justinian/ark"
)

func TestReadSave(t *testing.T) {
	archive, err := ark.OpenArchive("Ragnarok.ark")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open save file:\n%s\n", err)
		os.Exit(1)
	}

	_, err = ark.ReadSaveGame(archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read save game:\n%+v\n", err)
		os.Exit(1)
	}
}
