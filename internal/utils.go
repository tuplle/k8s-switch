package internal

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// GetFilesFromDir retrieves all regular files directly inside the specified directory
// (not subdirectories or symlinks) and returns a map of file names to their absolute paths.
func GetFilesFromDir(dirPath string) (map[string]string, error) {
	files := make(map[string]string)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		absPath, err := filepath.Abs(filepath.Join(dirPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		files[entry.Name()] = absPath
	}

	return files, nil
}

// FilterByExtension returns the subset of files whose name has the given extension (e.g. ".yaml").
func FilterByExtension(files map[string]string, ext string) map[string]string {
	filtered := make(map[string]string)
	for name, path := range files {
		if filepath.Ext(name) == ext {
			filtered[name] = path
		}
	}
	return filtered
}

// CopyFile atomically copies the contents of the source file to the destination file,
// creating or replacing it with permissions 0600 (destination files may hold credentials,
// e.g. a kubeconfig). It writes to a temporary file in the destination's directory and
// renames it into place, so an interruption or write failure never leaves a partially
// written destination file behind.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	tmpFile, err := os.CreateTemp(filepath.Dir(dst), "."+filepath.Base(dst)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // no-op once the rename below succeeds

	if _, err := io.Copy(tmpFile, sourceFile); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, dst)
}

// RunAnotherTUI executes another terminal UI as a subprocess, redirecting stdin, stdout, and stderr to the current terminal.
func RunAnotherTUI(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// This blocks the program until another TUI finishes
	return cmd.Run()
}
