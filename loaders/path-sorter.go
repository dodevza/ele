package loaders

// TODO: Not particulary happy with this implementation
// Double sorting on the files by name (Ignore case) and then to so again  for before and after keywords
// Would have prefered to use regex patterns to allow more customization

import (
	"ele/config"
	"fmt"
	"sort"
	"strings"
)

// PathSorter ...
type PathSorter struct {
	before []string
	after  []string
}

type pathSorterNode struct {
}

// Sort ...
func (srt *PathSorter) Sort(files ...string) []string {
	sortedByName := srt.SortIgnoreCase(files...)

	pathMap := make(map[string]string, 0)
	lowerPaths := make([]string, 0)
	for uniqueIndex, p := range sortedByName {
		key := srt.createKey(p, uniqueIndex)
		pathMap[key] = p
		lowerPaths = append(lowerPaths, key)
	}

	sort.Strings(lowerPaths)

	result := make([]string, 0)

	for _, key := range lowerPaths {
		original := pathMap[key]

		result = append(result, original)
	}

	return result
}

// SortIgnoreCase ...
func (srt *PathSorter) SortIgnoreCase(files ...string) []string {
	pathMap := make(map[string]string, 0)
	lowerPaths := make([]string, 0)
	for _, p := range files {
		key := strings.ToLower(p)
		pathMap[key] = p
		lowerPaths = append(lowerPaths, key)
	}

	sort.Strings(lowerPaths)

	result := make([]string, 0)

	for _, key := range lowerPaths {
		original := pathMap[key]

		result = append(result, original)
	}

	return result
}

func (srt *PathSorter) createKey(path string, unqCounter int) string {
	lowerKey := strings.ToLower(path)
	for _, a := range srt.after {
		index := strings.Index(lowerKey, a)
		if index == 0 {
			return fmt.Sprintf("~~~~~%d", unqCounter)
		}
		if index > -1 {
			return fmt.Sprintf("%s~~~%d", lowerKey[:index-1], unqCounter)
		}
	}

	for _, a := range srt.before {
		index := strings.Index(lowerKey, a)
		if index == 0 {

			return fmt.Sprintf("     %d", unqCounter)
		}
		if index > -1 {
			return fmt.Sprintf("%s   %d", lowerKey[:index-1], unqCounter)
		}
	}

	return lowerKey
}

// NewPathSorter ...
func NewPathSorter(hooks *config.HooksConfig) *PathSorter {

	return &PathSorter{before: hooks.Before, after: hooks.After}
}

func newPathSorter() *PathSorter {
	return &PathSorter{before: []string{"before."}, after: []string{"after."}}
}
