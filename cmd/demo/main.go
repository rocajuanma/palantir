package main

import (
	"github.com/rocajuanma/palantir"
)

func main() {
	// Initialize the default output handler(all features enabled)
	handler := palantir.NewDefaultOutputHandler()

	// Showcases the different output levels when default configurations are used.
	handler.PrintHeader("Palantir Demo(Default)")
	handler.PrintInfo("This is an info message")
	handler.PrintSuccess("Operation completed successfully!")
	handler.PrintWarning("This is a warning message")
	handler.PrintError("This is an error message")
	handler.PrintStage("Processing stage 1")
	handler.PrintAlreadyAvailable("Feature is already available")
	handler.PrintProgress(3, 10, "Processing items")

	// Tests the user's confirmation and success/failure scenarios
	if handler.Confirm("Do you want to continue?") {
		handler.PrintSuccess("User confirmed!")
	} else {
		handler.PrintInfo("User declined")
	}

	// Setup configuration that only colours the output level indicator
	levelColoursConfig := &palantir.OutputConfig{
		UseColors:         true,
		UseEmojis:         false,
		UseFormatting:     true,
		DisableOutput:     false,
		ColorizeLevelOnly: true,
	}

	levelColours := palantir.NewOutputHandler(levelColoursConfig)
	levelColours.PrintHeader("Palantir Demo(Level Colours Only)")
	levelColours.PrintInfo("This is an info message")
	levelColours.PrintSuccess("Operation completed successfully!")
	levelColours.PrintWarning("This is a warning message")
	levelColours.PrintError("This is an error message")
	levelColours.PrintStage("Processing stage 1")
	levelColours.PrintAlreadyAvailable("Feature is already available")
	levelColours.PrintProgress(3, 10, "Processing items")
	if levelColours.Confirm("Do you want to continue?") {
		levelColours.PrintSuccess("User confirmed!")
	} else {
		levelColours.PrintInfo("User declined")
	}

	// Setup configurations with colours only
	coloursOnlyConfig := &palantir.OutputConfig{
		UseColors:     true,
		UseEmojis:     false,
		UseFormatting: true,
		DisableOutput: false,
	}

	onlyColours := palantir.NewOutputHandler(coloursOnlyConfig)
	onlyColours.PrintHeader("Palantir Demo(Colours Only)")
	onlyColours.PrintInfo("This is an info message")
	onlyColours.PrintSuccess("Operation completed successfully!")
	onlyColours.PrintWarning("This is a warning message")
	onlyColours.PrintError("This is an error message")
	onlyColours.PrintStage("Processing stage 1")
	onlyColours.PrintAlreadyAvailable("Feature is already available")
	onlyColours.PrintProgress(3, 10, "Processing items")
	if onlyColours.Confirm("Do you want to continue?") {
		onlyColours.PrintSuccess("User confirmed!")
	} else {
		onlyColours.PrintInfo("User declined")
	}

	// Setup configurations without colours
	noColoursConfig := &palantir.OutputConfig{
		UseColors:     false,
		UseEmojis:     false,
		UseFormatting: false,
		DisableOutput: false,
	}
	noColours := palantir.NewOutputHandler(noColoursConfig)
	noColours.PrintHeader("Palantir Demo(Without Colours)")
	noColours.PrintInfo("This is an info message")
	noColours.PrintSuccess("Operation completed successfully!")
	noColours.PrintWarning("This is a warning message")
	noColours.PrintError("This is an error message")
	noColours.PrintStage("Processing stage 1")
	noColours.PrintAlreadyAvailable("Feature is already available")
	noColours.PrintProgress(3, 10, "Processing items")
	if noColours.Confirm("Do you want to continue?") {
		noColours.PrintSuccess("User confirmed!")
	} else {
		noColours.PrintInfo("User declined")
	}

	// File Tree demo
	handler.PrintHeader("File/Directory Tree Visualization")
	handler.PrintStage("Tree with different output configurations")

	// Show tree of current directory
	handler.PrintInfo("Displaying tree structure of current directory:")
	err, hasHierarchy := palantir.ShowHierarchy(".", "")
	if err != nil {
		handler.PrintError("Failed to display tree: %v", err)
	} else if hasHierarchy {
		handler.PrintSuccess("Tree displayed successfully!\n")
	} else {
		handler.PrintInfo("No hierarchy to display (single file)")
	}

	// Tree with colors disabled
	handler.PrintInfo("Tree with colours disabled:")
	palantir.SetGlobalOutputHandler(noColours)
	err, hasHierarchy = palantir.ShowHierarchy(".", "")
	if err != nil {
		handler.PrintError("Failed to display tree: %v", err)
	}

	// Sample YAML content
	yamlContent := []byte(`
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
  features:
    - authentication
    - logging
    - monitoring
redis:
  host: redis-server
  port: 6379
  database: 0
`)

	// YAML Tree demo
	handler.PrintHeader("YAML Tree Visualization")
	err = palantir.ShowYAMLHierarchy(yamlContent)
	if err != nil {
		handler.PrintError("Failed to display YAML tree: %v", err)
	} else {
		handler.PrintSuccess("YAML tree displayed successfully!\n")
	}

	handler.PrintSuccess("Tree system demonstration completed!")
}
