package lock

import "strings"

func ParsePyPIL(data string) []string {
	dependencies := []string{}

	split_data := strings.Split(data, "\n")
	for _, d := range split_data {
		// changes == to @, for more consistency
		dependencies = append(dependencies, strings.Replace(d, "==", "@", 1))
	}

	return dependencies
}
