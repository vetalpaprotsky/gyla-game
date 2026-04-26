# Gyla Game

A 2D single-player card game built with [Ebitengine](https://ebitengine.org/), powered by the [Gyla](https://github.com/vetalpaprotsky/gyla) game engine.

Play a traditional trick-taking card game against 3 AI opponents.

## Running

```sh
go run .
```

## How to Play

1. **Choose trump** — at the start of each round, click one of the four suit buttons to select the trump suit.
2. **Play cards** — click a card from your hand to play it. Playable cards are highlighted with a red border.
3. AI opponents will automatically play their turns after yours.

## Game Rules

See the full rules in the [Gyla engine README](https://github.com/vetalpaprotsky/gyla#game-rules).

## Project Structure

- `main.go` — Ebitengine game loop and input handling.
- `draw.go` — rendering logic (cards, table, UI elements).
- `client.go` — `GameClient` interface with `LocalClient` (single-player) and `OnlineClient` (multiplayer stub).

## Dependencies

- [Ebitengine](https://github.com/hajimehoshi/ebiten) — 2D game engine for Go.
- [Gyla](https://github.com/vetalpaprotsky/gyla) — card game logic engine.

## License

[MIT](LICENSE)