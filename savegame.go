package main

import (
	"fmt"
	"log"
)

type SaveGame struct {
	GameTime     float32
	SaveCount    int
	DataFiles    []string
	EmbeddedData []*Embed
	ObjectMap    map[int][]string
	Objects      []*GameObject
}

func ReadSaveGame(a *Archive) (*SaveGame, error) {
	var err error
	var s SaveGame

	s.GameTime, err = a.readFloat()
	if err != nil {
		return nil, err
	}

	s.SaveCount, err = a.readInt()
	if err != nil {
		return nil, err
	}

	err = s.readDataFiles(a)
	if err != nil {
		return nil, err
	}

	err = s.readEmbeds(a)
	if err != nil {
		return nil, err
	}

	err = s.readObjectMap(a)
	if err != nil {
		return nil, err
	}

	err = s.readGameObjects(a)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *SaveGame) readDataFiles(a *Archive) error {
	numFiles, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of data files: %w", err)
	}

	log.Printf("Reading %d data files", numFiles)

	s.DataFiles = make([]string, numFiles)
	for i := range s.DataFiles {
		file, err := a.readString()
		if err != nil {
			return fmt.Errorf("Reading data file entry: %w", err)
		}
		s.DataFiles[i] = file
	}

	return nil
}

func (s *SaveGame) readEmbeds(a *Archive) error {
	numEmbeds, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of embedded entries: %w", err)
	}

	log.Printf("Reading %d embedded data items", numEmbeds)

	s.EmbeddedData = make([]*Embed, numEmbeds)
	for i := range s.EmbeddedData {
		embed, err := ReadEmbed(a)
		if err != nil {
			return fmt.Errorf("Reading embedded entry: %w", err)
		}
		s.EmbeddedData[i] = embed
	}

	return nil
}

func (s *SaveGame) readObjectMap(a *Archive) error {
	mapCount, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of object map entries: %w", err)
	}

	log.Printf("Reading %d object map entries", mapCount)

	s.ObjectMap = make(map[int][]string)
	for i := 0; i < mapCount; i++ {
		level, err := a.readInt()
		if err != nil {
			return fmt.Errorf("Reading object map level: %w", err)
		}

		count, err := a.readInt()
		if err != nil {
			return fmt.Errorf("Reading object map count: %w", err)
		}

		log.Printf("   Object map for level %d: %d entries", level, count)

		names := make([]string, count)
		for j := range names {
			names[j], err = a.readString()
			if err != nil {
				return fmt.Errorf("Reading object map name: %w", err)
			}
		}

		s.ObjectMap[level] = names
	}

	return nil
}

func (s *SaveGame) readGameObjects(a *Archive) error {
	count, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of objects: %w", err)
	}

	log.Printf("Reading %d game objects", count)

	s.Objects = make([]*GameObject, count)
	for i := range s.Objects {
		obj, err := readGameObject(a)
		if err != nil {
			return fmt.Errorf("Reading object: %w", err)
		}

		s.Objects[i] = obj
	}

	totalProperties := 0
	for _, o := range s.Objects {
		o.Properties, err = a.readProperties(o.propertiesOffset)
		if err != nil {
			return err
		}

		totalProperties += len(o.Properties)
		log.Println(o)
	}

	return nil
}
