package terminal

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func setupSupportedTerminal(t *testing.T) {
	oldTerm := os.Getenv("TERM")
	os.Setenv("TERM", "xterm-256color")
	t.Cleanup(func() {
		os.Setenv("TERM", oldTerm)
	})
}

func setupUnsupportedTerminal(t *testing.T) {
	oldTerm := os.Getenv("TERM")
	os.Setenv("TERM", "dumb")
	t.Cleanup(func() {
		os.Setenv("TERM", oldTerm)
	})
}

func TestFormatMessage_AllConfigurations(t *testing.T) {
	setupSupportedTerminal(t)

	configs := []struct {
		name     string
		config   *OutputConfig
		expected map[OutputLevel]string
	}{
		{
			"WithAllFeatures",
			&OutputConfig{UseColors: true, UseEmojis: true, UseFormatting: true, DisableOutput: false},
			map[OutputLevel]string{
				LevelHeader:  fmt.Sprintf("\n%s%s=== Test Header ===%s\n", ColorBold, ColorCyan, ColorReset),
				LevelStage:   fmt.Sprintf("%s%s🔧 Test Stage%s\n", ColorBold, ColorBlue, ColorReset),
				LevelSuccess: fmt.Sprintf("%s%s✅ Test Success%s\n", ColorBold, ColorGreen, ColorReset),
				LevelError:   fmt.Sprintf("%s%s❌ Test Error%s\n", ColorBold, ColorRed, ColorReset),
				LevelWarning: fmt.Sprintf("%s%s⚠️  Test Warning%s\n", ColorBold, ColorYellow, ColorReset),
				LevelInfo:    fmt.Sprintf("%s%sTest Info%s\n", ColorBold, "", ColorReset),
			},
		},
		{
			"WithLevelOnlyColours",
			&OutputConfig{UseColors: true, UseEmojis: true, UseFormatting: true, DisableOutput: false, ColorizeLevelOnly: true},
			map[OutputLevel]string{
				LevelHeader:  fmt.Sprintf("\n%s%s=== Test Header ===%s\n", ColorBold, ColorCyan, ColorReset),
				LevelStage:   fmt.Sprintf("%s%s🔧 %sTest Stage\n", ColorBold, ColorBlue, ColorReset),
				LevelSuccess: fmt.Sprintf("%s%s✅ %sTest Success\n", ColorBold, ColorGreen, ColorReset),
				LevelError:   fmt.Sprintf("%s%s❌ %sTest Error\n", ColorBold, ColorRed, ColorReset),
				LevelWarning: fmt.Sprintf("%s%s⚠️  %sTest Warning\n", ColorBold, ColorYellow, ColorReset),
				LevelInfo:    fmt.Sprintf("%sTest Info%s\n", ColorBold, ColorReset),
			},
		},
		{
			"WithColorsOnly",
			&OutputConfig{UseColors: true, UseEmojis: false, UseFormatting: true, DisableOutput: false},
			map[OutputLevel]string{
				LevelHeader:  fmt.Sprintf("\n%s%s=== Test Header ===%s\n", ColorBold, ColorCyan, ColorReset),
				LevelStage:   fmt.Sprintf("%s%s[STAGE] Test Stage%s\n", ColorBold, ColorBlue, ColorReset),
				LevelSuccess: fmt.Sprintf("%s%s[SUCCESS] Test Success%s\n", ColorBold, ColorGreen, ColorReset),
				LevelError:   fmt.Sprintf("%s%s[ERROR] Test Error%s\n", ColorBold, ColorRed, ColorReset),
				LevelWarning: fmt.Sprintf("%s%s[WARNING] Test Warning%s\n", ColorBold, ColorYellow, ColorReset),
				LevelInfo:    fmt.Sprintf("%s%sTest Info%s\n", ColorBold, "", ColorReset),
			},
		},
		{
			"WithoutColors",
			&OutputConfig{UseColors: false, UseEmojis: false, UseFormatting: false, DisableOutput: false},
			map[OutputLevel]string{
				LevelHeader:  "\n=== Test Header ===\n",
				LevelStage:   "[STAGE] Test Stage\n",
				LevelSuccess: "[SUCCESS] Test Success\n",
				LevelError:   "[ERROR] Test Error\n",
				LevelWarning: "[WARNING] Test Warning\n",
				LevelInfo:    "Test Info\n",
			},
		},
	}

	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			handler := NewOutputHandler(config.config)

			for level := range config.expected {
				t.Run(fmt.Sprintf("Level_%d", level), func(t *testing.T) {
					message := fmt.Sprintf("Test %s", levelNames[level])
					result := handler.FormatMessage(level, message)
					expectedMsg := generateExpectedOutput(level, message, config.config)

					if result != expectedMsg {
						t.Errorf("FormatMessage() = %q, want %q", result, expectedMsg)
					}
				})
			}
		})
	}
}

