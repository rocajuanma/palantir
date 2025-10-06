package palantir

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Tree display constants
const (
	Branch   = "├── "
	Last     = "└── "
	Vertical = "│   "
	Space    = "    "
)

// TreeNode represents a simple tree node for display purposes only
type TreeNode struct {
	Name     string
	Data     interface{} // Can be FileNode or YAMLNode
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

	root := &TreeNode{
		Name: rootInfo.Name(),
		Data: FileNode{
			Name:    rootInfo.Name(),
			Path:    basePath,
			IsDir:   rootInfo.IsDir(),
			Size:    rootInfo.Size(),
			ModTime: rootInfo.ModTime().Unix(),
		},
		Children: nil,
	}

	// Build tree structure by walking filesystem
	err = buildTree(root, basePath)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err), false
	}

	// Check if tree has only one node and it's not a directory
	if len(root.Children) == 1 && !getIsDir(root.Children[0].Data) {
		return nil, false // No hierarchy needed
	}

	// Directories first, then alphabetically
	sortTree(root)
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
			// Use map for O(1) lookup
			childMap := make(map[string]*TreeNode)
			for _, child := range current.Children {
				if getIsDir(child.Data) {
					childMap[child.Name] = child
				}
			}

			if existingChild, found := childMap[part]; found {
				current = existingChild
			} else {
				// Create intermediate directory
				newDir := &TreeNode{
					Name: part,
					Data: FileNode{
						Name:  part,
						Path:  filepath.Join(dirPath, strings.Join(parts[:i+1], string(filepath.Separator))),
						IsDir: true,
					},
					Children: nil,
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
			Children: nil,
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
		// Get IsDir from the appropriate data type
		iIsDir := getIsDir(node.Children[i].Data)
		jIsDir := getIsDir(node.Children[j].Data)

		if iIsDir != jIsDir {
			return iIsDir // directories come first
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		sortTree(child)
	}
}

// getIsDir extracts IsDir from either FileNode or YAMLNode
func getIsDir(data interface{}) bool {
	if fileNode, ok := data.(FileNode); ok {
		return fileNode.IsDir
	}
	if yamlNode, ok := data.(YAMLNode); ok {
		return yamlNode.IsDir
	}
	return false
}

// printTree recursively prints a tree node with ASCII art and colors
func printTree(node *TreeNode, prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		// Choose the appropriate tree character
		var treeChar string
		if isLast {
			treeChar = Last
		} else {
			treeChar = Branch
		}

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
					childPrefix = prefix + Space
				} else {
					childPrefix = prefix + Vertical
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
		return node.Name
	}

	// Handle FileNode
	if fileNode, ok := node.Data.(FileNode); ok {
		if fileNode.IsDir {
			return fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, fileNode.Name, ColorReset)
		}

		// Color customized based on extension
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

	// Handle YAMLNode
	if yamlNode, ok := node.Data.(YAMLNode); ok {
		if yamlNode.IsDir {
			return fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, yamlNode.Name, ColorReset)
		}

		// Color based on node type
		switch yamlNode.NodeType {
		case "object":
			return fmt.Sprintf("%s%s%s%s", ColorBold, ColorBlue, yamlNode.Name, ColorReset)
		case "array":
			return fmt.Sprintf("%s%s%s", ColorYellow, yamlNode.Name, ColorReset)
		case "scalar":
			return fmt.Sprintf("%s%s%s", ColorGreen, yamlNode.Name, ColorReset)
		default:
			return yamlNode.Name
		}
	}

	// Fallback
	return node.Name
}

// YAMLNode represents a YAML data node for tree visualization
type YAMLNode struct {
	Name     string
	Value    interface{}
	IsDir    bool
	NodeType string // "object", "array", "scalar"
}

// ParseYAMLToTree converts YAML content to TreeNode structure
func ParseYAMLToTree(yamlContent []byte) (*TreeNode, error) {
	var data interface{}
	if err := yaml.Unmarshal(yamlContent, &data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	root := &TreeNode{
		Name:     "root",
		Data:     YAMLNode{Name: "root", Value: data, IsDir: true, NodeType: "object"},
		Children: nil,
	}

	return buildYAMLTree(root, data), nil
}

// buildYAMLTree recursively builds a tree structure from YAML data
func buildYAMLTree(node *TreeNode, data interface{}) *TreeNode {
	switch v := data.(type) {
	case map[string]interface{}:
		// Handle objects
		for key, value := range v {
			child := &TreeNode{
				Name:     key,
				Data:     YAMLNode{Name: key, Value: value, IsDir: true, NodeType: "object"},
				Children: nil,
			}
			node.Children = append(node.Children, buildYAMLTree(child, value))
		}
	case []interface{}:
		// Handle arrays
		for i, item := range v {
			// Create a name with just the value for array items
			var itemName string
			switch itemValue := item.(type) {
			case string:
				itemName = itemValue
			case int, int64, float64:
				itemName = fmt.Sprintf("%v", itemValue)
			case bool:
				itemName = fmt.Sprintf("%t", itemValue)
			default:
				itemName = fmt.Sprintf("[%d]", i)
			}

			child := &TreeNode{
				Name:     itemName,
				Data:     YAMLNode{Name: itemName, Value: item, IsDir: false, NodeType: "array"},
				Children: nil,
			}
			// Only recursively build if the item is a complex type (map or slice)
			switch item.(type) {
			case map[string]interface{}, []interface{}:
				node.Children = append(node.Children, buildYAMLTree(child, item))
			default:
				// For scalar values, just add the child as-is
				node.Children = append(node.Children, child)
			}
		}
	default:
		// Handle scalar values
		node.Data = YAMLNode{Name: node.Name, Value: v, IsDir: false, NodeType: "scalar"}
	}
	return node
}

// ShowYAMLHierarchy displays YAML content as a tree structure
func ShowYAMLHierarchy(yamlContent []byte) error {
	root, err := ParseYAMLToTree(yamlContent)
	if err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	sortTree(root)
	printTree(root, "", true, true)
	return nil
}

// ShowYAMLHierarchyFromFile reads and displays a YAML file as a tree structure
func ShowYAMLHierarchyFromFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}
	return ShowYAMLHierarchy(content)
}
