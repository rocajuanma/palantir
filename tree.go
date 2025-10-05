package palantir

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TreeNode represents a simple tree node for display purposes only
type TreeNode struct {
	Name     string
	Data     FileNode
	Children []*TreeNode
}

// FileNode represents a file or directory in the filesystem tree
type FileNode struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime int64
}

// ShowHierarchy displays a tree structure of files/directories
func ShowHierarchy(basePath, targetDir string) (error, bool) {
	// Get root directory info
	rootInfo, err := os.Stat(basePath)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err), false
	}

	// Create root node
	root := &TreeNode{
		Name: rootInfo.Name(),
		Data: FileNode{
			Name:    rootInfo.Name(),
			Path:    basePath,
			IsDir:   rootInfo.IsDir(),
			Size:    rootInfo.Size(),
			ModTime: rootInfo.ModTime().Unix(),
		},
		Children: []*TreeNode{},
	}

	// Build tree structure by walking filesystem
	err = buildTree(root, basePath)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err), false
	}

	// Check if tree has only one node and it's not a directory
	if len(root.Children) == 1 && !root.Children[0].Data.IsDir {
		return nil, false // No hierarchy needed
	}

	// Sort tree (directories first, then alphabetically)
	sortTree(root)

	// Render the tree
	printTree(root, "", true, true)

	return nil, true
}

// buildTree recursively builds a tree structure from the filesystem
func buildTree(node *TreeNode, dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == dirPath {
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
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		parts := strings.Split(relPath, string(filepath.Separator))

		// Find or create the parent node
		current := node
		for i, part := range parts[:len(parts)-1] {
			found := false
			for _, child := range current.Children {
				if child.Name == part && child.Data.IsDir {
					current = child
					found = true
					break
				}
			}
			if !found {
				// Create intermediate directory
				newDir := &TreeNode{
					Name: part,
					Data: FileNode{
						Name:  part,
						Path:  filepath.Join(dirPath, strings.Join(parts[:i+1], string(filepath.Separator))),
						IsDir: true,
					},
					Children: []*TreeNode{},
				}
				current.Children = append(current.Children, newDir)
				current = newDir
			}
		}

		// Add the final node
		finalNode := &TreeNode{
			Name: parts[len(parts)-1],
			Data: FileNode{
				Name:    info.Name(),
				Path:    path,
				IsDir:   info.IsDir(),
				Size:    info.Size(),
				ModTime: info.ModTime().Unix(),
			},
			Children: []*TreeNode{},
		}
		current.Children = append(current.Children, finalNode)

		return nil
	})
}

// sortTree recursively sorts all children in the tree (directories first, then files, both alphabetically)
func sortTree(node *TreeNode) {
	if len(node.Children) == 0 {
		return
	}

	// Sort children: directories first, then files, both alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].Data.IsDir != node.Children[j].Data.IsDir {
			return node.Children[i].Data.IsDir // directories come first
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		sortTree(child)
	}
}

// printTree recursively prints a tree node with ASCII art and colors
func printTree(node *TreeNode, prefix string, isLast bool, isRoot bool) {
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
func styleFileNode(node *TreeNode) string {
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
