package lock

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type PNPMLock struct {
	Importers map[string]PNPMImporter `yaml:"importers"`
	Packages  map[string]PNPMPackage  `yaml:"packages"`
}

type PNPMImporter struct {
	Dependencies    map[string]PNPMDependency `yaml:"dependencies"`
	DevDependencies map[string]PNPMDependency `yaml:"devDependencies"`
}

type PNPMDependency struct {
	Specifier string `yaml:"specifier"`
	Version   string `yaml:"version"`
}

type PNPMPackage struct {
	Dependencies map[string]string `yaml:"dependencies"`
}

func ParsePNPML(data []byte) []string {
	dependencies := []string{}

	var parsed PNPMLock
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return dependencies
	}

	for name, pkg := range parsed.Packages {
		fmt.Println(name)

		for dependency, version := range pkg.Dependencies {
			fmt.Printf("-> %s@%s", dependency, version)
		}
	}

	return dependencies
}
