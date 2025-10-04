package palantir

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"  // Reset all attributes
	ColorRed    = "\033[31m" // Red foreground
	ColorGreen  = "\033[32m" // Green foreground
	ColorYellow = "\033[33m" // Yellow foreground
	ColorBlue   = "\033[34m" // Blue foreground
	ColorPurple = "\033[35m" // Magenta (sometimes called purple) foreground
	ColorCyan   = "\033[36m" // Cyan foreground
	ColorWhite  = "\033[37m" // White foreground
	ColorBold   = "\033[1m"  // Bold text
)

var (
	// outputColors is a map of output levels to their corresponding colors
	outputColors = map[OutputLevel]string{
		LevelHeader:  ColorCyan,
		LevelStage:   ColorBlue,
		LevelSuccess: ColorGreen,
		LevelError:   ColorRed,
		LevelWarning: ColorYellow,
		LevelInfo:    "",
	}

	// outputEmojis is a map of output levels to their corresponding emojis
	outputEmojis = map[OutputLevel]string{
		LevelHeader:  "",
		LevelStage:   "🔧 ",
		LevelSuccess: "✅ ",
		LevelError:   "❌ ",
		LevelWarning: "⚠️  ",
		LevelInfo:    "",
	}

	// outputPrefixes is a map of output levels to their corresponding prefixes
	outputPrefixes = map[OutputLevel]string{
		LevelHeader:  headerFormat,
		LevelStage:   "[STAGE] ",
		LevelSuccess: "[SUCCESS] ",
		LevelError:   "[ERROR] ",
		LevelWarning: "[WARNING] ",
		LevelInfo:    "",
	}

	coloredHeaderFormat = "\n%s%s=== %s ===%s\n"
	headerFormat        = "\n=== %s ===\n"
)
