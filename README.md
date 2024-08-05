# Space Invaders
An emulation of [Space Invaders (1978)](https://www.wikiwand.com/en/Space_Invaders).

![demo](./demo.webp)

This project contains an 8080 emulator and the necessary fake hardware bits to run the original Space Invaders arcade ROM.

## Installation
Get a pre-built binary from the [Releases](https://github.com/braheezy/space-invaders/releases) page.

Or, install with Go:

    go install github.com/braheezy/space-invaders

Or, clone the repository and use `make`:

    git clone https://github.com/braheezy/space-invaders.git
    cd space-invaders
    make run

## Usage
Run the binary to start the game:

    space-invaders

Here are the controls:
| Action | Key |
| --- | --- |
| Insert credit | `c` |
| 1 player start | `1` |
| 2 player start | `2` |
| Shoot | `Spacebar` |
| Move | Arrow keys or WASD

The `cpm` command runs a pre-bundled test ROM to verify the 8080 CPU emulator. That can be executed as follows:

    > space-invaders cpm
    MICROCOSM ASSOCIATES 8080/8085 CPU DIAGNOSTIC
    VERSION 1.0  (C) 1980

    CPU IS OPERATIONAL

## Development
You need Go and the dependencies that [Ebiten engine](https://ebitengine.org/en/documents/install.html) has.

Run `make` for various commands to run.

## Roadmap
- [ ] DIP settings
- [ ] Persistent high score
