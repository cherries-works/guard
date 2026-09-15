package analyzer

import (
	"fmt"
	"strings"

	"github.com/cherries-works/guard/internal/parser"
	"github.com/cherries-works/guard/internal/types"
)

type Finding struct {
	DependencyID string
	ID           string
	Severity     string
	Summary      string
	FixedVersion string
}

type Dependency struct {
	ID            string
	Name          string
	Semver        string
	Version       string
	LatestVersion string
	Vulnerable    bool
	Outdated      bool
}

func GetDependencies(ecosystem types.Ecosystem, manifest string) ([]Dependency, []Finding) {
	dependencies := []Dependency{}
	findings := []Finding{}

	string_parsed_dependencies := parser.ParseManifest(ecosystem, manifest)
	for _, string_dependency := range string_parsed_dependencies {
		x := strings.Split(string_dependency, "@")
		name := x[0]
		semver := x[0]
		version := strings.TrimLeft(x[1], "v^=")
		latest_version := strings.TrimLeft(GetLatest(ecosystem, name), "v^=")
		vulns := GetVulns(types.EcosystemMapped[ecosystem], name, version)

		id := fmt.Sprintf("%d-%s-%s", ecosystem, name, version)
		for _, vuln := range vulns {
			fixed_version := ""
			for _, affected := range vuln.Affected {
				for _, affectedRange := range affected.Ranges {
					for _, event := range affectedRange.Events {
						if event.Fixed != "" {
							fixed_version = event.Fixed
							break
						}
					}
					if fixed_version != "" {
						break
					}
				}
				if fixed_version != "" {
					break
				}
			}

			finding := Finding{
				ID:           vuln.ID,
				DependencyID: id,
				Severity:     vuln.DatabaseSpecifics.Severity,
				Summary:      vuln.Summary,
				FixedVersion: strings.TrimLeft(fixed_version, "v^="),
			}
			findings = append(findings, finding)
		}

		version = strings.TrimLeft(version, "v^=")
		latest_version = strings.TrimLeft(latest_version, "v^=")

		dependency := Dependency{
			ID:            id,
			Name:          name,
			Semver:        semver,
			Version:       version,
			LatestVersion: latest_version,
			Vulnerable:    len(vulns) > 0,
			Outdated:      version != latest_version,
		}

		dependencies = append(dependencies, dependency)
	}

	return dependencies, findings
}
