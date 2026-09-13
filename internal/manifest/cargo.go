package manifest

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type CargoM struct {
	Package struct {
		Name    string
		Version string
		Edition string
	}

	Dependencies map[string]interface{}
}

func ParseCargoM(data string) []string {
	dependencies := []string{}

	var cargo CargoM
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
