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

func createDatabase(dbname string, specfiles, savefiles []string) error {
	db, err := openDatabase(dbname)
	if err != nil {
		return fmt.Errorf("Error opening database %s:\n%w", dbname, err)
	}
	defer db.Close()

	saves := make([]*ark.SaveGame, len(savefiles))
	for i, f := range savefiles {
		save, err := loadSavefile(f)
		if err != nil {
			return fmt.Errorf("Loading '%s':\n%w", f, err)
		}
		saves[i] = save
	}

	classNames, err := readSpecFiles(specfiles...)
	if err != nil {
		return fmt.Errorf("Reading spec files:\n%w", err)
	}

	err = processSaves(db, saves, classNames)
	if err != nil {
		return fmt.Errorf("Writing database:\n%w", err)
	}

	return nil
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

	tmp := output + ".tmp"
	err := createDatabase(tmp, specfiles, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = os.Rename(tmp, output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Renaming output file: %v\n", err)
		os.Exit(1)
	}
}
