# Space Invaders
An emulation of [Space Invaders (1978)](https://www.wikiwand.com/en/Space_Invaders).

![demo](./demo.webp)

| Action | Key |
| --- | --- |
| Insert credit | `c` |
| 1 player start | `1` |
| 2 player start | `2` |
| Shoot | `Spacebar` |
| Move | Arrow keys or WASD

## Installation
Get a pre-built binary from the on the [Releases](https://github.com/braheezy/space-invaders/releases) page.

Or, install with Go:

    go install github.com/braheezy/space-invaders

Or, clone the repository and use `make`:

    git clone https://github.com/braheezy/space-invaders.git
    cd space-invaders
    make run

## Development
You need Go and the dependencies that [Ebiten engine](https://ebitengine.org/en/documents/install.html) has.

Run `make` for various commands to run.

## Roadmap
- [ ] DIP settings
- [ ] Persistent high score
