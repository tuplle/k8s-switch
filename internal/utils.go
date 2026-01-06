package internal

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// GetFilesFromDir retrieves all files from the specified directory and returns a map of file names to their absolute paths.
func GetFilesFromDir(dirPath string) (map[string]string, error) {
	files := make(map[string]string)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
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

// CopyFile copies the contents of the source file to the destination file.
// It returns an error if the source file cannot be opened, the destination file cannot be created, or the copy operation fails.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
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
