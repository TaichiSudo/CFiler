package app

import "cfiler/internal/pane"

type DirLoadedMsg struct {
	Entries []pane.FileEntry
	Path    string
	PaneID  int
}

type DirLoadErrorMsg struct {
	Err    error
	PaneID int
}

type PreviewLoadedMsg struct {
	Content  string
	IsBinary bool
}

type FileOpResultMsg struct {
	Err error
	Op  string
}

type DialogResultMsg struct {
	Confirmed bool
	Text      string
	Action    string
}

type StatusMsg struct {
	Text string
}
