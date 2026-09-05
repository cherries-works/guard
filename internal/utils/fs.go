package utils

import (
	"fmt"
	"log"
	"os"
)

var IgnoredEntries = []string{
	"node_modules",
	".git",
	".github",
}

func GetFiles(pwd string) []string {
	stored_files := []string{}

	files, err := os.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// not a proper way of doing things
		// but whatever
		if Includes(IgnoredEntries, file.Name()) {
			continue
		}

		full_filename := fmt.Sprintf("%s/%s", pwd, file.Name())
		stored_files = append(stored_files, full_filename)

		if file.IsDir() {
			get_files := GetFiles(full_filename)
			for _, f := range get_files {
				stored_files = append(stored_files, f)
			}
			continue
		}
	}

	return stored_files
}

func FindFile(pwd string, filename string) string {
	files, err := os.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		full_filename := fmt.Sprintf("%s/%s", pwd, file.Name())
		if file.Name() == filename {
			return full_filename
		}

		if file.IsDir() {
			check_file := FindFile(full_filename, filename)
			if check_file != "" {
				return check_file
			}
			continue
		}
	}

	return ""
}
