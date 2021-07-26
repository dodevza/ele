package environments

import (
	"fmt"
	"strings"
)

// Activate ...
func (collection *EnvironmentState) Activate(name string) error {
	collection.lazyLoad()
	key := strings.ToLower(name)
	_, foundFile := collection.Map[key]
	if foundFile {
		collection.active = name
		return collection.SaveActiveState()
	}

	_, foundGroup := collection.Groups[key]

	if foundGroup {
		collection.active = name
		return collection.SaveActiveState()
	}

	return fmt.Errorf("Did not find environment %s to activate", name)
}

// Deactivate ...
func (collection *EnvironmentState) Deactivate() error {
	collection.active = ""
	return collection.SaveActiveState()
}
