# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Version Mapping

本项目基于 [React Ink](https://github.com/vadimdemedes/ink) 进行复刻实现，版本对应关系：

| Go-React-Ink | React Ink | React Ink Commit | Sync Status |
|--------------|-----------|------------------|-------------|
| v0.1.0 | v6.8.0 | `be1b1bb6ec65056e2ed60ef3c5ae642704b82d31` | ✅ Initial implementation |

> 当 React Ink 发布新版本时，我们会持续同步更新。

## [0.1.0] - 2026-04-08

### Added

#### Compiler (Phase 1)
- **Scanner**: JSX region identification in Go source code
- **Lexer**: Tokenization of JSX content with 15 token types
- **Parser**: Recursive descent parser for JSX AST generation
- **Code Generator**: Go code generation from AST
- **CLI Tool**: `gox` command-line compiler

#### Core Runtime (Phase 2)
- **Fiber Reconciler**: React Fiber architecture implementation
  - Work loop with 16ms time slicing
  - Simplified 3-level priority system
  - beginWork/completeWork phases
  - Diff algorithm and commit phase
- **Terminal Renderer**: Output buffer with diff calculation
  - Cell-based rendering
  - ANSI escape sequence generation
  - Style support (color, bold, italic, underline)
- **Host Config**: React Reconciler host configuration

#### Components (Phase 3)
- **Box**: Flexbox container component
- **Text**: Text content with styling
- **Spacer**: Flexible whitespace
- **Newline**: Line break
- **Static**: Static content (excluded from re-render)
- **Transform**: Text transformation wrapper
- **Fragment**: Group children without wrapper

#### Hooks (Phase 3)
- **Core Hooks**: useState, useEffect, useLayoutEffect, useRef, useMemo, useCallback
- **Input Hooks**: useInput, useApp, useFocus, useFocusManager, useCursor, useAnimation, useStdout, useStdin, useStderr, useWindowSize, useBoxMetrics, usePaste, useIsScreenReaderEnabled

#### Layout Engine (Phase 4)
- **Flexbox**: Pure Go implementation
  - Direction (row, column, reverse)
  - JustifyContent (flex-start, center, flex-end, space-between, space-around)
  - AlignItems and AlignSelf
  - FlexGrow, FlexShrink, FlexBasis
  - Padding, Margin, Border
  - Performance: <2ns for simple layout, <1ms for 500 nodes

#### Tooling (Phase 5)
- **VSCode Extension**: Syntax highlighting and auto-compile for .gox files
  - LSP integration configuration
- **LSP Server**: Language server for .gox files
  - Diagnostics with compile error detection
  - Completion for components and hooks
  - Hover documentation for all APIs

#### Testing (Phase 6)
- **Comparison Tests**: React Ink behavior verification
  - Component tests (Box, Text, Spacer, Newline, Static, Transform)
  - Layout tests (Direction, Justify, Align, FlexGrow, Padding, Margin)
  - Performance benchmarks
  - MeasureFunc for text nodes

#### Tooling (Phase 5)
- **VSCode Extension**: Syntax highlighting and auto-compile for .gox files

### API Compatibility

| React Ink Feature | Status |
|-------------------|--------|
| Box | ✅ |
| Text | ✅ |
| Spacer | ✅ |
| Newline | ✅ |
| Static | ✅ |
| Transform | ✅ |
| useState | ✅ |
| useEffect | ✅ |
| useRef | ✅ |
| useMemo | ✅ |
| useInput | ✅ |
| useApp | ✅ |
| useFocus | ✅ |
| useCursor | ✅ |
| useAnimation | ✅ |
| Flexbox | ✅ |

### Test Coverage

- **208 unit tests** passing
- **13 integration tests** passing
- **5 e2e tests** passing
- **20 comparison tests** passing
- **8 performance benchmarks** passing
- **Total: 254 tests** passing

### Project Structure

```
go-ink/
├── cmd/gox/              # CLI compiler
├── internal/
│   ├── scanner/          # Source code scanner
│   ├── lexer/            # Lexical analyzer
│   ├── parser/           # Syntax parser
│   ├── codegen/          # Code generator
│   └── compiler/         # Compiler integration
├── pkg/
│   ├── core/             # Core types
│   ├── components/       # Built-in components
│   ├── hooks/            # React-like hooks
│   ├── input/            # Input handling hooks
│   ├── fiber/            # Fiber reconciler
│   ├── layout/           # Flexbox layout engine
│   ├── renderer/         # Terminal renderer
│   ├── hostconfig/       # Host configuration
│   └── tcell/            # Terminal abstraction
├── tools/
│   ├── vscode-gox/       # VSCode extension
│   └── lsp-server/       # LSP server
├── test/
│   ├── integration/      # Integration tests
│   ├── e2e/              # End-to-end tests
│   └── comparison/       # React Ink comparison tests
├── examples/             # Example applications
├── test/                 # Integration and e2e tests
└── docs/                 # Documentation
```

### Dependencies

- Go 1.21+
- golang.org/x/tools (code formatting)

### Known Limitations

- tcell integration uses standalone interface (Phase 2)
- Some advanced React Ink features may differ

### Breaking Changes

None (initial release)

### Security

No known security issues.

### Migration Guide

N/A (initial release)

---

## Future Plans

### [0.2.0] - Planned

- tcell terminal integration
- LSP server with completion and hover
- Performance optimization
- Additional examples and documentation

### [1.0.0] - Planned

- Full React Ink API compatibility
- Production-ready stability
- Comprehensive documentation
- Example gallery