package ebitenaudio

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hoani/sndit"
)

// Context wraps an ebiten *audio.Context and implements sndit.Context.
type Context struct {
	ctx *audio.Context
}

// New creates a new Context from an ebiten *audio.Context.
func New(ctx *audio.Context) *Context {
	return &Context{ctx: ctx}
}

// NewPlayer decodes WAV data and creates an ebiten audio player.
func (c *Context) NewPlayer(data []byte) (sndit.Player, error) {
	stream, err := wav.DecodeF32(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("ebitenaudio: decode wav: %w", err)
	}
	p, err := c.ctx.NewPlayerF32(stream)
	if err != nil {
		return nil, fmt.Errorf("ebitenaudio: new player: %w", err)
	}
	return p, nil
}

// NewLoopPlayer decodes WAV data and creates a looping ebiten audio player.
func (c *Context) NewLoopPlayer(data []byte) (sndit.Player, error) {
	stream, err := wav.DecodeF32(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("ebitenaudio: decode wav: %w", err)
	}
	loop := audio.NewInfiniteLoopF32(stream, stream.Length())
	p, err := c.ctx.NewPlayerF32(loop)
	if err != nil {
		return nil, fmt.Errorf("ebitenaudio: new player: %w", err)
	}
	return p, nil
}
