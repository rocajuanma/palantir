package palantir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTree(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "palantir_build_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFiles := []string{
		"file1.txt",
		"dir1/file2.go",
		"dir1/subdir/file3.md",
		"dir2/file4.json",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Test buildTree function
	root := &TreeNode{
		Name:     filepath.Base(tempDir),
		Data:     FileNode{Name: filepath.Base(tempDir), Path: tempDir, IsDir: true},
		Children: nil,
	}

	err = buildTree(root, tempDir)
	if err != nil {
		t.Fatalf("buildTree() error = %v", err)
	}

	// Verify tree structure
	if len(root.Children) != 3 { // file1.txt, dir1, dir2
		t.Errorf("Expected 3 children, got %d", len(root.Children))
	}

	// Find dir1 and verify its structure
	var dir1 *TreeNode
	for _, child := range root.Children {
		if child.Name == "dir1" {
			dir1 = child
			break
		}
	}

	if dir1 == nil {
		t.Fatal("dir1 not found in tree")
	}

	if len(dir1.Children) != 2 { // file2.go, subdir
		t.Errorf("Expected dir1 to have 2 children, got %d", len(dir1.Children))
	}

	// Find subdir and verify its structure
	var subdir *TreeNode
	for _, child := range dir1.Children {
		if child.Name == "subdir" {
			subdir = child
			break
		}
	}

	if subdir == nil {
		t.Fatal("subdir not found in tree")
	}

	if len(subdir.Children) != 1 { // file3.md
		t.Errorf("Expected subdir to have 1 child, got %d", len(subdir.Children))
	}
}

func TestSortTree(t *testing.T) {
	// Create a test tree with mixed files and directories
	root := &TreeNode{
		Name: "root",
		Data: FileNode{Name: "root", IsDir: true},
		Children: []*TreeNode{
			{Name: "file3.txt", Data: FileNode{Name: "file3.txt", IsDir: false}},
			{Name: "dir1", Data: FileNode{Name: "dir1", IsDir: true}},
			{Name: "file1.go", Data: FileNode{Name: "file1.go", IsDir: false}},
			{Name: "dir2", Data: FileNode{Name: "dir2", IsDir: true}},
			{Name: "file2.md", Data: FileNode{Name: "file2.md", IsDir: false}},
		},
	}

	sortTree(root)

	// Verify directories come first
	expectedOrder := []string{"dir1", "dir2", "file1.go", "file2.md", "file3.txt"}
	actualOrder := make([]string, len(root.Children))
	for i, child := range root.Children {
		actualOrder[i] = child.Name
	}

	for i, expected := range expectedOrder {
		if actualOrder[i] != expected {
			t.Errorf("Expected child %d to be %q, got %q", i, expected, actualOrder[i])
		}
	}

	// Verify directories are first
	for i, child := range root.Children {
		if i < 2 && !getIsDir(child.Data) {
			t.Errorf("Expected directories to come first, but found file %q at position %d", child.Name, i)
		}
		if i >= 2 && getIsDir(child.Data) {
			t.Errorf("Expected files to come after directories, but found directory %q at position %d", child.Name, i)
		}
	}
}

