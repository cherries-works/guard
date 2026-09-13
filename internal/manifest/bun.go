package manifest

import (
	"encoding/json"
	"fmt"
)

type BunLockM struct {
	LockFileVersion string `json:"lockfileVersion"`
	ConfigVersion   string `json:"configVersion"`

	Workspaces map[string]struct {
		Name            string            `json:"name"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	} `json:"workspaces"`
}

func ParseBunM(data []byte) []string {
	dependencies := []string{}

	// parsing package.json
	if !json.Valid(data) {
		return dependencies
	}

	var parsed BunLockM
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
