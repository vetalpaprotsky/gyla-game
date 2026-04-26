package main

import "github.com/vetalpaprotsky/gyla/game"

// GameClient abstracts the communication between the UI and the game logic.
// SingleplayerClient calls game methods directly.
// MultiplayerClient will send/receive messages over WebSocket.
type GameClient interface {
	Start() ([]game.GameEvent, error)
	Apply(action game.Action) ([]game.GameEvent, error)
}

// LocalClient wraps game.Game directly. All opponents are AI.
type LocalClient struct {
	game *game.Game
}

func NewLocalClient() *LocalClient {
	g := game.NewGame(
		"You", "Bot 1", "Bot 2", "Bot 3",
		"Your Team", "Opponent Team",
		false, true, true, true,
	)
	return &LocalClient{game: &g}
}

func (c *LocalClient) Start() ([]game.GameEvent, error) {
	return c.game.Start()
}

func (c *LocalClient) Apply(action game.Action) ([]game.GameEvent, error) {
	return c.game.Apply(action)
}

// OnlineClient will communicate with the server over WebSocket using protobuf.
// For now, all methods are stubs.
type OnlineClient struct{}

func NewOnlineClient() *OnlineClient {
	return &OnlineClient{}
}

func (c *OnlineClient) Start() ([]game.GameEvent, error) {
	return nil, nil
}

func (c *OnlineClient) Apply(action game.Action) ([]game.GameEvent, error) {
	return nil, nil
}
