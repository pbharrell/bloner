package main

import "github.com/pbharrell/bloner-server/connection"

type teamColor uint8

const (
	Black teamColor = iota
	Red
)

type Team struct {
	points    int
	teamColor teamColor
	players   [2]Player
}

func (t *Team) Arrange(clientId int, clientPos PlayPos) {
	t.players[0].Arrange(clientId, clientPos)
	t.players[1].Arrange(clientId, clientPos)
}

func (t *Team) Decode(teamColor teamColor, teamState connection.TeamState) {
	t.points = teamState.TricksWon
	t.players[0].Decode(teamColor, 0, teamState.PlayerState[0])
	t.players[1].Decode(teamColor, 1, teamState.PlayerState[1])
}

func (t *Team) Encode() connection.TeamState {
	return connection.TeamState{
		TricksWon: t.points,
		PlayerState: [2]connection.PlayerState{
			t.players[0].Encode(), t.players[1].Encode(),
		},
	}
}
