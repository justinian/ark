package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type Archive struct {
	nameTable []string
	stream    io.ReadSeeker

	version           int16
	hibernationOffset int32
	propertiesOffset  int32
}

func NewArchive(path string) (*Archive, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var header struct {
		Version               int16
		HibernationOffset     int32
		Reserved0             int32
		NameTableOffset       int32
		PropertiesBlockOffset int32
	}

	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	archive := Archive{
		stream: file,

		version:           header.Version,
		hibernationOffset: header.HibernationOffset,
		propertiesOffset:  header.PropertiesBlockOffset,
	}

	log.Printf("Parsed archive header version %d", archive.version)

	if archive.version <= 5 {
		return nil, fmt.Errorf("Save format version %d file is too old.", archive.version)
	}

	err = archive.readNameTable(header.NameTableOffset)
	if err != nil {
		return nil, fmt.Errorf("Reading name table: %w", err)
	}

	return &archive, nil
}

func (a *Archive) readNameTable(offset int32) error {
	savedPosition, err := a.stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("Determining current file position: %w", err)
	}

	_, err = a.stream.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("Seeking to name table offset: %w", err)
	}

	nameCount, err := a.readInt()
	if err != nil {
		return fmt.Errorf("Reading name table count: %w", err)
	}

	log.Printf("Reading %d name table entries", nameCount)

	a.nameTable = make([]string, nameCount)
	for i := range a.nameTable {
		a.nameTable[i], err = a.readString()
		if err != nil {
			return fmt.Errorf("Reading name entry: %w", err)
		}
	}

	_, err = a.stream.Seek(savedPosition, io.SeekStart)
	if err != nil {
		return fmt.Errorf("Returning from name table offset: %w", err)
	}
	return nil
}

func (a *Archive) readInt() (int, error) {
	var number uint32
	if err := binary.Read(a, binary.LittleEndian, &number); err != nil {
		return 0, err
	}
	return int(number), nil
}

func (a *Archive) readIntOfSize(size int) (int64, error) {
	var err error
	var value int64

	switch size {
	case 1:
		var tmp int8
		err = binary.Read(a, binary.LittleEndian, &tmp)
		value = int64(tmp)
	case 2:
		var tmp int16
		err = binary.Read(a, binary.LittleEndian, &tmp)
		value = int64(tmp)
	case 4:
		var tmp int32
		err = binary.Read(a, binary.LittleEndian, &tmp)
		value = int64(tmp)
	case 8:
		err = binary.Read(a, binary.LittleEndian, &value)
	default:
		err = fmt.Errorf("Invalid int size: %d", size)
	}

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (a *Archive) readString() (string, error) {
	var length int32
	if err := binary.Read(a, binary.LittleEndian, &length); err != nil {
		return "", err
	}

	if length < 0 {
		length *= -2
	}

	data := make([]byte, length)
	if err := binary.Read(a, binary.LittleEndian, &data); err != nil {
		return "", err
	}

	return string(data[:len(data)-1]), nil
}

func (a *Archive) readBool() (bool, error) {
	number, err := a.readInt()
	if err != nil {
		return false, err
	}
	return number != 0, nil
}

func (a *Archive) readName() (Name, error) {
	index, err := a.readInt()
	if err != nil {
		return Name{}, err
	}

	index -= 1 // why, Ark?

	if index < 0 || index >= len(a.nameTable) {
		return Name{}, fmt.Errorf("Invalid nameTable index %d", index)
	}

	instance, err := a.readInt()
	if err != nil {
		return Name{}, err
	}

	return Name{
		Name:     a.nameTable[index],
		Instance: instance,
	}, nil
}

func (a *Archive) readFloat() (float32, error) {
	var number float32
	if err := binary.Read(a, binary.LittleEndian, &number); err != nil {
		return 0, err
	}
	return number, nil
}

func (a *Archive) readProperties(offset int) (PropertyMap, error) {
	savedPosition, err := a.stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Determining current file position: %w", err)
	}

	totalOffset := int64(offset) + int64(a.propertiesOffset)
	_, err = a.stream.Seek(totalOffset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("Seeking to property offset: %w", err)
	}

	properties, err := readPropertyMap(a)
	if err != nil {
		return nil, fmt.Errorf("Reading property map: %w", err)
	}

	_, err = a.stream.Seek(savedPosition, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("Returning from property offset: %w", err)
	}

	return properties, nil
}

func (a *Archive) Read(b []byte) (int, error) {
	return a.stream.Read(b)
}
