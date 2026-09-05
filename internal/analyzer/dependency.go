package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cherries-works/guard/internal/parser"
	"github.com/cherries-works/guard/internal/types"
)

type DependencyStatus struct {
	CurrentVersion string
	LatestVersion  string
	Outdated       bool
	Vulnerable     bool
	Vulns          []OSVVulnerability
}

type Dependency struct {
	Name    string
	Version string
	Status  DependencyStatus
}

func GetDependencies(ecosystem types.Ecosystem, manifest string) []Dependency {
	dependencies := []Dependency{}

	string_parsed_dependencies := parser.ParseManifest(ecosystem, manifest)
	for _, string_dependency := range string_parsed_dependencies {
		x := strings.Split(string_dependency, "@")
		name := x[0]
		version := x[1]
		latest_version := GetLatest(ecosystem, name)
		vulns := GetVulns(types.EcosystemMapped[ecosystem], name, version)
		status := DependencyStatus{
			CurrentVersion: version,
			LatestVersion:  latest_version,
			Outdated:       latest_version != version,
			Vulnerable:     len(vulns) > 0,
			Vulns:          vulns,
		}

		dependency := Dependency{
			Name:    name,
			Version: version,
			Status:  status,
		}

		dependencies = append(dependencies, dependency)
	}

	return dependencies
}

type OSVResponse struct {
	Vulns []OSVVulnerability `json:"vulns"`
}

type OSVVulnerability struct {
	ID                string               `json:"id"`
	Summary           string               `json:"summary"`
	Details           string               `json:"details"`
	Modified          string               `json:"modified"`
	Published         string               `json:"published"`
	Affected          []OSVAffected        `json:"affected"`
	DatabaseSpecifics OSVDatabaseSpecifics `json:"database_specifics"`
}

type OSVAffected struct {
	Package OSVPackage `json:"package"`
	Ranges  []OSVRange `json:"ranges"`
}

type OSVPackage struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
}

type OSVDatabaseSpecifics struct {
	Severity         string `json:"severity,omitempty"`
	GithubReviewedAt string `json:"github_reviewed_at,omitempty"`
	GithubReviewed   bool   `json:"github_reviewed,omitempty"`
}

type OSVRange struct {
	Type   string     `json:"type"`
	Events []OSVEvent `json:"events"`
}

type OSVEvent struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}

type OSVQuery struct {
	Package OSVPackage `json:"package"`
	Version string     `json:"version"`
}

func GetVulns(ecosystem string, name string, version string) []OSVVulnerability {
	query := OSVQuery{
		Package: OSVPackage{
			Name:      name,
			Ecosystem: ecosystem,
		},
		Version: version,
	}

	body, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", "https://api.osv.dev/v1/query", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	response_body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var osv_reponse OSVResponse
	err = json.Unmarshal([]byte(response_body), &osv_reponse)
	if err != nil {
		fmt.Println(err)
		return []OSVVulnerability{}
	}

	for _, vuln := range osv_reponse.Vulns {
		for _, affected := range vuln.Affected {
			for _, r := range affected.Ranges {
				for _, event := range r.Events {
					_ = event
				}
			}
		}
	}

	return osv_reponse.Vulns
}
