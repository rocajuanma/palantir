package terminal

import (
	"fmt"
	"os"

	"github.com/rocajuanma/palantir/interfaces"
)

// OutputLevel represents different levels of output
type OutputLevel int

const (
	LevelInfo OutputLevel = iota
	LevelWarning
	LevelError
	LevelSuccess
	LevelStage
	LevelHeader
)

// OutputConfig holds configuration for output formatting
type OutputConfig struct {
	UseColors     bool
	UseEmojis     bool
	UseFormatting bool
	DisableOutput bool
	VerboseMode   bool
}

// outputHandler implements the OutputHandler interface
type outputHandler struct {
	config *OutputConfig
}

// NewDefaultOutputHandler creates a new outputHandler with default configurations
func NewDefaultOutputHandler() interfaces.OutputHandler {
	return &outputHandler{
		config: &OutputConfig{
			UseColors:     true,
			UseEmojis:     true,
			UseFormatting: true,
			DisableOutput: false,
			VerboseMode:   false,
		},
	}
}

// NewOutputHandler creates a new outputHandler with a custom configurations
func NewOutputHandler(config *OutputConfig) *outputHandler {
	return &outputHandler{config: config}
}

// FormatMessage formats a message according to the output level
func (oh *outputHandler) FormatMessage(level OutputLevel, message string) string {
	if oh.config.DisableOutput {
		return ""
	}

	if !oh.IsSupported() {
		return message
	}

	var prefix, color string

	if oh.config.UseColors && oh.config.UseEmojis && oh.config.UseFormatting {
		prefix = outputEmojis[level]
		color = outputColors[level]
		if level == LevelHeader {
			return fmt.Sprintf(coloredHeaderFormat, ColorBold, color, message, ColorReset)
		}
	} else if oh.config.UseColors {
		prefix = outputPrefixes[level]
		color = outputColors[level]
		if level == LevelHeader {
			return fmt.Sprintf(coloredHeaderFormat, ColorBold, color, message, ColorReset)
		}
	} else {
		prefix = outputPrefixes[level]
		if level == LevelHeader {
			return fmt.Sprintf(prefix, message)
		}
	}

	if oh.config.UseColors && oh.config.UseFormatting {
		return fmt.Sprintf("%s%s%s%s%s\n", ColorBold, color, prefix, message, ColorReset)
	}

	return fmt.Sprintf("%s%s\n", prefix, message)
}

// PrintWithLevel prints a message with the specified level
func (oh *outputHandler) PrintWithLevel(level OutputLevel, format string, args ...interface{}) {
	if oh.config.DisableOutput {
		return
	}

	message := fmt.Sprintf(format, args...)
	formatted := oh.FormatMessage(level, message)
	fmt.Print(formatted)
}

// Implementation of OutputHandler interface methods

func (oh *outputHandler) PrintHeader(message string) {
	oh.PrintWithLevel(LevelHeader, message)
}

func (oh *outputHandler) PrintStage(message string) {
	oh.PrintWithLevel(LevelStage, message)
}

func (oh *outputHandler) PrintSuccess(message string) {
	oh.PrintWithLevel(LevelSuccess, message)
}

func (oh *outputHandler) PrintError(format string, args ...interface{}) {
	oh.PrintWithLevel(LevelError, format, args...)
}

func (oh *outputHandler) PrintWarning(format string, args ...interface{}) {
	oh.PrintWithLevel(LevelWarning, format, args...)
}

func (oh *outputHandler) PrintInfo(format string, args ...interface{}) {
	oh.PrintWithLevel(LevelInfo, format, args...)
}

func (oh *outputHandler) PrintAlreadyAvailable(format string, args ...interface{}) {
	if oh.config.DisableOutput {
		return
	}

	message := fmt.Sprintf(format, args...)

	if oh.config.UseColors && oh.config.UseEmojis && oh.config.UseFormatting {
		fmt.Printf("%s%s💙 %s%s\n", ColorBold, ColorBlue, message, ColorReset)
	} else if oh.config.UseColors {
		fmt.Printf("%s%s[AVAILABLE] %s%s\n", ColorBold, ColorBlue, message, ColorReset)
	} else {
		fmt.Printf("[AVAILABLE] %s\n", message)
	}
}

func (oh *outputHandler) PrintProgress(current, total int, message string) {
	if oh.config.DisableOutput {
		return
	}

	percentage := float64(current) / float64(total) * 100

	if oh.config.UseColors && oh.config.UseFormatting {
		fmt.Printf("\r%s%s[%d/%d] %.0f%% - %s%s\n", ColorBold, ColorCyan, current, total, percentage, message, ColorReset)
	} else {
		fmt.Printf("\r[%d/%d] %.0f%% - %s\n", current, total, percentage, message)
	}
}

func (oh *outputHandler) Confirm(message string) bool {
	if oh.config.DisableOutput {
		return false
	}

	if oh.config.UseColors && oh.config.UseFormatting {
		fmt.Printf("%s%s? %s (y/N): %s", ColorBold, ColorYellow, message, ColorReset)
	} else {
		fmt.Printf("? %s (y/N): ", message)
	}

	var response string
	fmt.Scanln(&response)

	switch response {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}

func (oh *outputHandler) IsSupported() bool {
	return os.Getenv("TERM") != "dumb"
}

// Disable disables all output
func (oh *outputHandler) Disable() {
	oh.config.DisableOutput = true
}

// Global output handler instance
var globalOutputHandler interfaces.OutputHandler = NewDefaultOutputHandler()

// SetGlobalOutputHandler sets the global output handler
func SetGlobalOutputHandler(handler interfaces.OutputHandler) {
	globalOutputHandler = handler
}

// GetGlobalOutputHandler returns the global output handler
func GetGlobalOutputHandler() interfaces.OutputHandler {
	if globalOutputHandler == nil {
		globalOutputHandler = NewDefaultOutputHandler()
	}
	return globalOutputHandler
}
