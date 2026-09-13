package lock

import (
	"encoding/xml"
	"fmt"
)

type Dependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Project struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

func ParseMavenL(data []byte) []string {
	dependencies := []string{}

	var project Project
	err := xml.Unmarshal(data, &project)
	if err != nil {
		fmt.Printf("%s\n", err)
		return dependencies
	}

	for _, dependency := range project.Dependencies {
		fmt.Printf("%s@%s\n", dependency.ArtifactID, dependency.Version)
		// dependencies = append(dependencies, d)
	}

	return dependencies
}
