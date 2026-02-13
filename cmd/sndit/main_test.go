package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	dir := t.TempDir()

	sfxDir := filepath.Join(dir, "sfx_play")
	if err := os.Mkdir(sfxDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sfxDir, "click.wav"), []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".", "-dir", dir, "-module", "github.com/hoani/testgame")
	cmd.Dir = filepath.Join(findModuleRoot(t), "cmd", "sndit")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(filepath.Join(sfxDir, "sounds_gen.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)

	if !strings.Contains(got, "package sfx_play") {
		t.Error("expected package sfx_play in output")
	}
	if !strings.Contains(got, "Click Sound = iota") {
		t.Error("expected Click const in output")
	}
	if !strings.Contains(got, `//go:embed "click.wav"`) {
		t.Error("expected embed directive in output")
	}
}

func TestCLI_MissingModule(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-dir", t.TempDir())
	cmd.Dir = filepath.Join(findModuleRoot(t), "cmd", "sndit")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected CLI to fail without -module flag")
	}
	if !strings.Contains(string(out), "-module flag is required") {
		t.Errorf("expected module error message, got: %s", out)
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find module root")
		}
		dir = parent
	}
}
