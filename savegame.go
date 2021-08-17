package main

import (
	"fmt"
)

type SaveGame struct {
	GameTime      float32                  `json:"gameTime"`
	SaveCount     int                      `json:"saveCount"`
	DataFiles     []string                 `json:"dataFiles"`
	EmbeddedData  []*Embed                 `json:"embeddedData"`
	DataFileMap   map[int][]string         `json:"dataFileMap"`
	Objects       []*GameObject            `json:"gameObjects"`
	ObjectMap     map[string][]*GameObject `json:"objectMap"`
	StoredObjects []*GameObject            `json:"storedObjects"`
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

	s.DataFiles, err = a.readStringTable()
	if err != nil {
		return nil, err
	}

	err = s.readEmbeds(a)
	if err != nil {
		return nil, err
	}

	err = s.readDataFileMap(a)
	if err != nil {
		return nil, err
	}

	err = s.readGameObjects(a)
	if err != nil {
		return nil, err
	}

	s.findCryopods(a)

	return &s, nil
}

func (s *SaveGame) readEmbeds(a *Archive) error {
	numEmbeds, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of embedded entries:\n%w", err)
	}

	fmt.Printf("Reading %d embedded data items\n", numEmbeds)

	s.EmbeddedData = make([]*Embed, numEmbeds)
	for i := range s.EmbeddedData {
		embed, err := ReadEmbed(a)
		if err != nil {
			return fmt.Errorf("Reading embedded entry:\n%w", err)
		}
		s.EmbeddedData[i] = embed
	}

	return nil
}

func (s *SaveGame) readDataFileMap(a *Archive) error {
	mapCount, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of object map entries:\n%w", err)
	}

	fmt.Printf("Reading %d object map entries\n", mapCount)

	s.DataFileMap = make(map[int][]string)
	for i := 0; i < mapCount; i++ {
		level, err := a.readInt()
		if err != nil {
			return fmt.Errorf("Reading object map level:\n%w", err)
		}

		count, err := a.readInt()
		if err != nil {
			return fmt.Errorf("Reading object map count:\n%w", err)
		}

		fmt.Printf("   Object map for level %d: %d entries\n", level, count)

		names := make([]string, count)
		for j := range names {
			names[j], err = a.readString()
			if err != nil {
				return fmt.Errorf("Reading object map name:\n%w", err)
			}
		}

		s.DataFileMap[level] = names
	}

	return nil
}

func (s *SaveGame) readGameObjects(a *Archive) error {
	count, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of objects:\n%w", err)
	}

	fmt.Printf("Reading %d game objects\n", count)

	s.Objects = make([]*GameObject, count)
	s.ObjectMap = make(map[string][]*GameObject, count)
	for i := range s.Objects {
		obj, err := readGameObject(a)
		if err != nil {
			return fmt.Errorf("Reading object:\n%w", err)
		}

		s.Objects[i] = obj
		for i := range obj.Names {
			key := obj.Names[i].String()
			s.ObjectMap[key] = append(s.ObjectMap[key], obj)
		}
	}

	for _, o := range s.Objects {
		o.Properties, err = a.readProperties(o.propertiesOffset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SaveGame) readDino(vr valueReader) error {
	return nil
}

func (s *SaveGame) findCryopods(a *Archive) {
	for _, o := range s.Objects {
		if o.isCryopod() {
			customItemDatas := o.Properties.Get("CustomItemDatas", 0)
			if customItemDatas == nil {
				continue
			}

			customItemDataArray, ok := customItemDatas.(*arrayProperty)
			if !ok {
				fmt.Printf("customItemDatas is not an array property!\n")
				continue
			}

			//var dino PropertyMap
			for _, prop := range customItemDataArray.Properties {
				propMap, ok := prop.(PropertyMap)
				if !ok {
					fmt.Printf("%v is not an property map!\n", prop)
					continue
				}

				nameProp := propMap.Get("CustomDataName", 0)
				if nameProp == nil {
					fmt.Printf("%v has no CustomDataName\n", prop)
					continue
				}

				name, ok := nameProp.(Name)
				if !ok {
					fmt.Printf("%v is not an name property!\n", nameProp)
					continue
				}

				if name.Name != "Dino" {
					continue
				}

				customBytesProp := propMap.Get("CustomDataBytes", 0)
				if customBytesProp == nil {
					continue
				}

				customBytesMap, ok := customBytesProp.(PropertyMap)
				if !ok {
					continue
				}

				byteArraysProp := customBytesMap.Get("ByteArrays", 0)
				if byteArraysProp == nil {
					continue
				}

				byteArraysArray, ok := byteArraysProp.(*arrayProperty)
				if !ok {
					continue
				}

				if len(byteArraysArray.Properties) == 0 {
					continue
				}

				creatureBytesMap, ok := byteArraysArray.Properties[0].(PropertyMap)
				if !ok {
					continue
				}

				creatureBytesProp := creatureBytesMap.Get("Bytes", 0)
				if creatureBytesProp == nil {
					continue
				}

				creatureBytes, ok := creatureBytesProp.(byteArrayProperty)
				if !ok {
					continue
				}

				vr := &sliceValueReader{
					data:      []byte(creatureBytes),
					nameTable: a.nameTable,
				}
				s.readDino(vr)
			}
		}
	}

}
