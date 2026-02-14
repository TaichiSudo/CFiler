package bookmark

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cfiler/internal/config"
)

const bookmarkFile = "bookmarks.json"

type Entry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func Load() ([]Entry, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, bookmarkFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func Save(entries []Entry) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(dir, bookmarkFile)
	return os.WriteFile(path, data, 0644)
}

func Add(name, path string) error {
	entries, err := Load()
	if err != nil {
		return err
	}

	// Don't add duplicates
	for _, e := range entries {
		if e.Path == path {
			return nil
		}
	}

	entries = append(entries, Entry{Name: name, Path: path})
	return Save(entries)
}

func Remove(path string) error {
	entries, err := Load()
	if err != nil {
		return err
	}

	var filtered []Entry
	for _, e := range entries {
		if e.Path != path {
			filtered = append(filtered, e)
		}
	}
	return Save(filtered)
}