// Helper map for level names
var levelNames = map[OutputLevel]string{
	LevelHeader:  "Header",
	LevelStage:   "Stage",
	LevelSuccess: "Success",
	LevelError:   "Error",
	LevelWarning: "Warning",
	LevelInfo:    "Info",
}

// generateExpectedOutput is a helper function to generate expected output for FormatMessage
func generateExpectedOutput(level OutputLevel, message string, config *OutputConfig) string {
	if config.DisableOutput {
		return ""
	}

	// Handle unsupported terminal case
	if os.Getenv("TERM") == "dumb" {
		return message
	}

	if level == LevelHeader {
		if config.UseColors {
			color := outputColors[level]
			return fmt.Sprintf(coloredHeaderFormat, ColorBold, color, message, ColorReset)
		}
		return fmt.Sprintf(headerFormat, message)
	}

	var prefix string
	var color string

	if config.UseColors && config.UseEmojis && config.UseFormatting {
		prefix = outputEmojis[level]
		color = outputColors[level]
	} else {
		prefix = outputPrefixes[level]
		if config.UseColors {
			color = outputColors[level]
		}
	}

	if config.UseColors && config.UseFormatting {
		if config.ColorizeLevelOnly && color != "" && prefix != "" {
			coloredPrefix := fmt.Sprintf("%s%s%s%s", ColorBold, color, prefix, ColorReset)
			return fmt.Sprintf("%s%s\n", coloredPrefix, message)
		}
		return fmt.Sprintf("%s%s%s%s%s\n", ColorBold, color, prefix, message, ColorReset)
	}

	return fmt.Sprintf("%s%s\n", prefix, message)
}

func TestFormatMessage_EdgeCases(t *testing.T) {
	// Test disabled output
	handler := &outputHandler{
		config: &OutputConfig{DisableOutput: true},
	}
	result := handler.FormatMessage(LevelInfo, "Test Message")
	if result != "" {
		t.Errorf("FormatMessage() with disabled output = %q, want empty string", result)
	}

	// Test unsupported terminal
	setupUnsupportedTerminal(t)
	handler = &outputHandler{
		config: &OutputConfig{
			UseColors:     true,
			UseEmojis:     true,
			UseFormatting: true,
			DisableOutput: false,
		},
	}
	result = handler.FormatMessage(LevelInfo, "Test Message")
	expected := "Test Message"
	if result != expected {
		t.Errorf("FormatMessage() with unsupported terminal = %q, want %q", result, expected)
	}
}

func TestPrintMethods_AllVariations(t *testing.T) {
	setupSupportedTerminal(t)

	// Test basic print methods
	handler := &outputHandler{
		config: &OutputConfig{
			UseColors:     true,
			UseEmojis:     true,
			UseFormatting: true,
			DisableOutput: false,
		},
	}

	tests := []struct {
		name     string
		method   func(string)
		message  string
		expected string
	}{
		{
			"PrintHeader",
			handler.PrintHeader,
			"Test Header",
			fmt.Sprintf("\n%s%s=== Test Header ===%s\n", ColorBold, ColorCyan, ColorReset),
		},
		{
			"PrintStage",
			handler.PrintStage,
			"Test Stage",
			fmt.Sprintf("%s%s🔧 Test Stage%s\n", ColorBold, ColorBlue, ColorReset),
		},
		{
			"PrintSuccess",
			handler.PrintSuccess,
			"Test Success",
			fmt.Sprintf("%s%s✅ Test Success%s\n", ColorBold, ColorGreen, ColorReset),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.method(tt.message)
			})
			if output != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, output, tt.expected)
			}
		})
	}

	formatTests := []struct {
		name     string
		method   func(string, ...interface{})
		format   string
		args     []interface{}
		expected string
	}{
		{
			"PrintError",
			handler.PrintError,
			"Error: %s",
			[]interface{}{"test error"},
			fmt.Sprintf("%s%s❌ Error: test error%s\n", ColorBold, ColorRed, ColorReset),
		},
		{
			"PrintWarning",
			handler.PrintWarning,
			"Warning: %s",
			[]interface{}{"test warning"},
			fmt.Sprintf("%s%s⚠️  Warning: test warning%s\n", ColorBold, ColorYellow, ColorReset),
		},
		{
			"PrintInfo",
			handler.PrintInfo,
			"Info: %s",
			[]interface{}{"test info"},
			fmt.Sprintf("%s%sInfo: test info%s\n", ColorBold, "", ColorReset),
		},
	}

	for _, tt := range formatTests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.method(tt.format, tt.args...)
			})
			if output != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, output, tt.expected)
			}
		})
	}
}

