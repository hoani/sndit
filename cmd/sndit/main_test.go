package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()

	sfxDir := filepath.Join(dir, "sfx_play")
	require.NoError(t, os.Mkdir(sfxDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(sfxDir, "click.wav"), []byte("fake"), 0o644))

	err := run([]string{"-dir", dir, "-module", "github.com/hoani/testgame"})
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(sfxDir, "sfx_play.go"))
	require.NoError(t, err)
	got := string(data)

	require.Contains(t, got, "package sfx_play")
	require.Contains(t, got, "Click Sound = iota")
	require.Contains(t, got, `//go:embed "click.wav"`)
}

func TestRun_MissingModule(t *testing.T) {
	err := run([]string{"-dir", t.TempDir()})
	require.EqualError(t, err, "-module flag is required")
}

func TestRun_DirMissingValue(t *testing.T) {
	err := run([]string{"-dir"})
	require.EqualError(t, err, "-dir requires a value")
}

func TestRun_ModuleMissingValue(t *testing.T) {
	err := run([]string{"-module"})
	require.EqualError(t, err, "-module requires a value")
}

func TestRun_NoArgs(t *testing.T) {
	err := run([]string{})
	require.EqualError(t, err, "-module flag is required")
}

func TestRun_InvalidDir(t *testing.T) {
	err := run([]string{"-dir", "/nonexistent/path", "-module", "example.com/mod"})
	require.Error(t, err)
}
