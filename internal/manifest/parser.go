package manifest

import (
	"fmt"
	"os"

	"github.com/cherries-works/guard/internal/types"
)

func ParseManifest(ecosystem types.Ecosystem, path string) []string {
	manifest_read, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	switch ecosystem {
	case types.NPM:
		return ParsePackageJSONM(manifest_read)
	case types.PNPM:
		return ParsePackageJSONM(manifest_read)
	case types.Yarn:
		return ParsePackageJSONM(manifest_read)
	case types.Bun:
		return ParseBunM(manifest_read)
	case types.Go:
		return ParseGoM(manifest_read)
	case types.Java:
		return ParseMavenM(manifest_read)
	case types.Python:
		return ParsePyPIM(string(manifest_read))
	case types.Cargo:
		return ParseCargoM(string(manifest_read))
	case types.Unknown:
		return []string{}
	}

	return []string{}
}
