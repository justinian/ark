package ark

import (
	"fmt"
	"strings"
)

type SaveGame struct {
	GameTime     float32              `json:"gameTime"`
	SaveCount    int                  `json:"saveCount"`
	DataFiles    []string             `json:"dataFiles"`
	EmbeddedData []*Embed             `json:"embeddedData"`
	DataFileMap  map[int][]string     `json:"dataFileMap"`
	Objects      [][]*GameObject      `json:"objectLists"`
	NameMap      map[Name]*GameObject `json:"-"`
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

	objects, err := s.readGameObjects(a)
	if err != nil {
		return nil, err
	}

	s.NameMap = make(map[Name]*GameObject)
	for _, o := range objects {
		if len(o.Names) == 1 {
			s.NameMap[o.Names[0]] = o
		}
	}

	storedObjects, err := s.findCryopods(a, objects)
	if err != nil {
		return nil, err
	}

	s.Objects = make([][]*GameObject, len(storedObjects)+1)
	s.Objects = append(s.Objects, objects)
	s.Objects = append(s.Objects, storedObjects...)

	return &s, nil
}

func (s *SaveGame) readEmbeds(a *Archive) error {
	numEmbeds, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading number of embedded entries:\n%w", err)
	}

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

func (s *SaveGame) readGameObjects(vr valueReader) ([]*GameObject, error) {
	count, err := vr.readInt()
	if err != nil {
		return nil, fmt.Errorf("Reading number of objects:\n%w", err)
	}

	objects := make([]*GameObject, count)
	for i := range objects {
		obj, err := readGameObject(vr)
		if err != nil {
			return nil, fmt.Errorf("Reading object %d/%d:\n%w", i, count, err)
		}

		obj.Properties, err = vr.readProperties(obj.propertiesOffset)
		if err != nil {
			return nil, fmt.Errorf("Reading object properties:\n%w", err)
		}

		objects[i] = obj
	}

	return objects, nil
}

func (s *SaveGame) findCryopods(a *Archive, objects []*GameObject) ([][]*GameObject, error) {
	var allStoredObjects [][]*GameObject

	for _, o := range objects {
		if strings.Contains(o.ClassName.Name, "Cryop") ||
			strings.Contains(o.ClassName.Name, "SoulTrap_") {

			var parent *GameObject
			ownerInventoryProp := o.Properties.Get("OwnerInventory", 0)
			if ownerInventoryProp != nil {
				containerId := ownerInventoryProp.(*ObjectProperty).Id
				container := objects[containerId]
				parentName := container.Names[len(container.Names)-1]
				parent = s.NameMap[parentName]
			}

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
					data: []byte(creatureBytes),
				}

				storedObjects, err := s.readGameObjects(vr)
				if err != nil {
					return nil, fmt.Errorf("Reading stored objects:\n%w", err)
				}

				for _, o := range storedObjects {
					o.IsCryopod = true
					o.Parent = parent
				}

				allStoredObjects = append(allStoredObjects, storedObjects)
			}
		}
	}

	return allStoredObjects, nil
}
