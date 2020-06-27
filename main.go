package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <savefile>\n", os.Args[0])
		os.Exit(1)
	}

	archive, err := NewArchive(os.Args[1])
	if err != nil {
		log.Fatalf("Could not open save file: %v", err)
	}

	save, err := ReadSaveGame(archive)
	if err != nil {
		log.Fatalf("Could not read save game: %v", err)
	}

	_ = save
}
