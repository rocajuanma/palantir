package palantir

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// TreeRendererImpl implements TreeRenderer for generic trees
type TreeRendererImpl[T any] struct {
	styler       NodeStyler[T]
	outputConfig *OutputConfig
}

// NewTreeRenderer creates a new tree renderer
func NewTreeRenderer[T any](styler NodeStyler[T]) *TreeRendererImpl[T] {
	return &TreeRendererImpl[T]{
		styler:       styler,
		outputConfig: GetGlobalOutputHandler().(*outputHandler).config,
	}
}

// Render renders a tree to the writer
func (r *TreeRendererImpl[T]) Render(tree Tree[T], writer io.Writer) error {
	root := tree.Root()
	if root == nil {
		return fmt.Errorf("tree has no root")
	}

	// Check if tree has only one node and it's not a directory (for filesystem trees)
	if len(root.Children) == 1 && !r.isDirectory(root.Children[0]) {
		return nil // No hierarchy needed
	}

	// Render the tree
	return r.renderNode(root, writer, "", true, true)
}

// renderNode recursively renders a node and its children
func (r *TreeRendererImpl[T]) renderNode(node *Node[T], writer io.Writer, prefix string, isLast bool, isRoot bool) error {
	if !isRoot {
		// Get tree character
		treeChar := r.styler.GetTreeChar(node, isLast)

		// Get styled node name
		styledName := r.styler.Style(node)

		// Write the node
		_, err := fmt.Fprintf(writer, "%s%s%s\n", prefix, treeChar, styledName)
		if err != nil {
			return err
		}
	}

	// Render children
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1

			// Calculate prefix for child
			childPrefix := r.styler.GetPrefix(node, isChildLast, isRoot)

			// Recursively render child
			err := r.renderNode(child, writer, childPrefix, isChildLast, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// isDirectory checks if a node represents a directory (for filesystem trees)
func (r *TreeRendererImpl[T]) isDirectory(node *Node[T]) bool {
	// This is a type assertion that works for FileNode
	if fileNode, ok := any(node.Data).(FileNode); ok {
		return fileNode.IsDir
	}
	return false
}

// FileSystemStyler implements NodeStyler for filesystem trees
type FileSystemStyler struct {
	outputConfig *OutputConfig
}

// NewFileSystemStyler creates a new filesystem styler
func NewFileSystemStyler() *FileSystemStyler {
	return &FileSystemStyler{
		outputConfig: GetGlobalOutputHandler().(*outputHandler).config,
	}
}

// Style styles a filesystem node
func (s *FileSystemStyler) Style(node *Node[FileNode]) string {
	if !s.outputConfig.UseColors {
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

// GetTreeChar returns the appropriate tree character for a node
func (s *FileSystemStyler) GetTreeChar(node *Node[FileNode], isLast bool) string {
	if isLast {
		return "└── "
	}
	return "├── "
}

// GetPrefix returns the prefix for a node
func (s *FileSystemStyler) GetPrefix(node *Node[FileNode], isLast bool, isRoot bool) string {
	if isRoot {
		return ""
	}

	if isLast {
		return "    " // 4 spaces
	}
	return "│   " // vertical line + 3 spaces
}
