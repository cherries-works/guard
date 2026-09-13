package lock

import (
	"fmt"

	"github.com/iseki0/go-yarnlock"
)

type YarnPackage struct {
	Name         string
	Version      string
	Dependencies map[string]string
}

func ParseYarnL(data []byte) []string {
	dependencies := []string{}

	parsed, err := yarnlock.ParseLockFileData(data)
	if err != nil {
		return dependencies
	}

	for key, pkg := range parsed {
		fmt.Println("KEY:", key)
		fmt.Println("NAME:", pkg.Name)
		fmt.Println("VERSION:", pkg.Version)

		for name, version := range pkg.Dependencies {
			fmt.Println("  ->", name, version)
		}
	}

	return dependencies
}
