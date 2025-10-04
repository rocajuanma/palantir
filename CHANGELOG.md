# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

### Changed

### Fixed

## [1.0.0] - 2025-10-03

### Added
- Usage section in README with basic usage examples, custom configuration, and global handler examples
- Palantir logo in README header for better visual branding

### Changed
- **Major restructuring**: Consolidated all code into a single `palantir` package
  - Moved `constants.go`, `output.go`, and `output_test.go` to root package
  - Moved demo application from `cmd/terminal/` to `cmd/demo/`
  - Updated all imports to use the unified package structure
  - Simplified package structure for better usability

### Fixed
- Fixed image path in README from `cmd/terminal/terminal.png` to `cmd/demo/terminal.png`
- Fixed demo link in README from `cmd/terminal/README.md` to `cmd/demo/README.md`

## [0.5.0] - 2025-10-03

### Added
- New level-only coloured option (Thank you, @muyiwaolurin!)
- Workflows to auto-generated Release and CHANGELOG updates using rocajuanma/simple-release
  
### Changed

### Fixed
- Update correct workflows @main tag usage

## [0.2.0] - 2025-09-30

### Added
- Multiple output levels (Info, Warning, Error, Success, Stage, Header)
- Colored terminal output with ANSI escape codes
- Emoji support for visual indicators
- Configurable output formatting options
- Progress indicators and interactive confirmations
- Level-only color mode for subtle highlighting
- Comprehensive test coverage

### Features
- **Output Levels**: Six distinct output levels with appropriate colors and emojis
- **Color Support**: Full ANSI color support with customizable color schemes
- **Emoji Integration**: Visual emojis for different output types (✅ success, ❌ error, ⚠️ warning, 🔧 stage)
- **Flexible Configuration**: Toggle colors, emojis, formatting, and verbose mode
- **Terminal Demo**: Complete example application showcasing all features

### Technical Details
- Go 1.23.6+ compatibility
- Interface-based design for easy testing and mocking
- Configurable output handler with sensible defaults
- Support for both colored and plain text output
- Cross-platform terminal compatibility
