package main

import "github.com/pbharrell/bloner-server/connection"

type teamColor uint8

const (
	Black teamColor = iota
	Red
)

type Team struct {
	tricksWon int
	teamColor teamColor
	players   [2]Player
}

func (t *Team) Decode(teamState connection.TeamState) {
	t.tricksWon = teamState.TricksWon
	t.players[0].Decode(teamState.PlayerState[0])
	t.players[1].Decode(teamState.PlayerState[1])
}

func (t *Team) Encode() connection.TeamState {
	return connection.TeamState{
		TricksWon: t.tricksWon,
		PlayerState: [2]connection.PlayerState{
			t.players[0].Encode(), t.players[1].Encode(),
		},
	}
}
