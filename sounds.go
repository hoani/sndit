package sndit

import "fmt"

// Player controls playback of a single audio source.
type Player interface {
	Play()
	Pause()
	Rewind() error
	IsPlaying() bool
	SetVolume(float64)
}

// Context creates new audio players from raw WAV data.
type Context interface {
	NewPlayer(data []byte) (Player, error)
	NewLoopPlayer(data []byte) (Player, error)
}

// SfxEngine manages sound effects (one-shot, restartable).
type SfxEngine[S ~int] struct {
	ctx     Context
	players map[S]Player
}

// NewSfx creates a new SfxEngine.
func NewSfx[S ~int](ctx Context) *SfxEngine[S] {
	return &SfxEngine[S]{
		ctx:     ctx,
		players: make(map[S]Player),
	}
}

// Register creates a player for the given sound ID from raw WAV data.
func (e *SfxEngine[S]) Register(id S, data []byte) error {
	if _, ok := e.players[id]; ok {
		return fmt.Errorf("sound %d already registered", id)
	}
	p, err := e.ctx.NewPlayer(data)
	if err != nil {
		return err
	}
	e.players[id] = p
	return nil
}

// Play plays the sound effect with the given ID, restarting if already playing.
func (e *SfxEngine[S]) Play(id S) {
	p, ok := e.players[id]
	if !ok {
		return
	}
	if p.IsPlaying() {
		p.Pause()
	}
	p.Rewind()
	p.Play()
}

// MusicEngine manages music tracks (looping, one-at-a-time).
type MusicEngine[S ~int] struct {
	ctx     Context
	players map[S]Player
	current S
	playing bool
}

// NewMusic creates a new MusicEngine.
func NewMusic[S ~int](ctx Context) *MusicEngine[S] {
	return &MusicEngine[S]{
		ctx:     ctx,
		players: make(map[S]Player),
		current: -1,
	}
}

// Register creates a looping player for the given sound ID from raw WAV data.
func (e *MusicEngine[S]) Register(id S, data []byte) error {
	if _, ok := e.players[id]; ok {
		return fmt.Errorf("sound %d already registered", id)
	}
	p, err := e.ctx.NewLoopPlayer(data)
	if err != nil {
		return err
	}
	e.players[id] = p
	return nil
}

// Play starts the music track with the given ID, stopping any current track.
func (e *MusicEngine[S]) Play(id S) {
	if e.playing && e.current == id {
		return
	}
	e.stopCurrent()
	p, ok := e.players[id]
	if !ok {
		return
	}
	p.Rewind()
	p.Play()
	e.current = id
	e.playing = true
}

// Stop stops the currently playing music track.
func (e *MusicEngine[S]) Stop() {
	e.stopCurrent()
}

func (e *MusicEngine[S]) stopCurrent() {
	if !e.playing {
		return
	}
	if p, ok := e.players[e.current]; ok {
		p.Pause()
		p.Rewind()
	}
	e.playing = false
}
