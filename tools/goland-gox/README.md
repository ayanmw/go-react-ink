# GoX Plugin for GoLand 2024.3.6

IntelliJ IDEA plugin providing JSX syntax support for Go (.gox files).

## Features

- **Syntax Highlighting** - Go keywords, strings, comments, and JSX tags
- **Code Completion** - Components (Box, Text, etc.) and Hooks (useState, useEffect, etc.)
- **Error Diagnostics** - Basic validation for unclosed tags
- **Compile Actions** - Compile .gox files directly from IDE
- **Color Settings** - Customizable syntax colors in Settings → Editor → Color Scheme → GoX

## Installation

### From JetBrains Marketplace (Recommended)

1. Open GoLand → Settings → Plugins
2. Search for "GoX"
3. Click Install

### From Source

```bash
cd tools/goland-gox
./gradlew buildPlugin
```

The plugin will be at `build/distributions/goland-gox-0.1.0.zip`.

Install manually: Settings → Plugins → ⚙️ → Install Plugin from Disk

## Requirements

- GoLand 2024.3.6 or IntelliJ IDEA Ultimate with Go plugin
- JDK 17+

## Usage

### Compile Current File

- **Shortcut**: `Ctrl+Shift+G`
- **Menu**: Tools → Compile .gox File

### Compile All Files

- **Menu**: Tools → Compile All .gox Files

### Code Completion

Type inside JSX context for component completion:

```go
func App() core.Element {
    return <Bo  // Ctrl+Space → Box
}
```

## Configuration

Settings → Editor → Color Scheme → GoX

Customize colors for:
- Keywords, Strings, Numbers
- JSX Tags, Attributes, Expressions
- Comments, Operators

## Development

### Build

```bash
./gradlew buildPlugin
```

### Run in Development Mode

```bash
./gradlew runIde
```

### Publish

```bash
./gradlew publishPlugin
```

Requires `PUBLISH_TOKEN` environment variable.

## Project Structure

```
src/main/
├── kotlin/com/ayanmw/gox/
│   ├── GoxLanguage.kt          # Language definition
│   ├── GoxFileType.kt          # File type
│   ├── GoxLexer.kt             # Lexer (Go + JSX)
│   ├── GoxParser.kt            # Parser
│   ├── GoxSyntaxHighlighter.kt # Highlighting
│   ├── GoxCompletionContributor.kt # Code completion
│   ├── GoxAnnotator.kt         # Error annotations
│   └── actions/
│       └── CompileActions.kt   # Compile actions
├── resources/
│   ├── META-INF/plugin.xml     # Plugin manifest
│   └── icons/gox.svg           # File icon
```

## License

MIT License