func TestPrintAlreadyAvailable_AllConfigurations(t *testing.T) {
	setupSupportedTerminal(t)

	configs := []struct {
		name     string
		config   *OutputConfig
		expected string
	}{
		{
			"WithColoursAndEmojis",
			&OutputConfig{UseColors: true, UseEmojis: true, UseFormatting: true, DisableOutput: false},
			fmt.Sprintf("%s%s💙 Feature is available%s\n", ColorBold, ColorBlue, ColorReset),
		},
		{
			"WithColoursAndEmojis_LevelOnly",
			&OutputConfig{UseColors: true, UseEmojis: true, UseFormatting: true, DisableOutput: false, ColorizeLevelOnly: true},
			fmt.Sprintf("%s%s💙 %sFeature is available\n", ColorBold, ColorBlue, ColorReset),
		},
		{
			"WithColours",
			&OutputConfig{UseColors: true, UseEmojis: false, UseFormatting: true, DisableOutput: false},
			fmt.Sprintf("%s%s[AVAILABLE] Feature is available%s\n", ColorBold, ColorBlue, ColorReset),
		},
		{
			"WithColours_LevelOnly",
			&OutputConfig{UseColors: true, UseEmojis: false, UseFormatting: true, DisableOutput: false, ColorizeLevelOnly: true},
			fmt.Sprintf("%s%s[AVAILABLE] %sFeature is available\n", ColorBold, ColorBlue, ColorReset),
		},
		{
			"WithEmojisAndNoColours", // TODO: Currently not supported, emojis are supported only when colours are enabled
			&OutputConfig{UseColors: false, UseEmojis: true, UseFormatting: true, DisableOutput: false},
			"[AVAILABLE] Feature is available\n",
		},
		{
			"WithoutColoursAndEmojis",
			&OutputConfig{UseColors: false, UseEmojis: false, UseFormatting: false, DisableOutput: false},
			"[AVAILABLE] Feature is available\n",
		},
	}

	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			handler := NewOutputHandler(config.config)

			output := captureOutput(func() {
				handler.PrintAlreadyAvailable("Feature is available")
			})
			if output != config.expected {
				t.Errorf("PrintAlreadyAvailable() = %q, want %q", output, config.expected)
			}
		})
	}
}

func TestPrintProgress_AllScenarios(t *testing.T) {
	setupSupportedTerminal(t)

	t.Run("WithColors", func(t *testing.T) {
		handler := NewOutputHandler(&OutputConfig{
			UseColors:     true,
			UseEmojis:     true,
			UseFormatting: true,
			DisableOutput: false,
		})

		output := captureOutput(func() {
			handler.PrintProgress(3, 10, "Processing")
		})
		expected := fmt.Sprintf("\r%s%s[3/10] 30%% - Processing%s\n", ColorBold, ColorCyan, ColorReset)
		if output != expected {
			t.Errorf("PrintProgress() = %q, want %q", output, expected)
		}
	})

	t.Run("WithColorsLevelOnly", func(t *testing.T) {
		handler := NewOutputHandler(&OutputConfig{
			UseColors:         true,
			UseEmojis:         true,
			UseFormatting:     true,
			DisableOutput:     false,
			ColorizeLevelOnly: true,
		})

		output := captureOutput(func() {
			handler.PrintProgress(3, 10, "Processing")
		})
		expected := fmt.Sprintf("\r%s%s[3/10] 30%% - %sProcessing\n", ColorBold, ColorCyan, ColorReset)
		if output != expected {
			t.Errorf("PrintProgress() level-only = %q, want %q", output, expected)
		}
	})

	t.Run("WithoutColors", func(t *testing.T) {
		handler := NewOutputHandler(&OutputConfig{
			UseColors:     false,
			UseEmojis:     false,
			UseFormatting: false,
			DisableOutput: false,
		})

		output := captureOutput(func() {
			handler.PrintProgress(3, 10, "Processing")
		})
		expected := "\r[3/10] 30% - Processing\n"
		if output != expected {
			t.Errorf("PrintProgress() = %q, want %q", output, expected)
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		handler := NewOutputHandler(&OutputConfig{
			UseColors:     false,
			UseEmojis:     false,
			UseFormatting: false,
			DisableOutput: false,
		})

		tests := []struct {
			current  int
			total    int
			message  string
			expected string
		}{
			{0, 10, "Starting", "\r[0/10] 0% - Starting\n"},
			{10, 10, "Complete", "\r[10/10] 100% - Complete\n"},
			{1, 3, "One third", "\r[1/3] 33% - One third\n"},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%d_%d", tt.current, tt.total), func(t *testing.T) {
				output := captureOutput(func() {
					handler.PrintProgress(tt.current, tt.total, tt.message)
				})
				if output != tt.expected {
					t.Errorf("PrintProgress(%d, %d, %q) = %q, want %q", tt.current, tt.total, tt.message, output, tt.expected)
				}
			})
		}
	})
}