func TestStyleFileNode(t *testing.T) {
	// Test with colors enabled
	outputConfig := &OutputConfig{
		UseColors:         true,
		UseEmojis:         false,
		UseFormatting:     true,
		DisableOutput:     false,
		VerboseMode:       false,
		ColorizeLevelOnly: false,
	}

	// Set global output handler
	SetGlobalOutputHandler(NewOutputHandler(outputConfig))

	tests := []struct {
		name     string
		node     *TreeNode
		expected string
	}{
		{
			name: "Directory",
			node: &TreeNode{
				Name: "testdir",
				Data: FileNode{Name: "testdir", IsDir: true},
			},
			expected: "testdir", // Should contain color codes
		},
		{
			name: "Go file",
			node: &TreeNode{
				Name: "main.go",
				Data: FileNode{Name: "main.go", IsDir: false},
			},
			expected: "main.go", // Should contain color codes
		},
		{
			name: "Markdown file",
			node: &TreeNode{
				Name: "README.md",
				Data: FileNode{Name: "README.md", IsDir: false},
			},
			expected: "README.md", // Should contain color codes
		},
		{
			name: "JSON file",
			node: &TreeNode{
				Name: "config.json",
				Data: FileNode{Name: "config.json", IsDir: false},
			},
			expected: "config.json", // Should contain color codes
		},
		{
			name: "Shell script",
			node: &TreeNode{
				Name: "script.sh",
				Data: FileNode{Name: "script.sh", IsDir: false},
			},
			expected: "script.sh", // Should contain color codes
		},
		{
			name: "Unknown file type",
			node: &TreeNode{
				Name: "unknown.xyz",
				Data: FileNode{Name: "unknown.xyz", IsDir: false},
			},
			expected: "unknown.xyz", // Should not contain color codes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := styleFileNode(tt.node)

			// For directories and known file types, should contain color codes
			if getIsDir(tt.node.Data) ||
				strings.HasSuffix(tt.node.Name, ".go") ||
				strings.HasSuffix(tt.node.Name, ".md") ||
				strings.HasSuffix(tt.node.Name, ".json") ||
				strings.HasSuffix(tt.node.Name, ".sh") {
				if !strings.Contains(result, ColorReset) {
					t.Errorf("Expected styled output to contain color codes, got: %q", result)
				}
			} else {
				// For unknown file types, should not contain color codes
				if strings.Contains(result, ColorReset) {
					t.Errorf("Expected plain output for unknown file type, got: %q", result)
				}
			}

			// Should always contain the filename
			if !strings.Contains(result, tt.node.Name) {
				t.Errorf("Expected output to contain filename %q, got: %q", tt.node.Name, result)
			}
		})
	}
}

func TestStyleFileNodeNoColors(t *testing.T) {
	// Test with colors disabled
	outputConfig := &OutputConfig{
		UseColors:         false,
		UseEmojis:         false,
		UseFormatting:     true,
		DisableOutput:     false,
		VerboseMode:       false,
		ColorizeLevelOnly: false,
	}

	// Set global output handler
	SetGlobalOutputHandler(NewOutputHandler(outputConfig))

	node := &TreeNode{
		Name: "testfile.go",
		Data: FileNode{Name: "testfile.go", IsDir: false},
	}

	result := styleFileNode(node)
	expected := "testfile.go"

	if result != expected {
		t.Errorf("Expected plain output %q, got: %q", expected, result)
	}
}

