package palantir

import (
	"fmt"
	"io"
	"os"
	"sort"
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

// TreeBuilder defines the interface for building trees from different sources
type TreeBuilder[T any] interface {
	Build(source string) (Tree[T], error)
}

// TreeRenderer defines the interface for rendering trees in different formats
type TreeRenderer[T any] interface {
	Render(tree Tree[T], writer io.Writer) error
}

// NodeStyler defines the interface for customizing node appearance
type NodeStyler[T any] interface {
	Style(node *Node[T]) string
	GetTreeChar(node *Node[T], isLast bool) string
	GetPrefix(node *Node[T], isLast bool, isRoot bool) string
}

// ShowHierarchy displays a tree structure of files/directories using the new generic tree system
func ShowHierarchy(basePath, targetDir string) (error, bool) {
	// Create filesystem tree builder with default config
	builder := NewFileSystemTreeBuilder()

	// Build the tree
	tree, err := builder.Build(basePath)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err), false
	}

	// Check if tree has only one node and it's not a directory
	root := tree.Root()
	if len(root.Children) == 1 && !root.Children[0].Data.IsDir {
		return nil, false // No hierarchy needed
	}

	// Create filesystem styler with default config
	styler := NewFileSystemStyler()

	// Create tree renderer
	renderer := NewTreeRenderer(styler)

	// Render the tree
	err = renderer.Render(tree, os.Stdout)
	if err != nil {
		return fmt.Errorf("failed to render tree: %w", err), false
	}

	return nil, true
}
