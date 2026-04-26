package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vetalpaprotsky/gyla/game"
)

// App implements ebiten.Game interface.
// Ebitengine requires you to implement 3 methods:
//   - Update()  — called at a fixed 60 TPS (game logic goes here)
//   - Draw()    — called at the monitor's refresh rate (rendering goes here)
//   - Layout()  — returns the logical screen size
type App struct {
	client   GameClient
	view     *game.GameView
	started  bool
	clicking bool // prevents repeated clicks while mouse is held
}

// Update is called 60 times per second (TPS = ticks per second).
// This is where you handle input, update game state, etc.
func (a *App) Update() error {
	if !a.started {
		events, err := a.client.Start()
		if err != nil {
			return err
		}
		a.started = true
		if len(events) > 0 {
			a.updateView(events)
		}
	}

	if a.view == nil {
		return nil
	}

	// Handle mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if a.clicking {
			return nil // wait for release
		}
		a.clicking = true

		mx, my := ebiten.CursorPosition()
		next := a.view.NextAction

		if next.Name == "assign_trump" && next.Player == game.Player1 {
			if suit, ok := trumpButtonHit(mx, my); ok {
				a.applyAction(game.Action{
					Name:   game.AssignTrumpAction,
					Player: game.Player1,
					Suit:   suit,
				})
			}
		} else if next.Name == "play_card" && next.Player == game.Player1 {
			if hc, ok := handCardHit(mx, my, a.view); ok {
				a.applyAction(game.Action{
					Name:   game.PlayCardAction,
					Player: game.Player1,
					Rank:   hc.Card.Rank,
					Suit:   hc.Card.Suit,
				})
			}
		}
	} else {
		a.clicking = false
	}

	return nil
}

func (a *App) applyAction(action game.Action) {
	events, err := a.client.Apply(action)
	if err != nil {
		return
	}
	if len(events) > 0 {
		a.updateView(events)
	}
}

func (a *App) updateView(events []game.GameEvent) {
	lastState := events[len(events)-1].GameState
	view := lastState.ViewFor(game.Player1)
	a.view = &view
}

// Draw is called at the monitor's refresh rate (e.g. 60Hz, 144Hz).
// It runs independently from Update. The screen is cleared automatically before each call.
func (a *App) Draw(screen *ebiten.Image) {
	drawGame(screen, a.view)
}

// Layout returns the logical screen size.
// Ebitengine scales this to fit the actual window size.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Gyla")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := &App{
		client: NewLocalClient(),
	}

	// RunGame starts the game loop. It blocks until the window is closed.
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
