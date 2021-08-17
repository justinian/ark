package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type valueReader interface {
	io.Reader

	skip(size int) error
	subReader(length int) (valueReader, error)

	readInt() (int, error)
	readIntOfSize(size int) (int64, error)
	readFloat() (float32, error)
	readString() (string, error)
	readBool() (bool, error)
	readName() (Name, error)
}

type sliceValueReader struct {
	data      []byte
	nameTable []string
	offset    int
}

func (vr *sliceValueReader) skip(size int) error {
	vr.offset += size
	if vr.offset > len(vr.data) {
		return fmt.Errorf("Attempt to seek beyond end of sliceValueReader")
	}
	return nil
}

func (vr *sliceValueReader) readInt() (int, error) {
	i, err := vr.readIntOfSize(4)
	return int(i), err
}

func (vr *sliceValueReader) readIntOfSize(size int) (int64, error) {
	if vr.offset+size > len(vr.data) {
		return 0, fmt.Errorf("Attempt to read beyond end of sliceValueReader")
	}

	var err error
	var value int64

	switch size {
	case 1:
		var tmp int8
		err = binary.Read(vr, binary.LittleEndian, &tmp)
		value = int64(tmp)
	case 2:
		var tmp int16
		err = binary.Read(vr, binary.LittleEndian, &tmp)
		value = int64(tmp)
	case 4:
		var tmp int32
		err = binary.Read(vr, binary.LittleEndian, &tmp)
		value = int64(tmp)
	case 8:
		err = binary.Read(vr, binary.LittleEndian, &value)
	default:
		err = fmt.Errorf("Invalid int size: %d", size)
	}

	return value, err
}

func (vr *sliceValueReader) readFloat() (float32, error) {
	if vr.offset+4 > len(vr.data) {
		return 0, fmt.Errorf("Attempt to read beyond end of sliceValueReader")
	}

	var number float32
	if err := binary.Read(vr, binary.LittleEndian, &number); err != nil {
		return 0, err
	}
	return number, nil
}

func (vr *sliceValueReader) readString() (string, error) {
	length, err := vr.readInt()
	if err != nil {
		return "", fmt.Errorf("Reading string length:\n%w", err)
	}

	if length == 0 {
		return "", nil
	} else if length < 0 {
		length *= -2
	}

	data := make([]byte, length)
	if err := binary.Read(vr, binary.LittleEndian, &data); err != nil {
		return "", fmt.Errorf("Reading %d-byte string:\n%w", length, err)
	}

	return string(data[:len(data)-1]), nil
}

func (vr *sliceValueReader) readBool() (bool, error) {
	number, err := vr.readInt()
	if err != nil {
		return false, err
	}
	return number != 0, nil
}

func (vr *sliceValueReader) readName() (Name, error) {
	index, err := vr.readInt()
	if err != nil {
		return Name{}, err
	}

	index -= 1 // why, Ark?

	if index < 0 || index >= len(vr.nameTable) {
		return Name{}, fmt.Errorf("Invalid nameTable index %d", index)
	}

	instance, err := vr.readInt()
	if err != nil {
		return Name{}, err
	}

	return Name{
		Name:     vr.nameTable[index],
		Instance: instance,
	}, nil
}

func (vr *sliceValueReader) subReader(length int) (valueReader, error) {
	svr := &sliceValueReader{
		data:      make([]byte, length),
		nameTable: vr.nameTable,
	}

	_, err := io.ReadFull(vr, svr.data)
	if err != nil {
		return nil, err
	}

	return svr, nil
}

func (vr *sliceValueReader) Read(out []byte) (int, error) {
	if vr.offset+len(out) > len(vr.data) {
		return 0, fmt.Errorf("Attempt to read beyond end of sliceValueReader")
	}

	r := bytes.NewReader(vr.data[vr.offset:])
	n, err := r.Read(out)
	vr.offset += n
	return n, err
}
