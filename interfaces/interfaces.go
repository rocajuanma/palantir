package interfaces

// OutputHandler defines the interface for terminal output operations
type OutputHandler interface {
	PrintHeader(message string)
	PrintStage(message string)
	PrintSuccess(message string)
	PrintError(format string, args ...interface{})
	PrintWarning(format string, args ...interface{})
	PrintInfo(format string, args ...interface{})
	PrintAlreadyAvailable(format string, args ...interface{})
	PrintProgress(current, total int, message string)
	Confirm(message string) bool
	IsSupported() bool
	Disable()
}
