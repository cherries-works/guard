package manifest

import (
	"encoding/json"
	"fmt"
)

type PackageJSONM struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParsePackageJSONM(data []byte) []string {
	dependencies := []string{}

	// parsing package.json
	if !json.Valid(data) {
		return dependencies
	}

	var parsed PackageJSONM
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range parsed.Dependencies {
		d := fmt.Sprintf("%s@%s", key, value)
		dependencies = append(dependencies, d)
	}

	return dependencies
}
