package lock

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PackageLock struct {
	Packages map[string]PackageLockEntry `json:"packages"`
}

type PackageLockEntry struct {
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

func ParseNPML(data []byte) []string {
	dependencies := []string{}

	if !json.Valid(data) {
		return dependencies
	}

	var parsed PackageLock
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		fmt.Println(err)
	}

	for path, pkg := range parsed.Packages {
		if path == "" {
			// root project's dependencies
			continue
		}

		name := strings.TrimPrefix(path, "node_modules/")

		fmt.Println(name, pkg.Version)

		for dependency := range pkg.Dependencies {
			fmt.Println("  depends on:", dependency)
		}
	}

	return dependencies
}
