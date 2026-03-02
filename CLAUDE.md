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
- `session` package (no Model) for session state persistence

### Message Flow

Async operations (directory loading, file operations, preview loading) return `tea.Cmd` functions that produce typed messages (`DirLoadedMsg`, `FileOpResultMsg`, `PreviewLoadedMsg`, `DialogResultMsg`). Errors surface through the status bar.

### Internal Clipboard

Copy/move uses an internal clipboard (`clipboard []string` + `clipAction`) rather than OS clipboard. The OS clipboard (`atotto/clipboard`) is used only for copying the current path.

### Package Layout

All code lives under `internal/` (not importable externally). Each feature has its own package: `app`, `pane`, `preview`, `statusbar`, `dialog`, `bookmark`, `fileops`, `config`, `session`.

## Key Conventions

- Bubble Tea component structs are named `Model`
- Message types use `Msg` suffix
- Key bindings are centralized in `internal/app/keys.go`
- Styling uses Tokyo Night color scheme, defined in `internal/app/styles.go`
- Platform-specific behavior (file open) is in `fileops/open.go` using `runtime.GOOS`
- Bookmark data persists as JSON in the OS-appropriate config directory (resolved by `config.ConfigDir()`)
- Session state (pane dirs, active pane, cursor positions) persists as `session.json` in the same config directory, loaded on startup and saved on quit/tab-switch/dir-load

## Implementation History

### Session State Persistence (2026-03-03)

起動時に前回のディレクトリ・アクティブペイン・カーソル位置を復元する機能を追加。

**新規ファイル**
- `internal/session/session.go` — `State` 構造体と `Load()`/`Save()` 関数。`bookmark` パッケージと同じ永続化パターンを踏襲し `{ConfigDir}/session.json` に保存。ファイル未存在時は `nil, nil` を返す。

**変更ファイル**
- `internal/pane/pane.go` — `SetCursor(n int)` メソッドを追加。`SetEntries()` はカーソルを 0 にリセットするため、初期ロード後にカーソルを復元するために必要。
- `internal/app/app.go`
  - `App` 構造体に `initCursor [2]int` フィールドを追加（-1 = 復元不要）
  - `New()` でセッションをロードし、前回の左右ディレクトリ・アクティブペイン・カーソル位置を初期値として使用
  - `resolveStartDirs()` ヘルパー追加：保存パスが存在しない場合は cwd にフォールバック。空文字列は Windows のドライブ一覧として有効扱い
  - `saveSession()` ヘルパーメソッド追加
  - `DirLoadedMsg` ハンドラ：エントリセット後に `initCursor` を適用し `saveSession()` を呼ぶ
  - Quit キー：`tea.Quit` 前に `saveSession()`
  - Tab キー：ペイン切替後に `saveSession()`
