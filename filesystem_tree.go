package palantir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileNode represents a file or directory in the filesystem tree
type FileNode struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime int64
}

// NewFileNode creates a new FileNode
func NewFileNode(name, path string, isDir bool, size int64, modTime int64) FileNode {
	return FileNode{
		Name:    name,
		Path:    path,
		IsDir:   isDir,
		Size:    size,
		ModTime: modTime,
	}
}

// FileSystemTreeBuilder implements TreeBuilder for filesystem trees
type FileSystemTreeBuilder struct {
	outputConfig *OutputConfig
}

// NewFileSystemTreeBuilder creates a new filesystem tree builder
func NewFileSystemTreeBuilder() *FileSystemTreeBuilder {
	return &FileSystemTreeBuilder{
		outputConfig: GetGlobalOutputHandler().(*outputHandler).config,
	}
}

// Build builds a tree from the filesystem
func (b *FileSystemTreeBuilder) Build(source string) (Tree[FileNode], error) {
	// Get absolute path
	absPath, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path does not exist: %w", err)
	}

	// Create root node
	rootData := NewFileNode(info.Name(), absPath, info.IsDir(), info.Size(), info.ModTime().Unix())
	tree := NewTree(rootData, info.Name())

	// Build tree recursively
	err = b.buildTreeRecursive(tree, absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}

	// Sort the tree (directories first, then alphabetically)
	tree.Sort(b.getSortFunction())

	return tree, nil
}

// buildTreeRecursive recursively builds the tree structure
func (b *FileSystemTreeBuilder) buildTreeRecursive(tree Tree[FileNode], dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		// Skip hidden files (simple default behavior)
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path from root
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Split the path into components
		parts := strings.Split(relPath, string(filepath.Separator))

		// Create file node data
		fileData := NewFileNode(info.Name(), path, info.IsDir(), info.Size(), info.ModTime().Unix())

		// Insert into tree
		err = tree.Insert(parts, fileData)
		if err != nil {
			return fmt.Errorf("failed to insert node: %w", err)
		}

		return nil
	})
}

// getSortFunction returns the sort function (directories first, then alphabetically)
func (b *FileSystemTreeBuilder) getSortFunction() func(a, b *Node[FileNode]) bool {
	return func(a, b *Node[FileNode]) bool {
		// Directories first
		if a.Data.IsDir != b.Data.IsDir {
			return a.Data.IsDir
		}
		// Alphabetical sort
		return a.Data.Name < b.Data.Name
	}
}
