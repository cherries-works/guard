package parser

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/cherries-works/guard/internal/types"
	"golang.org/x/mod/modfile"
)

func ParseManifest(ecosystem types.Ecosystem, path string) []string {
	manifest_read, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	switch ecosystem {
	case types.NPM:
		return ParseNPM(manifest_read)
	case types.PNPM:
		return ParsePNPM(manifest_read)
	case types.Yarn:
		return ParseYarn(manifest_read)
	case types.Bun:
		return ParseBun(manifest_read)
	case types.Go:
		return ParseGo(manifest_read)
	case types.Java:
		return ParseJava(manifest_read)
	case types.Python:
		return ParsePython(string(manifest_read))
	case types.Cargo:
		return ParseCargo(string(manifest_read))
	case types.Unknown:
		return []string{}
	}

	return []string{}
}

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParseNPM(data []byte) []string {
	dependencies := []string{}

	// parsing package.json
	if !json.Valid(data) {
		return dependencies
	}

	var parsed PackageJSON
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

func ParsePNPM(data []byte) []string {
	dependencies := []string{}

	// parsing package.json
	if !json.Valid(data) {
		return dependencies
	}

	var parsed PackageJSON
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

func ParseYarn(data []byte) []string {
	dependencies := []string{}

	return dependencies
}

type BunLock struct {
	LockFileVersion string `json:"lockfileVersion"`
	ConfigVersion   string `json:"configVersion"`

	Workspaces map[string]struct {
		Name            string            `json:"name"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	} `json:"workspaces"`
}

func ParseBun(data []byte) []string {
	dependencies := []string{}

	// parsing package.json
	if !json.Valid(data) {
		return dependencies
	}

	var parsed BunLock
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range parsed.Workspaces[""].Dependencies {
		d := fmt.Sprintf("%s@%s", key, value)
		dependencies = append(dependencies, d)
	}

	return dependencies
}

func ParseGo(data []byte) []string {
	dependencies := []string{}

	parsed, err := modfile.Parse("", data, nil)
	if err != nil {
		fmt.Println(err)
	}

	for _, req := range parsed.Require {
		dependencies = append(dependencies, req.Mod.String())
	}

	return dependencies
}

type Dependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Project struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

func ParseJava(data []byte) []string {
	dependencies := []string{}

	var project Project
	err := xml.Unmarshal(data, &project)
	if err != nil {
		fmt.Printf("%s\n", err)
		return dependencies
	}

	for _, dependency := range project.Dependencies {
		d := fmt.Sprintf("%s@%s", dependency.ArtifactID, dependency.Version)
		dependencies = append(dependencies, d)
	}

	return dependencies
}

func ParsePython(data string) []string {
	dependencies := []string{}

	split_data := strings.Split(data, "\n")
	for _, d := range split_data {
		// changes == to @, for more consistency
		dependencies = append(dependencies, strings.Replace(d, "==", "@", 1))
	}

	return dependencies
}

type Cargo struct {
	Package struct {
		Name    string
		Version string
		Edition string
	}

	Dependencies map[string]interface{}
}

func ParseCargo(data string) []string {
	dependencies := []string{}

	var cargo Cargo
	_, err := toml.Decode(data, &cargo)

	if err != nil {
		fmt.Printf("%s\n", err)
		return dependencies
	}

	for name, version := range cargo.Dependencies {
		d := fmt.Sprintf("%s@%s", name, version)
		dependencies = append(dependencies, d)
	}

	return dependencies
}