func TestDisabledOutput(t *testing.T) {
	handler := NewOutputHandler(&OutputConfig{
		DisableOutput: true,
	})

	// Test that all print methods return nothing when disabled
	methods := []func(){
		func() { handler.PrintHeader("test") },
		func() { handler.PrintStage("test") },
		func() { handler.PrintSuccess("test") },
		func() { handler.PrintError("test") },
		func() { handler.PrintWarning("test") },
		func() { handler.PrintInfo("test") },
		func() { handler.PrintAlreadyAvailable("test") },
		func() { handler.PrintProgress(1, 2, "test") },
	}

	for i, method := range methods {
		output := captureOutput(method)
		if output != "" {
			t.Errorf("Method %d should return empty string when disabled, got %q", i, output)
		}
	}
}

func TestIsSupported(t *testing.T) {
	handler := &outputHandler{}

	setupSupportedTerminal(t)
	if !handler.IsSupported() {
		t.Error("IsSupported() should return true for normal terminal")
	}

	// Test with dumb terminal
	os.Setenv("TERM", "dumb")
	if handler.IsSupported() {
		t.Error("IsSupported() should return false for dumb terminal")
	}
}

func TestGlobalHandler(t *testing.T) {
	handler := GetGlobalOutputHandler()
	if handler == nil {
		t.Error("GetGlobalOutputHandler() should not return nil")
	}

	customHandler := NewOutputHandler(&OutputConfig{
		UseColors:     false,
		UseEmojis:     false,
		UseFormatting: false,
		DisableOutput: false,
	})

	SetGlobalOutputHandler(customHandler)
	retrieved := GetGlobalOutputHandler()
	if retrieved != customHandler {
		t.Error("SetGlobalOutputHandler() should set the global handler")
	}

	SetGlobalOutputHandler(NewDefaultOutputHandler())
}

