package analyzer

import (
	"encoding/json"
	"net/http"

	"github.com/cherries-works/guard/internal/types"
)

func GetLatest(ecosystem types.Ecosystem, name string) string {
	latest := ""
	switch ecosystem {
	case types.Go:
		latest = GetLatestGoVersion(name)
	case types.NPM:
		latest = GetLatestNPMVersion(name)
	case types.PNPM:
		latest = GetLatestNPMVersion(name)
	case types.Yarn:
		latest = GetLatestNPMVersion(name)
	case types.Bun:
		latest = GetLatestNPMVersion(name)
	case types.Cargo:
		latest = GetLatestCargoVersion(name)
	case types.Python:
		latest = GetLatestPyPiVersion(name)
	}

	return latest
}

type GoLatestResponse struct {
	Version string `json:"Version"`
	Time    string `json:"Time"`
}

func GetLatestGoVersion(name string) string {
	url := "https://proxy.golang.org/" + name + "/@latest"

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result GoLatestResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Version
}

type NPMResponse struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
}

func GetLatestNPMVersion(name string) string {
	url := "https://registry.npmjs.org/" + name

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result NPMResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.DistTags.Latest
}

type PyPiResponse struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
}

func GetLatestPyPiVersion(name string) string {
	url := "https://pypi.org/pypi/" + name + "/json"

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result PyPiResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Info.Version
}

type CargoResponse struct {
	Crate struct {
		NewestVersion string `json:"newest_version"`
	} `json:"crate"`
}

func GetLatestCargoVersion(name string) string {
	url := "https://crates.io/api/v1/crates/" + name

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result CargoResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Crate.NewestVersion
}
