package manifest

import (
	"fmt"

	"golang.org/x/mod/modfile"
)

func ParseGoM(data []byte) []string {
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
