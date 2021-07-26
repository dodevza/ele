package migrations

import (
	"ele/constants"
	"strings"
)

// Version  ...
type Version struct {
	Name string

	next         *Version
	prev         *Version
	dependencies *DependencyTree
}

func newVersion(version string) Version {
	dependencies := EmptyDependencyTree()
	pointToDep := &dependencies
	return Version{Name: version, dependencies: pointToDep}
}

// VersionTree ...
type VersionTree struct {
	firstNode  *Version
	lastNode   *Version
	repeatable *Version
	tags       map[string]bool
}

// ToString ...
func (tree *VersionTree) ToString() string {
	if tree.firstNode == nil {
		return ""
	}

	node := tree.firstNode
	str := ""
	for node != nil {
		str += node.Name + ","
		node = node.next
	}

	if tree.repeatable != nil {
		str += constants.REPEATABLE + ","
	}

	return strings.TrimRight(str, ",")
}

// Add ...
func (tree *VersionTree) Add(migration Migration) {
	version := tree.addVersion(migration.Version)

	version.dependencies.Add(migration)
}

// ToArray ...
func (tree *VersionTree) ToArray() []*Version {
	result := make([]*Version, 0)
	node := tree.firstNode
	for node != nil {
		result = append(result, node)
		node = node.next
	}

	return result
}

// IsTag ..
func (tree *VersionTree) IsTag(name string) bool {
	return tree.tags[name]
}

// Forward Returns a list of migrations
func (tree *VersionTree) Forward(versionStart string, versionEnd string) MigrationCollection {
	collection := make(MigrationCollection, 0)
	node := tree.firstNode
	foundStart := false
	for node != nil {
		if versionEnd != "" && compareVersion(node.Name, versionEnd) > 0 {
			break
		}
		if versionStart == "" || compareVersion(node.Name, versionStart) >= 0 {
			foundStart = true
		}

		if foundStart {
			collection = append(collection, node.dependencies.InOrder()...)
		}

		if versionEnd != "" && compareVersion(node.Name, versionEnd) >= 0 {
			break
		}

		node = node.next
	}
	if tree.repeatable != nil {
		collection = append(collection, tree.repeatable.dependencies.InOrder()...)
	}

	return collection
}

func (tree *VersionTree) addVersion(version string) *Version {
	versionToUpper := strings.ToUpper(version)

	if versionToUpper == constants.REPEATABLE {
		repeatable := tree.repeatable
		if repeatable == nil {
			repeatableValue := newVersion(constants.REPEATABLE)
			repeatable = &repeatableValue
			tree.repeatable = &repeatableValue
		}
		return repeatable
	}
	if tree.firstNode == nil {
		v := newVersion(versionToUpper)
		tree.firstNode = &v
		tree.lastNode = &v
	}

	if tree.tags[versionToUpper] {
		return tree.addTagAtTheEnd(versionToUpper)
	}

	node := tree.firstNode
	for tree.tags[node.Name] == false && compareVersion(node.Name, versionToUpper) < 0 && node.next != nil {
		node = node.next
	}

	if node.Name == versionToUpper {
		return node
	}

	v := newVersion(versionToUpper)

	if tree.tags[node.Name] || compareVersion(node.Name, versionToUpper) > 0 {
		prev := node.prev

		node.prev = &v

		if prev != nil {
			v.prev = prev
			prev.next = &v
		} else {
			tree.firstNode = &v
		}
		v.next = node

	} else {
		next := node.next

		node.next = &v

		if next != nil {
			v.next = next
			next.prev = &v

		} else {
			tree.lastNode = &v

		}
		v.prev = node
	}
	return &v

}

func compareVersion(a string, b string) int {

	aPrefixCount := 0
	aChar := a[0]
	for aChar == '_' {
		aPrefixCount++
		aChar = a[aPrefixCount]
	}

	bPrefixCount := 0
	bChar := b[0]
	for bChar == '_' {
		bPrefixCount++
		bChar = b[bPrefixCount]
	}

	// Having a prefix should move the version to the top
	if aPrefixCount != 0 && bPrefixCount == 0 {
		return -1
	} else if bPrefixCount != 0 && aPrefixCount == 0 {
		return 1
	}

	// Having more prefixes should lower the version in the sequence
	// _init
	// __then
	// ___something_else

	if aPrefixCount > bPrefixCount {
		return 1
	} else if bPrefixCount > aPrefixCount {
		return -1
	}

	astring := a[aPrefixCount:]
	bstring := b[bPrefixCount:]
	if astring == bstring {
		return 0
	}

	if astring > bstring {
		return 1
	}

	return -1

}

func (tree *VersionTree) addTagAtTheEnd(version string) *Version {
	node := tree.lastNode
	for tree.tags[node.Name] && node.Name != version && node.prev != nil {
		node = node.prev
	}

	if node.Name == version {
		return node
	}

	// If its not found, append at the end of tree
	// If we want the tags to be in the sequence the user specified, we need to pre-populate the tags in that order

	v := newVersion(version)

	lastNode := tree.lastNode
	v.prev = lastNode
	lastNode.next = &v
	tree.lastNode = &v

	return &v
}

// EmptyVersionTreeWithTags ...
func EmptyVersionTreeWithTags(tags []string) VersionTree {
	tagMap := make(map[string]bool, 0)

	for _, tag := range tags {
		key := strings.ToUpper(tag)
		tagMap[key] = true
	}

	tree := VersionTree{tags: tagMap}

	// Prepopulate tags in order they were provided
	for _, tag := range tags {
		tree.addVersion(tag)
	}

	return tree
}

// EmptyVersionTree ...
func EmptyVersionTree() VersionTree {
	return EmptyVersionTreeWithTags([]string{})
}
