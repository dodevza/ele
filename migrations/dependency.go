package migrations

import (
	"sort"
	"strings"
)

// DependencyNode ...
type DependencyNode struct {
	level      int
	name       string
	moduleKeys []string
	modules    map[string]*DependencyNode
	migrations MigrationCollection

	parentNode *DependencyNode
}

func (node *DependencyNode) setLevel(level int) int {
	node.level = level
	maxLevel := level
	for _, module := range node.modules {
		moduleDepth := module.setLevel(level + 1)
		if moduleDepth > maxLevel {
			maxLevel = moduleDepth
		}
	}

	return maxLevel
}

func newDependencyNode(name string) DependencyNode {
	modules := make(map[string]*DependencyNode)
	moduleKeys := make([]string, 0)
	return DependencyNode{name: name, modules: modules, moduleKeys: moduleKeys}
}

// DependencyTree ...
type DependencyTree struct {
	firstNode  *DependencyNode
	moduleKeys []string
	modules    map[string]*DependencyNode
	maxLevel   int
}

// Add ...
func (tree *DependencyTree) Add(migration Migration) {
	node := tree.addModule(migration.Path)

	node.migrations = append(node.migrations, &migration)

	// Add dependencies
	for _, requirement := range migration.Requires {
		rNode := tree.addModule(requirement)
		// tree.firstNode.addDependency(tree, rNode)
		rNode.addDependency(tree, node)
	}
	// tree.firstNode.addDependency(tree, node)

	parentModules := getParentModules(migration.Path)
	if len(parentModules) == 0 {
		return
	}

	prevNode := tree.addModule(parentModules[0])
	parentIndex := 1
	for parentIndex < len(parentModules) {
		rNode := tree.addModule(parentModules[parentIndex])
		rNode.addDependency(tree, prevNode)
		prevNode = rNode
		parentIndex++
	}
	node.addDependency(tree, prevNode)
}

func (tree *DependencyTree) addModule(name string) *DependencyNode {
	innernode, ok := tree.modules[name]
	if ok == false {
		nodeValue := newDependencyNode(name)
		innernode = &nodeValue
		tree.modules[name] = innernode
		tree.updateKeys()
	}

	return innernode
}

func getParentModules(path string) []string {
	list := make([]string, 0)
	lastIndex := strings.LastIndex(path, "/")
	max := 10
	for lastIndex > 0 && max > 0 {
		path = path[0:lastIndex]
		list = append([]string{path}, list...)
		lastIndex = strings.LastIndex(path, "/")
		max--
	}
	return list
}

func (node *DependencyNode) addDependency(tree *DependencyTree, dep *DependencyNode) {
	if dep.level > node.level {
		return
	}
	maxLevel := dep.setLevel(node.level + 1)

	if tree.maxLevel < maxLevel {
		tree.maxLevel = maxLevel
	}
	// if dep.parentNode != nil {
	// 	dep.parentNode.removeDependency(dep)
	// 	dep.parentNode.updateKeys()
	// }
	dep.parentNode = node

	_, ok := node.modules[dep.name]
	if ok == false {
		node.modules[dep.name] = dep
		node.updateKeys()
	}

}

func (node *DependencyNode) removeDependency(dep *DependencyNode) {
	delete(node.modules, dep.name)
}

func (node *DependencyNode) updateKeys() {
	keys := make([]string, 0)
	for m := range node.modules {
		keys = append(keys, m)
	}
	sort.Strings(keys)
	node.moduleKeys = keys
}

func (tree *DependencyTree) updateKeys() {
	keys := make([]string, 0)
	for m := range tree.modules {
		keys = append(keys, m)
	}
	sort.Strings(keys)
	tree.moduleKeys = keys
}

// InOrder ...
func (tree *DependencyTree) InOrder() MigrationCollection {
	list := make(MigrationCollection, 0)
	level := 0
	for level <= tree.maxLevel {
		for _, key := range tree.moduleKeys {
			module := tree.modules[key]
			if module.level == level {
				list = append(list, module.migrations...)
			}
		}
		level++
	}
	return list
}

func (node *DependencyNode) inOrder() MigrationCollection {
	list := make(MigrationCollection, 0)

	list = append(list, node.migrations...)

	for _, key := range node.moduleKeys {
		module := node.modules[key]
		if module.level == node.level+1 {
			list = append(list, module.migrations...)
		}
	}

	for _, key := range node.moduleKeys {
		module := node.modules[key]
		if module.level == node.level+1 {
			list = append(list, module.inOrderDependencies()...)
		}
	}

	return list
}

func (node *DependencyNode) inOrderDependencies() MigrationCollection {
	list := make(MigrationCollection, 0)
	for _, module := range node.modules {
		if module.level == node.level+1 {
			list = append(list, module.inOrder()...)
		}
	}

	return list
}

// EmptyDependencyTree ...
func EmptyDependencyTree() DependencyTree {
	firstNode := newDependencyNode("")
	modules := make(map[string]*DependencyNode)
	moduleKeys := make([]string, 0)
	return DependencyTree{firstNode: &firstNode, modules: modules, moduleKeys: moduleKeys}
}
