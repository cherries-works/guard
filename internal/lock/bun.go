package lock

import (
	"encoding/json"
	"fmt"
)

type BunLock struct {
	LockfileVersion int                     `json:"lockfileVersion"`
	Workspaces      map[string]BunWorkspace `json:"workspaces"`
	Packages        map[string]BunPackage   `json:"packages"`
}

type BunWorkspace struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type BunPackage struct {
	Resolution string
	Info       BunPackageInfo
}

type BunPackageInfo struct {
	Dependencies         map[string]string `json:"dependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
}

func ParseBunL(data []byte) []string {
	dependencies := []string{}

	var lock struct {
		Packages map[string][]json.RawMessage `json:"packages"`
	}

	if err := json.Unmarshal(data, &lock); err != nil {
		return dependencies
	}

	for _, pkg := range lock.Packages {
		if len(pkg) < 3 {
			continue
		}

		var resolved string
		if err := json.Unmarshal(pkg[0], &resolved); err != nil {
			continue
		}

		var info struct {
			Dependencies map[string]string `json:"dependencies"`
		}

		if err := json.Unmarshal(pkg[2], &info); err != nil {
			continue
		}

		for name, version := range info.Dependencies {
			fmt.Printf("%s -> %s@%s\n", resolved, name, version)
		}
	}

	return dependencies
}
