package sndit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockPlayer struct {
	playing bool
	rewound bool
	paused  bool
	volume  float64
	playN   int
	pauseN  int
	rewindN int
}

func (m *mockPlayer) Play()               { m.playing = true; m.playN++ }
func (m *mockPlayer) Pause()              { m.playing = false; m.paused = true; m.pauseN++ }
func (m *mockPlayer) Rewind() error       { m.rewound = true; m.rewindN++; return nil }
func (m *mockPlayer) IsPlaying() bool     { return m.playing }
func (m *mockPlayer) SetVolume(v float64) { m.volume = v }

type mockContext struct {
	players []*mockPlayer
	err     error
}

func (m *mockContext) NewPlayer(data []byte) (Player, error) {
	if m.err != nil {
		return nil, m.err
	}
	p := &mockPlayer{}
	m.players = append(m.players, p)
	return p, nil
}

func (m *mockContext) NewLoopPlayer(data []byte) (Player, error) {
	return m.NewPlayer(data)
}

type testSound int

const (
	soundA testSound = iota
	soundB
)

func TestSfxEngine_Play(t *testing.T) {
	ctx := &mockContext{}
	engine := NewSfx[testSound](ctx)
	engine.Register(soundA, []byte("fake"))
	engine.Register(soundB, []byte("fake"))

	engine.Play(soundA)

	pa := ctx.players[0]
	require.True(t, pa.playing, "expected soundA to be playing")
	require.Equal(t, 1, pa.playN)
}

func TestSfxEngine_PlayRestartsWhilePlaying(t *testing.T) {
	ctx := &mockContext{}
	engine := NewSfx[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)
	engine.Play(soundA) // should restart

	pa := ctx.players[0]
	require.Equal(t, 1, pa.pauseN)
	require.Equal(t, 2, pa.rewindN)
	require.Equal(t, 2, pa.playN)
}

func TestSfxEngine_PlayUnregisteredIsNoop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewSfx[testSound](ctx)

	// Should not panic.
	engine.Play(soundA)
}

func TestMusicEngine_Play(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)

	pa := ctx.players[0]
	require.True(t, pa.playing, "expected soundA to be playing")
}

func TestMusicEngine_PlaySwitchesTracks(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))
	engine.Register(soundB, []byte("fake"))

	engine.Play(soundA)
	engine.Play(soundB)

	pa := ctx.players[0]
	pb := ctx.players[1]
	require.False(t, pa.playing, "expected soundA to be stopped")
	require.True(t, pb.playing, "expected soundB to be playing")
}

func TestMusicEngine_PlaySameTrackIsNoop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)
	engine.Play(soundA) // should not restart

	pa := ctx.players[0]
	require.Equal(t, 1, pa.playN)
}

func TestMusicEngine_Stop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)
	engine.Stop()

	pa := ctx.players[0]
	require.False(t, pa.playing, "expected soundA to be stopped")
	require.Equal(t, 2, pa.rewindN) // 1 from Play + 1 from Stop
}

func TestMusicEngine_StopWhenNotPlayingIsNoop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)

	// Should not panic.
	engine.Stop()
}

func TestSfxEngine_RegisterError(t *testing.T) {
	ctx := &mockContext{err: errors.New("bad data")}
	engine := NewSfx[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	// Player was not registered, so Play is a noop.
	engine.Play(soundA)
	require.Empty(t, ctx.players)
}

func TestMusicEngine_RegisterError(t *testing.T) {
	ctx := &mockContext{err: errors.New("bad data")}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	// Player was not registered, so Play is a noop.
	engine.Play(soundA)
	require.Empty(t, ctx.players)
}

func TestMusicEngine_PlayUnregisteredIsNoop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)

	// Should not panic.
	engine.Play(soundA)
}
