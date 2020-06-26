package main

import (
	"fmt"
)

type Embed struct {
	Path  string
	Blobs [][][]byte
}

func ReadEmbed(a *Archive) (*Embed, error) {
	path, err := a.readString()
	if err != nil {
		return nil, err
	}

	partCount, err := a.readInt()
	if err != nil {
		return nil, err
	}

	parts := make([][][]byte, partCount)
	for i := range parts {
		blobCount, err := a.readInt()
		if err != nil {
			return nil, err
		}

		blobs := make([][]byte, blobCount)
		for j := range blobs {
			blobSize, err := a.readInt()
			if err != nil {
				return nil, err
			}

			blobSize *= 4 // blobSize is a count of 32bit values

			blob := make([]byte, blobSize)
			n, err := a.Read(blob)
			if err != nil {
				return nil, err
			} else if n != blobSize {
				return nil, fmt.Errorf("Failed to read blob size %d, got %d", blobSize, n)
			}

			blobs[j] = blob
		}

		parts[i] = blobs
	}

	return &Embed{path, parts}, nil
}
