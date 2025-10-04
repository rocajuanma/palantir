# Palantir

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocajuanma/palantir)](https://goreportcard.com/report/github.com/rocajuanma/palantir)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23.6-blue.svg)](https://golang.org/)
[![GitHub release](https://img.shields.io/github/release/rocajuanma/palantir.svg)](https://github.com/rocajuanma/palantir/releases)

A lightweight Go package for enhanced terminal output, featuring colored text, emoji indicators, and consistent formatting.

> In *The Lord of the Rings* lore, a Palantír is a seeing-stone that enables users to communicate and observe distant events. Similarly, this package empowers you to gain clearer insights into your program's output, making it easier to monitor and understand what's happening.


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

## Features

- Multiple output levels
- Colored output: fully-coloured line, level-only or no-colour
- Emoji support
- Progress indicators
- Interactive confirmations


<p align="center">
  <img src="./cmd/demo/terminal.png" alt="Palantir Demo">
</p>


See the [Palantir terminal demo](cmd/demo/README.md) for usage examples.
