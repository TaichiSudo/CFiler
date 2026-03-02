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

The root `App` model (`internal/app/app.go`) operates in distinct modes that control input routing:

| Mode | Constant | Description |
|------|----------|-------------|
| Normal | `modeNormal` | File navigation and commands |
| Dialog | `modeDialog` | Confirm/input dialog open |
| Search | `modeSearch` | Search bar active |
| Bookmark | `modeBookmark` | Bookmark selector open |
| Help | `modeHelp` | Help overlay |

`Update()` dispatches to `handleKey()`, which routes to mode-specific handlers (e.g., `handleNormalKey()`, `handleSearchKey()`).

### Component Composition

`App` composes sub-models, each following the Bubble Tea pattern:
- Two `pane.Model` instances (left/right file lists)
- `preview.Model` for text file preview (viewport-based, max 64KB, with binary/UTF-8 detection)
- `statusbar.Model` for status display
- `dialog.Dialog` interface (confirm/input dialogs)
- `bookmark.Model` for bookmark management
- `session` package (no Model) for session state persistence

### Message Flow

Async operations return `tea.Cmd` functions that produce typed messages handled in `Update()`:

| Message | Trigger |
|---------|---------|
| `DirLoadedMsg` | Directory listing complete |
| `FileOpResultMsg` | Copy/move/delete done |
| `PreviewLoadedMsg` | File preview loaded |
| `DialogResultMsg` | User confirmed/cancelled dialog |

Errors surface through the status bar rather than panicking.

### Dialog Action Strings

Dialog results encode the operation as an action string parsed with `strings.SplitN(action, ":", 2)`:
- `"delete:path"` — delete confirmation
- `"mkdir:dir"` — new directory
- `"rename:oldname"` — rename operation
- `"goto:0"` / `"goto:1"` — jump to pane

### Internal Clipboard

Copy/move uses an internal clipboard (`clipboard []string` + `clipAction`) rather than OS clipboard. The OS clipboard (`atotto/clipboard`) is used only for copying the current path.

### Pane Model Details

`pane.Model` (`internal/pane/pane.go`) key points:
- `marked map[string]bool` uses **filename** (not full path) as key
- `SetEntries()` always resets cursor to 0; use `SetCursor(n)` afterward to restore position
- Empty string `""` for `dir` is treated as the Windows drive list (not an error)
- `Entries()` returns filtered results when searching, otherwise all entries

### Package Layout

All code lives under `internal/` (not importable externally):

| Package | Responsibility |
|---------|---------------|
| `app` | Root model, Update/View, key dispatch |
| `pane` | File list, navigation, multi-select, search |
| `preview` | Text file preview viewport |
| `statusbar` | Bottom status line |
| `dialog` | Confirm and text-input dialog interface |
| `bookmark` | Bookmark list UI and persistence |
| `fileops` | Copy/move/delete/rename/mkdir; platform file open |
| `config` | `ConfigDir()` → `{os.UserConfigDir()}/cfiler/` |
| `session` | `session.json` save/load for startup state restoration |

## Key Conventions

- Bubble Tea component structs are named `Model`
- Message types use `Msg` suffix
- Key bindings are centralized in `internal/app/keys.go` — never add bindings elsewhere
- Styling uses Tokyo Night color scheme, defined in `internal/app/styles.go`
- Platform-specific behavior (file open) is in `fileops/open.go` using `runtime.GOOS` (`cmd /c start`, `open`, `xdg-open`)
- Both bookmark data and session state persist as JSON in `config.ConfigDir()`
- Session state (pane dirs, active pane, cursor positions) is saved on quit, tab-switch, and directory load; restored on startup via `initCursor [2]int` field (`-1` = no restore needed)
- File operations detect cross-volume moves: when `os.Rename` fails, fall back to copy + delete
- Search is case-insensitive substring matching — no regex library
