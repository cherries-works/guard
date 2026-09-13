package lock

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

type Cargo struct {
	Packages []struct {
		Name         string   `toml:"name"`
		Version      string   `toml:"version"`
		Dependencies []string `toml:"dependencies"`
	} `toml:"package"`
}

func ParseCargoL(data []byte) []string {
	dependencies := []string{}

	var lock Cargo
	if _, err := toml.Decode(string(data), &lock); err != nil {
		return dependencies
	}

	for _, pkg := range lock.Packages {
		parent := fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)

		for _, dependency := range pkg.Dependencies {
			dep := strings.Fields(dependency)

			if len(dep) == 0 {
				continue
			}

			if len(dep) == 1 {
				fmt.Printf("%s -> %s\n", parent, dep[0])
			} else {
				fmt.Printf("%s -> %s@%s\n", parent, dep[0], dep[1])
			}
		}
	}

	return dependencies
}
