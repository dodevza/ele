package fileio

import (
	"strings"
)

// AppendPathChar ...
func AppendPathChar(path string) string {
	if len(path) == 0 {
		return ""
	}

	if path == "/" {
		return path
	}

	return path + "/"
}

// PathFromRoot ...
func PathFromRoot(root string, path string) string {
	result := strings.TrimLeft(path, root)
	result = strings.TrimLeft(result, "/")
	return result
	//removedRoot:= strings.Replace(path, root, "", 1)
}
