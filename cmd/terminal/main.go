package main

import (
	"github.com/rocajuanma/palantir/terminal"
)

func main() {
	// Initialize the global output handler
	handler := terminal.GetGlobalOutputHandler()

	// Showcases the different output levels
	handler.PrintHeader("Palantir Terminal Output Demo")
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
}
