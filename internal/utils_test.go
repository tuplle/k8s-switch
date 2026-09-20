package internal

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetFilesFromDir(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.yaml", "b.conf"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("data"), 0o644); err != nil {
			t.Fatalf("failed to write fixture file %s: %v", name, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatalf("failed to create fixture subdir: %v", err)
	}

	files, err := GetFilesFromDir(dir)
	if err != nil {
		t.Fatalf("GetFilesFromDir returned error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(files), files)
	}

	for _, name := range []string{"a.yaml", "b.conf"} {
		path, ok := files[name]
		if !ok {
			t.Fatalf("expected %q in result, got %v", name, files)
		}
		if !filepath.IsAbs(path) {
			t.Errorf("expected absolute path for %q, got %q", name, path)
		}
	}

	if _, ok := files["subdir"]; ok {
		t.Errorf("expected subdirectories to be excluded, got %v", files)
	}
}

func TestGetFilesFromDirMissingDir(t *testing.T) {
	_, err := GetFilesFromDir(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

func TestFilterByExtension(t *testing.T) {
	input := map[string]string{
		"prod.yaml":    "/configs/prod.yaml",
		"staging.yaml": "/configs/staging.yaml",
		"dev.conf":     "/configs/dev.conf",
		"notes.txt":    "/configs/notes.txt",
	}

	got := FilterByExtension(input, ".yaml")

	want := map[string]string{
		"prod.yaml":    "/configs/prod.yaml",
		"staging.yaml": "/configs/staging.yaml",
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d entries, got %d: %v", len(want), len(got), got)
	}
	for name, path := range want {
		if got[name] != path {
			t.Errorf("expected %q -> %q, got %q", name, path, got[name])
		}
	}
}

func TestFilterByExtensionNoMatches(t *testing.T) {
	input := map[string]string{"dev.conf": "/configs/dev.conf"}

	got := FilterByExtension(input, ".yaml")

	if len(got) != 0 {
		t.Fatalf("expected no matches, got %v", got)
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.yaml")
	dst := filepath.Join(dir, "dst.yaml")
	want := []byte("apiVersion: v1\nkind: Config\n")

	if err := os.WriteFile(src, want, 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile returned error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("expected destination content %q, got %q", want, got)
	}
}

func TestCopyFileMissingSource(t *testing.T) {
	dir := t.TempDir()
	err := CopyFile(filepath.Join(dir, "does-not-exist.yaml"), filepath.Join(dir, "dst.yaml"))
	if err == nil {
		t.Fatal("expected error for missing source file, got nil")
	}
}

func TestCopyFileOverwritesDestination(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.yaml")
	dst := filepath.Join(dir, "dst.yaml")

	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}
	if err := os.WriteFile(dst, []byte("old-longer-content"), 0o644); err != nil {
		t.Fatalf("failed to write destination file: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile returned error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}
	if string(got) != "new" {
		t.Errorf("expected destination to be fully overwritten with %q, got %q", "new", got)
	}
}

func TestRunAnotherTUI(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test relies on a POSIX shell command")
	}

	if err := RunAnotherTUI("true", nil); err != nil {
		t.Fatalf("expected no error running a successful command, got %v", err)
	}
}

func TestRunAnotherTUIMissingCommand(t *testing.T) {
	if err := RunAnotherTUI("k8s-switch-definitely-not-a-real-command", nil); err == nil {
		t.Fatal("expected error for a nonexistent command, got nil")
	}
}
