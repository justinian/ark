package ark

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type ArchiveHeader struct {
	Version               int16
	HibernationOffset     int32
	Reserved0             int32
	NameTableOffset       int32
	PropertiesBlockOffset int32
}

type Archive struct {
	sliceValueReader
	ArchiveHeader
}

func OpenArchive(path string) (*Archive, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading archive file:\n%w", err)
	}

	a := &Archive{
		sliceValueReader{data: data},
		ArchiveHeader{},
	}

	if err := binary.Read(a, binary.LittleEndian, &a.ArchiveHeader); err != nil {
		return nil, err
	}

	if a.Version <= 8 {
		return nil, fmt.Errorf("Save format version %d file is too old.", a.Version)
	}

	nameTableReader, _ := a.subReaderAt(int(a.NameTableOffset))
	a.nameTable, err = nameTableReader.readStringTable()
	if err != nil {
		return nil, fmt.Errorf("Reading archive name table:\n%w", err)
	}

	a.propertiesOffset = int(a.PropertiesBlockOffset)

	return a, nil
}