func TestShowHierarchyBasic(t *testing.T) {
	// Create a simple test directory
	tempDir, err := os.MkdirTemp("", "palantir_hierarchy_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple files to ensure hierarchy
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.go")
	if err := os.WriteFile(testFile1, []byte("test1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("test2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	// Test with a single file (should return true because it's a file, not a directory)
	err, hasHierarchy := ShowHierarchy(testFile1, "")
	if err != nil {
		t.Errorf("ShowHierarchy() error = %v", err)
	}
	// For a single file, ShowHierarchy will still return true because it's not a directory
	// The logic checks for children, but a single file has no children
	if !hasHierarchy {
		t.Errorf("ShowHierarchy() hasHierarchy = %v, want true for single file", hasHierarchy)
	}

	// Test with a directory containing multiple files (should return true for hierarchy)
	err, hasHierarchy = ShowHierarchy(tempDir, "")
	if err != nil {
		t.Errorf("ShowHierarchy() error = %v", err)
	}
	// The directory contains multiple files, so it should have hierarchy
	if !hasHierarchy {
		t.Errorf("ShowHierarchy() hasHierarchy = %v, want true for directory with multiple files", hasHierarchy)
	}
}

func TestBuildTreeErrorHandling(t *testing.T) {
	// Test with non-existent directory
	root := &TreeNode{
		Name:     "nonexistent",
		Data:     FileNode{Name: "nonexistent", Path: "/nonexistent/path", IsDir: true},
		Children: nil,
	}

	err := buildTree(root, "/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}

func TestShowHierarchyErrorHandling(t *testing.T) {
	// Test with non-existent path
	err, hasHierarchy := ShowHierarchy("/nonexistent/path", "")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
	if hasHierarchy {
		t.Error("Expected hasHierarchy to be false for non-existent path")
	}
}

func TestBuildTreeEmptyDirectory(t *testing.T) {
	// Create empty directory
	tempDir, err := os.MkdirTemp("", "palantir_empty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	root := &TreeNode{
		Name:     filepath.Base(tempDir),
		Data:     FileNode{Name: filepath.Base(tempDir), Path: tempDir, IsDir: true},
		Children: nil,
	}

	err = buildTree(root, tempDir)
	if err != nil {
		t.Fatalf("buildTree() error = %v", err)
	}

	// Empty directory should have no children
	if len(root.Children) != 0 {
		t.Errorf("Expected 0 children for empty directory, got %d", len(root.Children))
	}
}

func TestBuildTreeDeepStructure(t *testing.T) {
	// Create deep directory structure
	tempDir, err := os.MkdirTemp("", "palantir_deep_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create deep nested structure
	deepPath := filepath.Join(tempDir, "level1", "level2", "level3", "level4", "level5")
	if err := os.MkdirAll(deepPath, 0755); err != nil {
		t.Fatalf("Failed to create deep directory: %v", err)
	}

	// Create file at deep level
	deepFile := filepath.Join(deepPath, "deepfile.txt")
	if err := os.WriteFile(deepFile, []byte("deep content"), 0644); err != nil {
		t.Fatalf("Failed to create deep file: %v", err)
	}

	root := &TreeNode{
		Name:     filepath.Base(tempDir),
		Data:     FileNode{Name: filepath.Base(tempDir), Path: tempDir, IsDir: true},
		Children: nil,
	}

	err = buildTree(root, tempDir)
	if err != nil {
		t.Fatalf("buildTree() error = %v", err)
	}

	// Verify deep structure was built correctly
	if len(root.Children) != 1 {
		t.Errorf("Expected 1 child at root, got %d", len(root.Children))
	}

	level1 := root.Children[0]
	if level1.Name != "level1" {
		t.Errorf("Expected level1, got %s", level1.Name)
	}

	// Traverse down to verify depth
	current := level1
	for i := 2; i <= 5; i++ {
		if len(current.Children) != 1 {
			t.Errorf("Expected 1 child at level %d, got %d", i-1, len(current.Children))
		}
		current = current.Children[0]
		expectedName := fmt.Sprintf("level%d", i)
		if current.Name != expectedName {
			t.Errorf("Expected %s, got %s", expectedName, current.Name)
		}
	}

	// Verify file is at the deepest level
	if len(current.Children) != 1 {
		t.Errorf("Expected 1 file at deepest level, got %d", len(current.Children))
	}
	if current.Children[0].Name != "deepfile.txt" {
		t.Errorf("Expected deepfile.txt, got %s", current.Children[0].Name)
	}
}

func TestBuildTreeSpecialCharacters(t *testing.T) {
	// Create directory with special characters
	tempDir, err := os.MkdirTemp("", "palantir_special_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with special characters
	specialFiles := []string{
		"file with spaces.txt",
		"file-with-dashes.go",
		"file_with_underscores.md",
		"file.with.dots.json",
	}

	for _, file := range specialFiles {
		fullPath := filepath.Join(tempDir, file)
		if err := os.WriteFile(fullPath, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	root := &TreeNode{
		Name:     filepath.Base(tempDir),
		Data:     FileNode{Name: filepath.Base(tempDir), Path: tempDir, IsDir: true},
		Children: nil,
	}

	err = buildTree(root, tempDir)
	if err != nil {
		t.Fatalf("buildTree() error = %v", err)
	}

	// Verify all special files are present
	if len(root.Children) != len(specialFiles) {
		t.Errorf("Expected %d children, got %d", len(specialFiles), len(root.Children))
	}

	// Check that all special files are found
	foundFiles := make(map[string]bool)
	for _, child := range root.Children {
		foundFiles[child.Name] = true
	}

	for _, expectedFile := range specialFiles {
		if !foundFiles[expectedFile] {
			t.Errorf("Expected to find file %s, but it was not found", expectedFile)
		}
	}
}

func TestStyleFileNodeExtendedTypes(t *testing.T) {
	// Test with colors enabled
	outputConfig := &OutputConfig{
		UseColors:         true,
		UseEmojis:         false,
		UseFormatting:     true,
		DisableOutput:     false,
		VerboseMode:       false,
		ColorizeLevelOnly: false,
	}

	SetGlobalOutputHandler(NewOutputHandler(outputConfig))

	tests := []struct {
		name            string
		node            *TreeNode
		shouldHaveColor bool
	}{
		{"YAML file", &TreeNode{Name: "config.yaml", Data: FileNode{Name: "config.yaml", IsDir: false}}, true},
		{"XML file", &TreeNode{Name: "data.xml", Data: FileNode{Name: "data.xml", IsDir: false}}, false},      // Not supported
		{"CSS file", &TreeNode{Name: "style.css", Data: FileNode{Name: "style.css", IsDir: false}}, false},    // Not supported
		{"HTML file", &TreeNode{Name: "index.html", Data: FileNode{Name: "index.html", IsDir: false}}, false}, // Not supported
		{"Python file", &TreeNode{Name: "script.py", Data: FileNode{Name: "script.py", IsDir: false}}, false}, // Not supported
		{"JavaScript file", &TreeNode{Name: "app.js", Data: FileNode{Name: "app.js", IsDir: false}}, false},   // Not supported
		{"TypeScript file", &TreeNode{Name: "app.ts", Data: FileNode{Name: "app.ts", IsDir: false}}, false},   // Not supported
		{"Rust file", &TreeNode{Name: "main.rs", Data: FileNode{Name: "main.rs", IsDir: false}}, false},       // Not supported
		{"C file", &TreeNode{Name: "main.c", Data: FileNode{Name: "main.c", IsDir: false}}, false},            // Not supported
		{"C++ file", &TreeNode{Name: "main.cpp", Data: FileNode{Name: "main.cpp", IsDir: false}}, false},      // Not supported
		{"Java file", &TreeNode{Name: "Main.java", Data: FileNode{Name: "Main.java", IsDir: false}}, false},   // Not supported
		{"PHP file", &TreeNode{Name: "index.php", Data: FileNode{Name: "index.php", IsDir: false}}, false},    // Not supported
		{"Ruby file", &TreeNode{Name: "app.rb", Data: FileNode{Name: "app.rb", IsDir: false}}, false},         // Not supported
		{"File without extension", &TreeNode{Name: "README", Data: FileNode{Name: "README", IsDir: false}}, false},
		{"Hidden file", &TreeNode{Name: ".gitignore", Data: FileNode{Name: ".gitignore", IsDir: false}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := styleFileNode(tt.node)

			if tt.shouldHaveColor {
				if !strings.Contains(result, ColorReset) {
					t.Errorf("Expected %s to have color codes, got: %q", tt.name, result)
				}
			} else {
				if strings.Contains(result, ColorReset) {
					t.Errorf("Expected %s to not have color codes, got: %q", tt.name, result)
				}
			}

			// Should always contain the filename
			if !strings.Contains(result, tt.node.Name) {
				t.Errorf("Expected output to contain filename %q, got: %q", tt.node.Name, result)
			}
		})
	}
}

func TestSortTreeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		children []*TreeNode
		expected []string
	}{
		{
			name:     "Empty children",
			children: []*TreeNode{},
			expected: []string{},
		},
		{
			name: "Single file",
			children: []*TreeNode{
				{Name: "file.txt", Data: FileNode{Name: "file.txt", IsDir: false}},
			},
			expected: []string{"file.txt"},
		},
		{
			name: "Single directory",
			children: []*TreeNode{
				{Name: "dir", Data: FileNode{Name: "dir", IsDir: true}},
			},
			expected: []string{"dir"},
		},
		{
			name: "All directories",
			children: []*TreeNode{
				{Name: "dir3", Data: FileNode{Name: "dir3", IsDir: true}},
				{Name: "dir1", Data: FileNode{Name: "dir1", IsDir: true}},
				{Name: "dir2", Data: FileNode{Name: "dir2", IsDir: true}},
			},
			expected: []string{"dir1", "dir2", "dir3"},
		},
		{
			name: "All files",
			children: []*TreeNode{
				{Name: "file3.txt", Data: FileNode{Name: "file3.txt", IsDir: false}},
				{Name: "file1.go", Data: FileNode{Name: "file1.go", IsDir: false}},
				{Name: "file2.md", Data: FileNode{Name: "file2.md", IsDir: false}},
			},
			expected: []string{"file1.go", "file2.md", "file3.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &TreeNode{
				Name:     "root",
				Data:     FileNode{Name: "root", IsDir: true},
				Children: tt.children,
			}

			sortTree(root)

			actualOrder := make([]string, len(root.Children))
			for i, child := range root.Children {
				actualOrder[i] = child.Name
			}

			if len(actualOrder) != len(tt.expected) {
				t.Errorf("Expected %d children, got %d", len(tt.expected), len(actualOrder))
				return
			}

			for i, expected := range tt.expected {
				if actualOrder[i] != expected {
					t.Errorf("Expected child %d to be %q, got %q", i, expected, actualOrder[i])
				}
			}
		})
	}
}

func TestShowHierarchyEmptyDirectory(t *testing.T) {
	// Create empty directory
	tempDir, err := os.MkdirTemp("", "palantir_empty_hierarchy_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with empty directory
	err, hasHierarchy := ShowHierarchy(tempDir, "")
	if err != nil {
		t.Errorf("ShowHierarchy() error = %v", err)
	}
	// Empty directory should have hierarchy (current logic returns true for any directory)
	if !hasHierarchy {
		t.Errorf("ShowHierarchy() hasHierarchy = %v, want true for empty directory", hasHierarchy)
	}
}

func TestShowHierarchySingleFileInDirectory(t *testing.T) {
	// Create directory with single file
	tempDir, err := os.MkdirTemp("", "palantir_single_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create single file
	singleFile := filepath.Join(tempDir, "single.txt")
	if err := os.WriteFile(singleFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create single file: %v", err)
	}

	// Test with directory containing single file
	err, hasHierarchy := ShowHierarchy(tempDir, "")
	if err != nil {
		t.Errorf("ShowHierarchy() error = %v", err)
	}
	// Directory with single file should not have hierarchy (current logic returns false for single non-directory child)
	if hasHierarchy {
		t.Errorf("ShowHierarchy() hasHierarchy = %v, want false for directory with single file", hasHierarchy)
	}
}

func TestBuildTreePermissionError(t *testing.T) {
	// Test with a path that might cause permission issues
	// This test might not work on all systems, so we'll make it conditional

	// Try to access a system directory that might be restricted
	restrictedPath := "/root" // This should fail on most systems

	root := &TreeNode{
		Name:     "root",
		Data:     FileNode{Name: "root", Path: restrictedPath, IsDir: true},
		Children: nil,
	}

	err := buildTree(root, restrictedPath)
	// We expect this to fail due to permission issues
	if err == nil {
		t.Log("Warning: buildTree succeeded with restricted path, this might indicate running as root")
	} else {
		t.Logf("Expected error for restricted path: %v", err)
	}
}

func TestShowHierarchyInvalidPath(t *testing.T) {
	// Test with various invalid paths
	invalidPaths := []string{
		"",                   // Empty path
		"/nonexistent/path",  // Non-existent path
		"\x00invalid",        // Path with null character
		"path/with/\x00null", // Path with null character in middle
	}

	for _, invalidPath := range invalidPaths {
		t.Run(fmt.Sprintf("InvalidPath_%q", invalidPath), func(t *testing.T) {
			err, hasHierarchy := ShowHierarchy(invalidPath, "")

			// Should return an error for invalid paths
			if err == nil {
				t.Errorf("Expected error for invalid path %q, got nil", invalidPath)
			}

			// Should not have hierarchy for invalid paths
			if hasHierarchy {
				t.Errorf("Expected hasHierarchy=false for invalid path %q, got true", invalidPath)
			}
		})
	}
}

func TestBuildTreeWithNilNode(t *testing.T) {
	// Test buildTree with nil node
	// Note: Current implementation doesn't check for nil, so this will panic
	// This test documents the current behavior
	defer func() {
		if r := recover(); r != nil {
			t.Logf("buildTree with nil node panicked (expected): %v", r)
		}
	}()

	err := buildTree(nil, "/tmp")
	if err == nil {
		t.Log("buildTree with nil node succeeded (unexpected)")
	}
}

func TestBuildTreeWithInvalidNodeData(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "palantir_invalid_node_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with node that has invalid path data
	root := &TreeNode{
		Name: "invalid",
		Data: FileNode{
			Name:  "invalid",
			Path:  "", // Empty path should cause issues
			IsDir: true,
		},
		Children: nil,
	}

	err = buildTree(root, "")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

func TestShowHierarchyWithEmptyPath(t *testing.T) {
	// Test ShowHierarchy with empty path
	err, hasHierarchy := ShowHierarchy("", "")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
	if hasHierarchy {
		t.Error("Expected hasHierarchy=false for empty path, got true")
	}
}

func TestBuildTreeCircularReference(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "palantir_circular_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create a symlink that points to parent (circular reference)
	symlinkPath := filepath.Join(subDir, "parent_link")
	parentPath := tempDir

	// Note: This test might not work on all systems or might require special permissions
	// We'll make it conditional
	if err := os.Symlink(parentPath, symlinkPath); err != nil {
		t.Logf("Could not create symlink (might not be supported): %v", err)
		return
	}

	root := &TreeNode{
		Name:     filepath.Base(tempDir),
		Data:     FileNode{Name: filepath.Base(tempDir), Path: tempDir, IsDir: true},
		Children: nil,
	}

	// This should handle the circular reference gracefully
	err = buildTree(root, tempDir)
	if err != nil {
		t.Logf("buildTree handled circular reference with error (expected): %v", err)
	} else {
		t.Log("buildTree handled circular reference without error")
	}
}

func TestShowHierarchyWithFileInsteadOfDirectory(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "palantir_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Test ShowHierarchy with a file path instead of directory
	err, hasHierarchy := ShowHierarchy(tempFile.Name(), "")
	if err != nil {
		t.Errorf("ShowHierarchy() error = %v", err)
	}
	// A single file should still return true for hierarchy (as per current logic)
	if !hasHierarchy {
		t.Errorf("ShowHierarchy() hasHierarchy = %v, want true for single file", hasHierarchy)
	}
}

func TestBuildTreeWithBrokenSymlink(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "palantir_broken_symlink_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a broken symlink
	brokenSymlink := filepath.Join(tempDir, "broken_link")
	if err := os.Symlink("/nonexistent/target", brokenSymlink); err != nil {
		t.Logf("Could not create broken symlink (might not be supported): %v", err)
		return
	}

	root := &TreeNode{
		Name:     filepath.Base(tempDir),
		Data:     FileNode{Name: filepath.Base(tempDir), Path: tempDir, IsDir: true},
		Children: nil,
	}

	// This should handle the broken symlink gracefully
	err = buildTree(root, tempDir)
	if err != nil {
		t.Logf("buildTree handled broken symlink with error: %v", err)
	} else {
		t.Log("buildTree handled broken symlink without error")
	}
}

func TestParseYAMLToTree(t *testing.T) {
	tests := []struct {
		name             string
		yamlContent      []byte
		expectedRoot     string
		expectedSections []string
		expectedError    bool
	}{
		{
			name: "Valid YAML with nested structure",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  credentials:
    username: admin
    password: secret
  tables:
    - users
    - posts
    - comments
server:
  host: 0.0.0.0
  port: 8080
  debug: true
`),
			expectedRoot:     "root",
			expectedSections: []string{"database", "server"},
			expectedError:    false,
		},
		{
			name:             "Empty YAML",
			yamlContent:      []byte(""),
			expectedRoot:     "root",
			expectedSections: []string{},
			expectedError:    false,
		},
		{
			name: "Simple key-value pairs",
			yamlContent: []byte(`
name: test
value: 42
enabled: true
`),
			expectedRoot:     "root",
			expectedSections: []string{"name", "value", "enabled"},
			expectedError:    false,
		},
		{
			name: "Invalid YAML",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  invalid: [unclosed array
`),
			expectedRoot:     "",
			expectedSections: []string{},
			expectedError:    true,
		},
		{
			name:             "Nil YAML content",
			yamlContent:      nil,
			expectedRoot:     "root",
			expectedSections: []string{},
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := ParseYAMLToTree(tt.yamlContent)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseYAMLToTree() error = %v", err)
			}

			// Verify root structure
			if root.Name != tt.expectedRoot {
				t.Errorf("Expected root name %q, got %q", tt.expectedRoot, root.Name)
			}

			// Verify we have the expected sections
			if len(root.Children) != len(tt.expectedSections) {
				t.Errorf("Expected %d children, got %d", len(tt.expectedSections), len(root.Children))
			}

			// Verify section names using map for O(1) lookup
			actualSections := make(map[string]bool)
			for _, child := range root.Children {
				actualSections[child.Name] = true
			}

			for _, expected := range tt.expectedSections {
				if !actualSections[expected] {
					t.Errorf("Expected section %q not found", expected)
				}
			}
		})
	}
}

func TestParseYAMLToTreeWithDifferentDataTypes(t *testing.T) {
	tests := []struct {
		name           string
		yamlContent    []byte
		expectedArrays map[string][]string
		expectedError  bool
	}{
		{
			name: "Arrays with different data types",
			yamlContent: []byte(`
string_value: hello
number_value: 42
float_value: 3.14
boolean_value: true
array_of_strings:
  - first
  - second
  - third
array_of_numbers:
  - 1
  - 2
  - 3
array_of_booleans:
  - true
  - false
nested_object:
  level1:
    level2:
      value: deep
`),
			expectedArrays: map[string][]string{
				"array_of_strings":  {"first", "second", "third"},
				"array_of_numbers":  {"1", "2", "3"},
				"array_of_booleans": {"true", "false"},
			},
			expectedError: false,
		},
		{
			name: "Mixed array types",
			yamlContent: []byte(`
mixed_array:
  - string_item
  - 42
  - true
  - 3.14
`),
			expectedArrays: map[string][]string{
				"mixed_array": {"string_item", "42", "true", "3.14"},
			},
			expectedError: false,
		},
		{
			name: "Empty arrays",
			yamlContent: []byte(`
empty_strings: []
empty_numbers: []
`),
			expectedArrays: map[string][]string{
				"empty_strings": {},
				"empty_numbers": {},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := ParseYAMLToTree(tt.yamlContent)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseYAMLToTree() error = %v", err)
			}

			// Test each expected array
			for arrayName, expectedValues := range tt.expectedArrays {
				// Find the array
				var arrayNode *TreeNode
				for _, child := range root.Children {
					if child.Name == arrayName {
						arrayNode = child
						break
					}
				}

				if arrayNode == nil {
					t.Errorf("Array %q not found", arrayName)
					continue
				}

				// Verify array length
				if len(arrayNode.Children) != len(expectedValues) {
					t.Errorf("Expected array %q to have %d children, got %d", arrayName, len(expectedValues), len(arrayNode.Children))
					continue
				}

				// Verify array values
				for i, child := range arrayNode.Children {
					if i < len(expectedValues) && child.Name != expectedValues[i] {
						t.Errorf("Expected array %q item %d to be %q, got %q", arrayName, i, expectedValues[i], child.Name)
					}
				}
			}
		})
	}
}

func TestShowYAMLHierarchy(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent []byte
		expectError bool
	}{
		{
			name: "Valid YAML with nested structure",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  tables:
    - users
    - posts
`),
			expectError: false,
		},
		{
			name: "Simple key-value pairs",
			yamlContent: []byte(`
name: test
value: 42
enabled: true
`),
			expectError: false,
		},
		{
			name:        "Empty YAML",
			yamlContent: []byte(""),
			expectError: false,
		},
		{
			name:        "Nil YAML content",
			yamlContent: nil,
			expectError: false,
		},
		{
			name: "Invalid YAML",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  invalid: [unclosed array
`),
			expectError: true,
		},
		{
			name: "Malformed YAML syntax",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  invalid: [unclosed array
  extra: value
`),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ShowYAMLHierarchy(tt.yamlContent)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ShowYAMLHierarchy() error = %v", err)
				}
			}
		})
	}
}

func TestShowYAMLHierarchyFromFile(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent []byte
		expectError bool
	}{
		{
			name: "Valid YAML file",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  tables:
    - users
    - posts
`),
			expectError: false,
		},
		{
			name: "Simple YAML file",
			yamlContent: []byte(`
name: test
value: 42
enabled: true
`),
			expectError: false,
		},
		{
			name:        "Empty YAML file",
			yamlContent: []byte(""),
			expectError: false,
		},
		{
			name: "Invalid YAML file",
			yamlContent: []byte(`
database:
  host: localhost
  port: 5432
  invalid: [unclosed array
`),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary YAML file
			tempFile, err := os.CreateTemp("", "test_yaml_*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			// Write YAML content to file
			if _, err := tempFile.Write(tt.yamlContent); err != nil {
				t.Fatalf("Failed to write YAML content: %v", err)
			}
			tempFile.Close()

			// Test ShowYAMLHierarchyFromFile
			err = ShowYAMLHierarchyFromFile(tempFile.Name())

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ShowYAMLHierarchyFromFile() error = %v", err)
				}
			}
		})
	}

	// Test with non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		err := ShowYAMLHierarchyFromFile("/nonexistent/file.yaml")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
}

func TestYAMLNodeDataTypes(t *testing.T) {
	// Test YAML content with mixed data types
	yamlContent := []byte(`
database:
  host: localhost
  port: 5432
  tables:
    - users
    - posts
server:
  host: 0.0.0.0
  port: 8080
`)

	root, err := ParseYAMLToTree(yamlContent)
	if err != nil {
		t.Fatalf("ParseYAMLToTree() error = %v", err)
	}

	// Test root node
	if root.Name != "root" {
		t.Errorf("Expected root name 'root', got %q", root.Name)
	}

	// Find database section
	var database *TreeNode
	for _, child := range root.Children {
		if child.Name == "database" {
			database = child
			break
		}
	}

	if database == nil {
		t.Fatal("database section not found")
	}

	// Verify database YAMLNode data
	if yamlNode, ok := database.Data.(YAMLNode); ok {
		if yamlNode.Name != "database" {
			t.Errorf("Expected YAMLNode name 'database', got %q", yamlNode.Name)
		}
		if !yamlNode.IsDir {
			t.Error("Expected YAMLNode IsDir to be true for object")
		}
		if yamlNode.NodeType != "object" {
			t.Errorf("Expected YAMLNode NodeType 'object', got %q", yamlNode.NodeType)
		}
	} else {
		t.Error("Expected YAMLNode data type")
	}

	// Find tables array
	var tables *TreeNode
	for _, child := range database.Children {
		if child.Name == "tables" {
			tables = child
			break
		}
	}

	if tables == nil {
		t.Fatal("tables array not found")
	}

	// Verify array YAMLNode data
	if yamlNode, ok := tables.Data.(YAMLNode); ok {
		if yamlNode.Name != "tables" {
			t.Errorf("Expected YAMLNode name 'tables', got %q", yamlNode.Name)
		}
		if !yamlNode.IsDir {
			t.Error("Expected YAMLNode IsDir to be true for array")
		}
		if yamlNode.NodeType != "object" {
			t.Errorf("Expected YAMLNode NodeType 'object', got %q", yamlNode.NodeType)
		}
	} else {
		t.Error("Expected YAMLNode data type for array")
	}

	// Find first table item
	if len(tables.Children) == 0 {
		t.Fatal("Expected at least one table item")
	}

	firstTable := tables.Children[0]
	if yamlNode, ok := firstTable.Data.(YAMLNode); ok {
		if yamlNode.Name != "users" {
			t.Errorf("Expected YAMLNode name 'users', got %q", yamlNode.Name)
		}
		if yamlNode.IsDir {
			t.Error("Expected YAMLNode IsDir to be false for array item")
		}
		if yamlNode.NodeType != "array" {
			t.Errorf("Expected YAMLNode NodeType 'array', got %q", yamlNode.NodeType)
		}
	} else {
		t.Error("Expected YAMLNode data type for array item")
	}
}
