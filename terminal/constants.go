package terminal

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
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
