[![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/disgoplayer.svg)](https://pkg.go.dev/github.com/disgoorg/disgoplayer)
[![Go Report](https://goreportcard.com/badge/github.com/disgoorg/disgoplayer)](https://goreportcard.com/report/github.com/disgoorg/disgoplayer)
[![Go Version](https://img.shields.io/github/go-mod/go-version/disgoorg/disgoplayer)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgoplayer/blob/master/LICENSE)
[![DisGoPlayer Version](https://img.shields.io/github/v/tag/disgoorg/disgoplayer?label=release)](https://github.com/disgoorg/disgoplayer/releases/latest)
[![Support Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/zQ4u3CdU3J)

<img align="right" src="/.github/disgoplayer.png" width=192 alt="discord gopher">

# DisGoPlayer

The disgoplayer module provides opus/pcm/mp3 audio encoding/decoding/resampling as C bindings based on the [libopus](https://github.com/xiph/opus), [libsamplerate](http://mega-nerd.com/SRC) and [mpg123](https://mpg123.de) libraries.
It also lets you combine multiple pcm streams into a single pcm stream.
This module requires [CGO](https://go.dev/blog/cgo) to be enabled.

## Getting Started

### Installing

```sh
$ go get github.com/disgoorg/disgoplayer
```

### Usage

// TODO