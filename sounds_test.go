package sndit

import "testing"

type mockPlayer struct {
	playing  bool
	rewound  bool
	paused   bool
	volume   float64
	playN    int
	pauseN   int
	rewindN  int
}

func (m *mockPlayer) Play()              { m.playing = true; m.playN++ }
func (m *mockPlayer) Pause()             { m.playing = false; m.paused = true; m.pauseN++ }
func (m *mockPlayer) Rewind() error      { m.rewound = true; m.rewindN++; return nil }
func (m *mockPlayer) IsPlaying() bool    { return m.playing }
func (m *mockPlayer) SetVolume(v float64) { m.volume = v }

type mockContext struct {
	players []*mockPlayer
	idx     int
}

func (m *mockContext) NewPlayer(data []byte) (Player, error) {
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
	if !pa.playing {
		t.Error("expected soundA to be playing")
	}
	if pa.playN != 1 {
		t.Errorf("expected 1 play call, got %d", pa.playN)
	}
}

func TestSfxEngine_PlayRestartsWhilePlaying(t *testing.T) {
	ctx := &mockContext{}
	engine := NewSfx[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)
	engine.Play(soundA) // should restart

	pa := ctx.players[0]
	if pa.pauseN != 1 {
		t.Errorf("expected 1 pause call, got %d", pa.pauseN)
	}
	if pa.rewindN != 2 {
		t.Errorf("expected 2 rewind calls, got %d", pa.rewindN)
	}
	if pa.playN != 2 {
		t.Errorf("expected 2 play calls, got %d", pa.playN)
	}
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
	if !pa.playing {
		t.Error("expected soundA to be playing")
	}
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

	if pa.playing {
		t.Error("expected soundA to be stopped")
	}
	if !pb.playing {
		t.Error("expected soundB to be playing")
	}
}

func TestMusicEngine_PlaySameTrackIsNoop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)
	engine.Play(soundA) // should not restart

	pa := ctx.players[0]
	if pa.playN != 1 {
		t.Errorf("expected 1 play call, got %d", pa.playN)
	}
}

func TestMusicEngine_Stop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)
	engine.Register(soundA, []byte("fake"))

	engine.Play(soundA)
	engine.Stop()

	pa := ctx.players[0]
	if pa.playing {
		t.Error("expected soundA to be stopped")
	}
	if pa.rewindN != 2 { // 1 from Play + 1 from Stop
		t.Errorf("expected 2 rewind calls, got %d", pa.rewindN)
	}
}

func TestMusicEngine_StopWhenNotPlayingIsNoop(t *testing.T) {
	ctx := &mockContext{}
	engine := NewMusic[testSound](ctx)

	// Should not panic.
	engine.Stop()
}
