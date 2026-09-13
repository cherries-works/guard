package lock

import (
	"fmt"
	"os"

	"github.com/cherries-works/guard/internal/types"
)

func ParseLock(ecosystem types.Ecosystem, path string) []string {
	lock_read, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	switch ecosystem {
	case types.NPM:
		return ParseNPML(lock_read)
	case types.PNPM:
		return ParsePNPML(lock_read)
	case types.Yarn:
		return ParseYarnL(lock_read)
	case types.Bun:
		return ParseBunL(lock_read)
	case types.Go:
		return ParseGoL(lock_read)
	case types.Java:
		return ParseMavenL(lock_read)
	case types.Python:
		return ParsePyPIL(string(lock_read))
	case types.Cargo:
		return ParseCargoL(lock_read)
	case types.Unknown:
		return []string{}
	}

	return []string{}
}
