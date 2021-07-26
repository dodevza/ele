package fileio

import (
	"log"
	"os"
)

// CreateDirectoryIfNotExists ...
func CreateDirectoryIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Fatalf("Could not create directory: %s, %s", path, err)
		}
	}
}

// DeleteFileIfExists ...
func DeleteFileIfExists(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return
	}

	os.Remove(filename)
}

// DeleteDirectoryIfExists ...
func DeleteDirectoryIfExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return
	}

	os.Remove(folder)
}