func TestConfirm_AllScenarios(t *testing.T) {
	setupSupportedTerminal(t)

	handler := NewOutputHandler(&OutputConfig{
		UseColors:     true,
		UseEmojis:     true,
		UseFormatting: true,
		DisableOutput: false,
	})

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Yes_lowercase", "y", true},
		{"Yes_uppercase", "Y", true},
		{"Yes_word", "yes", true},
		{"Yes_word_capitalized", "Yes", true},
		{"No_lowercase", "n", false},
		{"No_uppercase", "N", false},
		{"No_word", "no", false},
		{"Empty_input", "", false},
		{"Invalid_input", "maybe", false},
		{"Partial_yes", "ye", false},
		{"Partial_no", "na", false},
		{"random_word", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original stdin
			oldStdin := os.Stdin
			defer func() {
				os.Stdin = oldStdin
			}()

			// Create a pipe to simulate stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write the test input
			go func() {
				w.WriteString(tt.input + "\n")
				w.Close()
			}()

			result := handler.Confirm("Test confirmation")
			if result != tt.expected {
				t.Errorf("Confirm() with input %q = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}

	// Test disabled output
	handler.Disable()
	result := handler.Confirm("Test confirmation")
	if result != false {
		t.Error("Confirm() should return false when output is disabled")
	}
}

func TestConfirm_LevelOnlyColours(t *testing.T) {
	setupSupportedTerminal(t)

	handler := NewOutputHandler(&OutputConfig{
		UseColors:         true,
		UseEmojis:         true,
		UseFormatting:     true,
		DisableOutput:     false,
		ColorizeLevelOnly: true,
	})

	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("y\n")
		w.Close()
	}()

	output := captureOutput(func() {
		handler.Confirm("Test confirmation")
	})

	expected := fmt.Sprintf("%s%s?%s Test confirmation (y/N): ", ColorBold, ColorYellow, ColorReset)
	if output != expected {
		t.Errorf("Confirm() level-only output = %q, want %q", output, expected)
	}
}

func TestPrintProgress_ExtendedEdgeCases(t *testing.T) {
	setupSupportedTerminal(t)

	handler := NewOutputHandler(&OutputConfig{
		UseColors:     false,
		UseEmojis:     false,
		UseFormatting: false,
		DisableOutput: false,
	})

	tests := []struct {
		name     string
		current  int
		total    int
		message  string
		expected string
	}{
		{name: "Zero_progress", current: 0, total: 10, message: "Starting", expected: "\r[0/10] 0% - Starting\n"},
		{name: "Complete_progress", current: 10, total: 10, message: "Complete", expected: "\r[10/10] 100% - Complete\n"},
		{name: "Half_progress", current: 5, total: 10, message: "Halfway", expected: "\r[5/10] 50% - Halfway\n"},
		{name: "Single_item", current: 1, total: 1, message: "One item", expected: "\r[1/1] 100% - One item\n"},
		{name: "Large_numbers", current: 999, total: 1000, message: "Almost done", expected: "\r[999/1000] 100% - Almost done\n"},
		{name: "Fractional_percentage", current: 1, total: 3, message: "One third", expected: "\r[1/3] 33% - One third\n"},
		{name: "Small_fraction", current: 1, total: 7, message: "Small fraction", expected: "\r[1/7] 14% - Small fraction\n"},
		{name: "Zero_total", current: 0, total: 0, message: "Zero total", expected: "\r[0/0] NaN% - Zero total\n"},
		{name: "Negative_current", current: -1, total: 10, message: "Negative", expected: "\r[-1/10] -10% - Negative\n"},
		{name: "Current_greater_than_total", current: 15, total: 10, message: "Overflow", expected: "\r[15/10] 150% - Overflow\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				handler.PrintProgress(tt.current, tt.total, tt.message)
			})
			if output != tt.expected {
				t.Errorf("PrintProgress(%d, %d, %q) = %q, want %q", tt.current, tt.total, tt.message, output, tt.expected)
			}
		})
	}

	t.Run("WithColors", func(t *testing.T) {
		coloredHandler := NewOutputHandler(&OutputConfig{
			UseColors:     true,
			UseEmojis:     true,
			UseFormatting: true,
			DisableOutput: false,
		})

		output := captureOutput(func() {
			coloredHandler.PrintProgress(3, 10, "Colored progress")
		})
		expected := fmt.Sprintf("\r%s%s[3/10] 30%% - Colored progress%s\n", ColorBold, ColorCyan, ColorReset)
		if output != expected {
			t.Errorf("PrintProgress with colors = %q, want %q", output, expected)
		}
	})

	t.Run("WithColorsLevelOnly", func(t *testing.T) {
		coloredHandler := NewOutputHandler(&OutputConfig{
			UseColors:         true,
			UseEmojis:         true,
			UseFormatting:     true,
			DisableOutput:     false,
			ColorizeLevelOnly: true,
		})

		output := captureOutput(func() {
			coloredHandler.PrintProgress(3, 10, "Colored progress")
		})
		expected := fmt.Sprintf("\r%s%s[3/10] 30%% - %sColored progress\n", ColorBold, ColorCyan, ColorReset)
		if output != expected {
			t.Errorf("PrintProgress with level-only colors = %q, want %q", output, expected)
		}
	})

	// Test disabled output
	t.Run("DisabledOutput", func(t *testing.T) {
		handler.Disable()

		output := captureOutput(func() {
			handler.PrintProgress(5, 10, "Should not appear")
		})
		if output != "" {
			t.Errorf("PrintProgress with disabled output = %q, want empty string", output)
		}
	})
}

func TestOutputFormatConsistency(t *testing.T) {
	setupSupportedTerminal(t)

	handler := NewOutputHandler(&OutputConfig{
		UseColors:     true,
		UseEmojis:     true,
		UseFormatting: true,
		DisableOutput: false,
	})

	message := "Test Message"
	output1 := handler.FormatMessage(LevelSuccess, message)
	output2 := handler.FormatMessage(LevelSuccess, message)

	// Output should be deterministic for the same input and config.
	if output1 != output2 {
		t.Error("FormatMessage should produce consistent output for the same input and configuration")
	}

	// Output should include the success emoji, the message, and a newline.
	if !strings.Contains(output1, "✅") {
		t.Error("Success output should include the success emoji (✅)")
	}
	if !strings.Contains(output1, message) {
		t.Error("Output should include the original message")
	}
	if !strings.HasSuffix(output1, "\n") {
		t.Error("Output should end with a newline character")
	}
}
