package palantir

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TreeNode represents a node in the tree
type TreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*TreeNode
}

// ShowHierarchy displays a tree structure of files/directories
func ShowHierarchy(basePath, targetDir string) (error, bool) {
	root, err := buildTree(basePath)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err), false
	}

	if len(root.Children) == 1 && !root.Children[0].IsDir {
		return nil, false // if there is only one node, return false because no hierarchy is needed
	}

	// Sort children for consistent display
	sortChildren(root)

	// Print the tree starting from root
	printTreeNode(root, "", true, true)
	return nil, true
}

// buildTree recursively builds a tree structure from the filesystem
func buildTree(dirPath string) (*TreeNode, error) {
	root := &TreeNode{
		Name:     filepath.Base(dirPath),
		Path:     dirPath,
		IsDir:    true,
		Children: []*TreeNode{},
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		// Get relative path from root
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Split the path into components
		parts := strings.Split(relPath, string(filepath.Separator))

		// Find or create the parent node
		current := root
		for i, part := range parts[:len(parts)-1] {
			found := false
			for _, child := range current.Children {
				if child.Name == part && child.IsDir {
					current = child
					found = true
					break
				}
			}
			if !found {
				// Create intermediate directory
				newDir := &TreeNode{
					Name:     part,
					Path:     filepath.Join(dirPath, strings.Join(parts[:i+1], string(filepath.Separator))),
					IsDir:    true,
					Children: []*TreeNode{},
				}
				current.Children = append(current.Children, newDir)
				current = newDir
			}
		}

		// Add the final node
		finalNode := &TreeNode{
			Name:  parts[len(parts)-1],
			Path:  path,
			IsDir: info.IsDir(),
		}
		if info.IsDir() {
			finalNode.Children = []*TreeNode{}
		}
		current.Children = append(current.Children, finalNode)

		return nil
	})

	return root, err
}

// sortChildren recursively sorts all children in the tree (directories first, then files, both alphabetically)
func sortChildren(node *TreeNode) {
	if node.Children == nil {
		return
	}

	// Sort children: directories first, then files, both alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir // directories come first
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		sortChildren(child)
	}
}

// printTreeNode prints a tree node with ASCII art and colors
func printTreeNode(node *TreeNode, prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		// Choose the appropriate tree character
		var treeChar string
		if isLast {
			treeChar = "└── "
		} else {
			treeChar = "├── "
		}

		// Color the output based on file type
		var coloredName string
		if node.IsDir {
			coloredName = fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, node.Name, ColorReset)
		} else {
			// Color files based on extension
			ext := strings.ToLower(filepath.Ext(node.Name))
			switch ext {
			case ".json", ".yaml", ".yml", ".toml":
				coloredName = fmt.Sprintf("%s%s%s", ColorGreen, node.Name, ColorReset)
			case ".md", ".txt", ".log":
				coloredName = fmt.Sprintf("%s%s%s", ColorCyan, node.Name, ColorReset)
			case ".sh", ".zsh", ".bash":
				coloredName = fmt.Sprintf("%s%s%s", ColorYellow, node.Name, ColorReset)
			default:
				coloredName = node.Name
			}
		}

		// Print the current node
		fmt.Printf("%s%s%s\n", prefix, treeChar, coloredName)
	}

	// Print children
	if node.Children != nil {
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

			printTreeNode(child, childPrefix, isChildLast, false)
		}
	}
}
