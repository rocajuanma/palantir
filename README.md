<div align="center">
  <img src="assets/palantir.png" alt="Palantir Logo" width="200" style="border-radius: 50%;">
  <h1>Palantir</h1>
</div>
<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocajuanma/palantir)](https://goreportcard.com/report/github.com/rocajuanma/palantir)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23.6-blue.svg)](https://golang.org/)
[![GitHub release](https://img.shields.io/github/release/rocajuanma/palantir.svg)](https://github.com/rocajuanma/palantir/releases)

A lightweight Go package for enhanced terminal output, featuring colored text, emoji indicators, and consistent formatting.
</div>

> In *The Lord of the Rings* lore, a Palantír is a seeing-stone that enables users to communicate and observe distant events. Similarly, this package empowers you to gain clearer insights into your program's output, making it easier to monitor and understand what's happening.

## Features
- **Multiple Output Levels:** Easily distinguish between informational messages, warnings, errors, successes, headers, and stages with dedicated methods for each.
- **Flexible Colored Output:** Choose between fully-colored lines, colorizing only the output level, or disabling colors entirely to suit your terminal or preferences.
- **Tree Visualization:** Visualize files, directories, and YAML content in a structured tree format for easier navigation and understanding.
- **Progress Indicators:** Display animated progress spinners or status updates to track ongoing operations in real time.
- **Interactive Confirmations:** Prompt users for confirmation or input directly in the terminal, supporting interactive workflows.
- **Customizable Configuration:** Fine-tune output behavior with options for verbosity, formatting, color usage, and emoji toggling.

## Installation

```bash
go get github.com/rocajuanma/palantir
```

## Usage

### Basic Usage

```go
package main

import (
    "github.com/rocajuanma/palantir"
)

func main() {
    // Create a default output handler
    handler := palantir.NewDefaultOutputHandler()
    
    // Use different output levels
    handler.PrintHeader("My Application")
    handler.PrintInfo("Starting process...")
    handler.PrintSuccess("Operation completed!")
    handler.PrintWarning("This is a warning")
    handler.PrintError("Something went wrong")
    handler.PrintStage("Processing stage 1")

    // Display directory tree structure
    err, _ := palantir.ShowHierarchy("/path/to/directory", "")
    if err != nil {
        panic(err)
    }

    // Display YAML content in tree structure
    yamlContent, err := os.ReadFile("config.yaml")
    if err != nil {
        panic(err)
    }
    
    // Display YAML as tree structure
    err = palantir.ShowYAMLHierarchy(yamlContent)
    if err != nil {
        panic(err)
    }
    
    // Or read from file directly
    err = palantir.ShowYAMLHierarchyFromFile("config.yaml")
    if err != nil {
        panic(err)
    }
}
```

### Custom Configuration

```go
// Create a custom configuration
config := &palantir.OutputConfig{
    UseColors:         true,
    UseEmojis:         false,
    UseFormatting:     true,
    DisableOutput:     false,
    VerboseMode:       false,
    ColorizeLevelOnly: true,
}

handler := palantir.NewOutputHandler(config)
```

Check out the [Palantir demo](cmd/demo/README.md) for detailed usage examples, advanced capabilities, and interactive feature showcases.

<p align="center">
  <img src="./assets/terminal.png" alt="Palantir Demo">
</p>



