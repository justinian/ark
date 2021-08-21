package main

import (
	"fmt"
	"os"

	"github.com/justinian/ark"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
)

func loadSavefile(filename string) (*ark.SaveGame, error) {
	archive, err := ark.OpenArchive(filename)
	if err != nil {
		return nil, fmt.Errorf("Could not open save file:\n%w", err)
	}

	save, err := ark.ReadSaveGame(archive)
	if err != nil {
		return nil, fmt.Errorf("Could not read save game:\n%w", err)
	}

	return save, nil
}

func main() {
	var output string
	var specfiles []string
	pflag.StringVarP(&output, "out", "o", "ark.db", "Filename of the database to create")
	pflag.StringArrayVarP(&specfiles, "spec", "s", nil, "JSON species/item files to load")
	pflag.Parse()

	args := pflag.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <savefile> ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s -h  for help\n", os.Args[0])
		os.Exit(1)
	}

	db, err := openDatabase(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database %s: %v\n", output, err)
		os.Exit(1)
	}
	defer db.Close()

	saves := make([]*ark.SaveGame, len(args))
	for i, f := range args {
		save, err := loadSavefile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Loading '%s':\n%v\n", f, err)
			os.Exit(1)
		}
		saves[i] = save
	}

	classNames, err := readSpecFiles(specfiles...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Reading spec files:\n%v\n", err)
		os.Exit(1)
	}

	err = processSaves(db, saves, classNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Writing database:\n%v\n", err)
		os.Exit(1)
	}
}
