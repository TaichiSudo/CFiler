# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CFiler is a keyboard-driven dual-pane file manager (TUI) built with Go and the Bubble Tea framework. It runs on Windows, macOS, and Linux.

## Build & Run Commands

```bash
go build -o cfiler.exe .   # Build (Windows)
go run main.go             # Run without building
go mod tidy                # Clean dependencies
```

No tests or linter configuration exist yet.

## Architecture

The application follows Bubble Tea's **Elm Architecture** (Model-Update-View) with message passing for all async operations.

### Mode-based Input Handling

The root `App` model (`internal/app/app.go`) operates in distinct modes that control input routing: `modeNormal`, `modeDialog`, `modeSearch`, `modeBookmark`, `modeHelp`. The current mode determines which component receives key events in `Update()`.

### Component Composition

`App` composes sub-models, each following the Bubble Tea pattern:
- Two `pane.Model` instances (left/right file lists)
- `preview.Model` for text file preview (viewport-based)
- `statusbar.Model` for status display
- `dialog.Dialog` interface (confirm/input dialogs)
- `bookmark.Model` for bookmark management

### Message Flow

Async operations (directory loading, file operations, preview loading) return `tea.Cmd` functions that produce typed messages (`DirLoadedMsg`, `FileOpResultMsg`, `PreviewLoadedMsg`, `DialogResultMsg`). Errors surface through the status bar.

### Internal Clipboard

Copy/move uses an internal clipboard (`clipboard []string` + `clipAction`) rather than OS clipboard. The OS clipboard (`atotto/clipboard`) is used only for copying the current path.

### Package Layout

All code lives under `internal/` (not importable externally). Each feature has its own package: `app`, `pane`, `preview`, `statusbar`, `dialog`, `bookmark`, `fileops`, `config`.

## Key Conventions

- Bubble Tea component structs are named `Model`
- Message types use `Msg` suffix
- Key bindings are centralized in `internal/app/keys.go`
- Styling uses Tokyo Night color scheme, defined in `internal/app/styles.go`
- Platform-specific behavior (file open) is in `fileops/open.go` using `runtime.GOOS`
- Bookmark data persists as JSON in the OS-appropriate config directory (resolved by `config.ConfigDir()`)
