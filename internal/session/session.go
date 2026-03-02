package session

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cfiler/internal/config"
)

const sessionFile = "session.json"

type State struct {
	LeftDir     string `json:"left_dir"`
	RightDir    string `json:"right_dir"`
	ActivePane  int    `json:"active_pane"`
	LeftCursor  int    `json:"left_cursor"`
	RightCursor int    `json:"right_cursor"`
}

func Load() (*State, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, sessionFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func Save(state State) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(dir, sessionFile)
	return os.WriteFile(path, data, 0644)
}
