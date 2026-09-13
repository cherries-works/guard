package lock

import "strings"

func ParseGoL(data []byte) []string {
	dependencies := []string{}

	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)

		if len(parts) < 2 || strings.HasSuffix(parts[1], "/go.mod") {
			continue
		}

		// fmt.Printf("%s@%s\n", parts[0], parts[1])
	}

	return dependencies
}
