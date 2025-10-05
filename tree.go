package palantir

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Tree represents a generic tree structure
type Tree[T any] interface {
	Root() *Node[T]
	Insert(path []string, data T) error
	Find(path []string) (*Node[T], bool)
	Traverse(visitor func(*Node[T]) bool)
	Sort(comparator func(a, b *Node[T]) bool)
	Size() int
}

// Node represents a node in the generic tree
type Node[T any] struct {
	Data     T
	Name     string
	Children []*Node[T]
	Parent   *Node[T]
}

// NewNode creates a new tree node
func NewNode[T any](name string, data T) *Node[T] {
	return &Node[T]{
		Data:     data,
		Name:     name,
		Children: make([]*Node[T], 0),
		Parent:   nil,
	}
}

// AddChild adds a child node to this node
func (n *Node[T]) AddChild(child *Node[T]) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// FindChild finds a child node by name
func (n *Node[T]) FindChild(name string) (*Node[T], bool) {
	for _, child := range n.Children {
		if child.Name == name {
			return child, true
		}
	}
	return nil, false
}

// IsLeaf returns true if the node has no children
func (n *Node[T]) IsLeaf() bool {
	return len(n.Children) == 0
}

// Depth returns the depth of the node in the tree
func (n *Node[T]) Depth() int {
	depth := 0
	current := n.Parent
	for current != nil {
		depth++
		current = current.Parent
	}
	return depth
}

// Path returns the path from root to this node
func (n *Node[T]) Path() []string {
	if n.Parent == nil {
		return []string{n.Name}
	}
	return append(n.Parent.Path(), n.Name)
}

// genericTree implements the Tree interface
type genericTree[T any] struct {
	root *Node[T]
}

// NewTree creates a new generic tree
func NewTree[T any](rootData T, rootName string) Tree[T] {
	return &genericTree[T]{
		root: NewNode(rootName, rootData),
	}
}

// Root returns the root node
func (t *genericTree[T]) Root() *Node[T] {
	return t.root
}

// Insert inserts a node at the specified path
func (t *genericTree[T]) Insert(path []string, data T) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	current := t.root
	for _, name := range path[:len(path)-1] {
		child, found := current.FindChild(name)
		if !found {
			// Create intermediate node with zero value
			var zero T
			child = NewNode(name, zero)
			current.AddChild(child)
		}
		current = child
	}

	// Add the final node
	finalName := path[len(path)-1]
	finalNode := NewNode(finalName, data)
	current.AddChild(finalNode)
	return nil
}

// Find finds a node at the specified path
func (t *genericTree[T]) Find(path []string) (*Node[T], bool) {
	current := t.root
	for _, name := range path {
		child, found := current.FindChild(name)
		if !found {
			return nil, false
		}
		current = child
	}
	return current, true
}

// Traverse traverses the tree and calls visitor for each node
func (t *genericTree[T]) Traverse(visitor func(*Node[T]) bool) {
	t.traverseNode(t.root, visitor)
}

// traverseNode recursively traverses a node and its children
func (t *genericTree[T]) traverseNode(node *Node[T], visitor func(*Node[T]) bool) {
	if !visitor(node) {
		return
	}
	for _, child := range node.Children {
		t.traverseNode(child, visitor)
	}
}

// Sort sorts the tree using the provided comparator
func (t *genericTree[T]) Sort(comparator func(a, b *Node[T]) bool) {
	t.sortNode(t.root, comparator)
}

// sortNode recursively sorts a node and its children
func (t *genericTree[T]) sortNode(node *Node[T], comparator func(a, b *Node[T]) bool) {
	if len(node.Children) == 0 {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		return comparator(node.Children[i], node.Children[j])
	})

	for _, child := range node.Children {
		t.sortNode(child, comparator)
	}
}

// Size returns the total number of nodes in the tree
func (t *genericTree[T]) Size() int {
	count := 0
	t.Traverse(func(node *Node[T]) bool {
		count++
		return true
	})
	return count
}

// FileNode represents a file or directory in the filesystem tree
type FileNode struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime int64
}

// ShowHierarchy displays a tree structure of files/directories using the generic tree system
func ShowHierarchy(basePath, targetDir string) (error, bool) {
	// Get root directory info
	rootInfo, err := os.Stat(basePath)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err), false
	}

	// Create root FileNode data
	rootData := FileNode{
		Name:    rootInfo.Name(),
		Path:    basePath,
		IsDir:   rootInfo.IsDir(),
		Size:    rootInfo.Size(),
		ModTime: rootInfo.ModTime().Unix(),
	}

	// Create generic tree with FileNode data
	tree := NewTree(rootData, rootInfo.Name())

	// Walk filesystem and build tree directly
	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == basePath {
			return nil // Skip root directory itself
		}

		// Skip hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path and split into components
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		parts := strings.Split(relPath, string(filepath.Separator))

		// Create FileNode data
		fileData := FileNode{
			Name:    info.Name(),
			Path:    path,
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		}

		// Insert into generic tree
		return tree.Insert(parts, fileData)
	})

	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err), false
	}

	// Check if tree has only one node and it's not a directory
	root := tree.Root()
	if len(root.Children) == 1 && !root.Children[0].Data.IsDir {
		return nil, false // No hierarchy needed
	}

	// Sort tree (directories first, then alphabetically)
	tree.Sort(func(a, b *Node[FileNode]) bool {
		if a.Data.IsDir != b.Data.IsDir {
			return a.Data.IsDir
		}
		return a.Data.Name < b.Data.Name
	})

	// Render the tree directly
	printTree(tree.Root(), "", true, true)

	return nil, true
}

// printTree recursively prints a tree node with ASCII art and colors
func printTree(node *Node[FileNode], prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		// Choose the appropriate tree character
		var treeChar string
		if isLast {
			treeChar = "└── "
		} else {
			treeChar = "├── "
		}

		// Style the node name
		styledName := styleFileNode(node)

		// Print the current node
		fmt.Printf("%s%s%s\n", prefix, treeChar, styledName)
	}

	// Print children
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1

			// Calculate prefix for child
			var childPrefix string
			if isRoot {
				childPrefix = ""
			} else {
				if isLast {
					childPrefix = prefix + "    "
				} else {
					childPrefix = prefix + "│   "
				}
			}

			printTree(child, childPrefix, isChildLast, false)
		}
	}
}

// styleFileNode styles a filesystem node based on OutputConfig
func styleFileNode(node *Node[FileNode]) string {
	outputConfig := GetGlobalOutputHandler().(*outputHandler).config

	if !outputConfig.UseColors {
		return node.Data.Name
	}

	// Get file node data
	fileNode := node.Data

	if fileNode.IsDir {
		return fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, fileNode.Name, ColorReset)
	}

	// Color files based on extension
	ext := strings.ToLower(filepath.Ext(fileNode.Name))
	switch ext {
	case ".json", ".yaml", ".yml", ".toml":
		return fmt.Sprintf("%s%s%s", ColorGreen, fileNode.Name, ColorReset)
	case ".md", ".txt", ".log":
		return fmt.Sprintf("%s%s%s", ColorCyan, fileNode.Name, ColorReset)
	case ".sh", ".zsh", ".bash":
		return fmt.Sprintf("%s%s%s", ColorYellow, fileNode.Name, ColorReset)
	case ".go":
		return fmt.Sprintf("%s%s%s", ColorPurple, fileNode.Name, ColorReset)
	default:
		return fileNode.Name
	}
}
