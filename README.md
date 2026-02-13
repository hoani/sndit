# sndit

[![CI](https://github.com/hoani/sndit/actions/workflows/ci.yml/badge.svg)](https://github.com/hoani/sndit/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/hoani/sndit/branch/main/graph/badge.svg)](https://codecov.io/gh/hoani/sndit)

Code generation tool and engine for sound in Go games.

## Install

```
go install github.com/hoani/sndit/cmd/sndit@latest
```

## Setup

Place `.wav` files in directories under your project:

```
mygame/
  sfx_play/       # sound effects (one-shot, restartable)
    move.wav
    die.wav
  mus_theme/      # music (looping, one-at-a-time)
    title.wav
```

Run the generator:

```
sndit -dir mygame -module github.com/you/mygame
```

This creates a `sounds_gen.go` in each directory with typed constants and `Init`/`Play`/`Stop` functions.

Or use `go:generate`:

```go
//go:generate sndit -dir . -module github.com/you/mygame
```

## Usage

### With Ebitengine

Use the `sndebiten` adapter package:

```go
import (
    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hoani/sndit/sndebiten"
)

audioCtx := audio.NewContext(44100)
ctx := sndebiten.New(audioCtx)
```

Then call the generated packages:

```go
if err := sfx_play.Init(ctx); err != nil {
    log.Fatal(err)
}

sfx_play.Play(sfx_play.Move)
```

```go
if err := mus_theme.Init(ctx); err != nil {
    log.Fatal(err)
}

mus_theme.Play(mus_theme.Title)
mus_theme.Stop()
```

### Other audio libraries

Implement `sndit.Context` and `sndit.Player` for your audio library, then pass it to the generated `Init` functions.
