package analyzer

import (
	"strings"

	"github.com/cherries-works/guard/internal/types"
	"github.com/cherries-works/guard/internal/utils"
)

type Ecosystem struct {
	Ecosystem    types.Ecosystem
	Manifests    []string
	Locks        []string
	Dependencies []Dependency
}

type Project struct {
	Root       string
	Ecosystems []Ecosystem
}

var ManifestFiles = []string{
	"package.json",
	"go.mod",
	"Cargo.toml",
	"requirements.txt",
	"pyproject.toml",
	"pom.xml",
	"build.gradle",
}

var ManifestFilesMapped = []types.Ecosystem{
	types.NPM,
	types.Go,
	types.Cargo,
	types.Python,
	types.Python,
	types.Java,
	types.Java,
}

var LockFiles = []string{
	"package-lock.json",
	"pnpm-lock.yaml",
	"yarn.lock",
	"bun.lockb",
	"bun.lock",
	"go.sum",
	"Cargo.lock",
}

var LockFilesMapped = []types.Ecosystem{
	types.NPM,
	types.PNPM,
	types.Yarn,
	types.Bun,
	types.Bun,
	types.Go,
	types.Cargo,
}

var LanguageExtensions = []string{
	"js",
	"cjs",
	"mjs",
	"ts",
	"go",
	"rs",
	"py",
	"java",
}

var LanguageExtensionsMapped = []types.Ecosystem{
	types.NPM,
	types.NPM,
	types.NPM,
	types.NPM,
	types.Go,
	types.Cargo,
	types.Python,
	types.Java,
}

func GetEcoSystems(pwd string) []Ecosystem {
	ecosystems := []Ecosystem{}
	files := utils.GetFiles(pwd)
	for _, filename := range files {
		for manifest_index, manifest_filename := range ManifestFiles {
			if strings.HasSuffix(filename, manifest_filename) {
				checked := false
				for ecosystem_index, _ecosystem := range ecosystems {
					if _ecosystem.Ecosystem == ManifestFilesMapped[manifest_index] {
						ecosystems[ecosystem_index].Manifests = append(ecosystems[ecosystem_index].Manifests, filename)
						checked = true
						break
					}
				}

				if checked {
					continue
				}

				ecosystem := Ecosystem{
					Ecosystem: ManifestFilesMapped[manifest_index],
					Manifests: []string{filename},
				}
				ecosystems = append(ecosystems, ecosystem)
			}
		}

		for lock_index, lock_filename := range LockFiles {
			if strings.HasSuffix(filename, lock_filename) {
				checked := false
				for ecosystem_index, _ecosystem := range ecosystems {
					if _ecosystem.Ecosystem == LockFilesMapped[lock_index] {
						ecosystems[ecosystem_index].Locks = append(ecosystems[ecosystem_index].Locks, filename)
						checked = true
						break
					}
				}

				if checked {
					continue
				}

				ecosystem := Ecosystem{
					Ecosystem: LockFilesMapped[lock_index],
					Locks:     []string{filename},
				}
				ecosystems = append(ecosystems, ecosystem)
			}
		}
	}

	return ecosystems
}

func GetManifests(pwd string) []string {
	manifests := []string{}
	files := utils.GetFiles(pwd)
	for _, filename := range files {
		for _, manifest_filename := range ManifestFiles {
			if strings.HasSuffix(filename, manifest_filename) {
				manifests = append(manifests, filename)
			}
		}
	}

	return manifests
}

func GetLocks(pwd string) []string {
	locks := []string{}
	files := utils.GetFiles(pwd)
	for _, filename := range files {
		for _, lock_filename := range LockFiles {
			if strings.HasSuffix(filename, lock_filename) {
				locks = append(locks, filename)
			}
		}
	}

	return locks
}

func GetProject(pwd string) Project {
	ecosystems := GetEcoSystems(pwd)

	project := Project{
		Root:       pwd,
		Ecosystems: ecosystems,
	}

	return project
}
