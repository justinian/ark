package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type classSpec struct {
	Name      string `json:"name"`
	Blueprint string `json:"bp"`
}

var keynames = []string{"items", "species"}
var clean = regexp.MustCompile(`\s+`)

func readSpecFiles(paths ...string) (map[string]string, error) {
	classNames := make(map[string]string)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("Opening spec file %f:\n%w", path, err)
		}

		jsonData, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("Reading spec file %f:\n%w", path, err)
		}

		var values map[string]json.RawMessage
		if err := json.Unmarshal(jsonData, &values); err != nil {
			return nil, fmt.Errorf("Loading spec file %f:\n%w", path, err)
		}

		for _, key := range keynames {
			if raw, ok := values[key]; ok {
				var specs []classSpec
				if err := json.Unmarshal(raw, &specs); err != nil {
					return nil, fmt.Errorf("Loading specs from file %f:\n%w", path, err)
				}

				for _, spec := range specs {
					parts := strings.Split(spec.Blueprint, ".")
					className := parts[len(parts)-1]
					classNames[className] = clean.ReplaceAllString(spec.Name, " ")
				}
			}
		}
	}

	return classNames, nil
}
