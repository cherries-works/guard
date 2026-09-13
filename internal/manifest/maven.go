package manifest

import (
	"encoding/xml"
	"fmt"
)

type DependencyM struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type ProjectM struct {
	Dependencies []DependencyM `xml:"dependencies>dependency"`
}

func ParseMavenM(data []byte) []string {
	dependencies := []string{}

	var project ProjectM
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
