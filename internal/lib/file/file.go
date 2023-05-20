package file

import (
	"fmt"
	"os"
	"time"
)

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(fp string) bool {
	f, e := os.Stat(fp)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func EnsureDirRW(dataDir string) error {
	err := EnsureDir(dataDir)
	if err != nil {
		return err
	}

	checkFile := fmt.Sprintf("%s/rw.%d", dataDir, time.Now().UnixNano())
	fd, err := Create(checkFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("open %s: rw permission denied", dataDir)
		}
		return err
	}

	if err := Close(fd); err != nil {
		return fmt.Errorf("close error: %s", err)
	}

	if err := Remove(checkFile); err != nil {
		return fmt.Errorf("remove error: %s", err)
	}

	return nil
}

// Create one file
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Remove one file
func Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll files
func RemoveAll(directory string) error {
	return os.RemoveAll(directory)
}

// Close fd
func Close(fd *os.File) error {
	return fd.Close()
}

// EnsureDir mkdir dir if not exist
func EnsureDir(fp string) error {
	return os.MkdirAll(fp, os.ModePerm)
}

// CreateTempDir creates a temporary folder with defined pattern
func CreateTempDir(pattern string) (string, error) {
	tempDir, err := os.MkdirTemp("", pattern)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}
