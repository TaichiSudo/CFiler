package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Copy(src, dstDir string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(dstDir, filepath.Base(src))

	if info.IsDir() {
		return copyDir(src, dstPath)
	}
	return copyFile(src, dstPath)
}

func Move(src, dstDir string) error {
	dstPath := filepath.Join(dstDir, filepath.Base(src))

	err := os.Rename(src, dstPath)
	if err != nil {
		// Cross-volume: copy + delete
		info, statErr := os.Stat(src)
		if statErr != nil {
			return statErr
		}
		if info.IsDir() {
			if copyErr := copyDir(src, dstPath); copyErr != nil {
				return copyErr
			}
		} else {
			if copyErr := copyFile(src, dstPath); copyErr != nil {
				return copyErr
			}
		}
		return os.RemoveAll(src)
	}
	return nil
}

func Delete(path string) error {
	return os.RemoveAll(path)
}

func Rename(oldPath, newName string) error {
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)
	return os.Rename(oldPath, newPath)
}

func Mkdir(parentDir, name string) error {
	path := filepath.Join(parentDir, name)
	return os.MkdirAll(path, 0755)
}

func copyFile(src, dst string) error {
	// Don't overwrite existing
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("destination already exists: %s", dst)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return nil
	}
	return os.Chmod(dst, info.Mode())
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
