package pane

import (
	"io/fs"
	"time"
)

type FileEntry struct {
	Name    string
	Size    int64
	ModTime time.Time
	IsDir   bool
	Mode    fs.FileMode
	IsLink  bool
}
