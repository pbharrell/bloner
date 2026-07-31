package main

import (
	"encoding/json"
	"fmt"

	"github.com/pbharrell/bloner-server/connection"
)

func (g *Game) HandleLobbyAssignMessage(data connection.LobbyAssign) {
	println("Player with id:", data.PlayerId)
	println("Lobby with id:", data.LobbyId)

	g.id = data.PlayerId
	g.lobbyId = data.LobbyId
	g.mode = LobbyAssigned
	g.PushPage(CreateLobbyAssignedPage(g))
	g.debugPrintln("Handled lobby assign message!")
}

func (g *Game) HandleGameStartMessage(data connection.GameStart) {
	gameActivePage := CreateGameActivePage(g)
	for i := range gameActivePage.state.teams {
		for j := range gameActivePage.state.teams[i].players {
			gameActivePage.state.teams[i].players[j].Id = data[i*2+j]
		}
	}

	g.PopPage()
	g.PushPage(gameActivePage)
}

func (g *Game) HandleStateRequestMessage() {
	g.SendStateResponse()
}

func (g *Game) HandleStateResponseMessage(data connection.StateResponse) {
	gameActivePage, ok := g.pageStack[len(g.pageStack)-1].(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	g.mode = GameActive
	// g.PushPage(createGameActivePage(g))
	g.DecodeGameState(data)
	for i, t := range gameActivePage.state.teams {
		for _, p := range t.players {
			println("Player on team", i, "with id", p.Id)
		}
	}
}

func (g *Game) HandleTurnInfoMessage(data connection.TurnInfo) {
	g.DecodeTurnInfo(data)
}

func (g *Game) HandleMessage(msg connection.Message) {
	raw, err := json.Marshal(msg.Data)
	if err != nil {
		println("marshal error:", err)
		return
	}

	switch msg.Type {
	case "lobby_assign":
		var lobbyAssign connection.LobbyAssign
		if err := json.Unmarshal(raw, &lobbyAssign); err != nil {
			println("LobbyAssign unmarshal error:", err)
			return
		}

		g.HandleLobbyAssignMessage(lobbyAssign)

	case "game_start":
		var gameStart connection.GameStart
		if err := json.Unmarshal(raw, &gameStart); err != nil {
			println("LobbyAssign unmarshal error:", err)
			return
		}

		g.HandleGameStartMessage(gameStart)

	case "state_req":
		g.HandleStateRequestMessage()

	case "state_res":
		var stateResponse connection.StateResponse
		if err := json.Unmarshal(raw, &stateResponse); err != nil {
			println("StateResponse unmarshal error:", err)
			return
		}

		g.HandleStateResponseMessage(stateResponse)

	case "turn_info":
		var turnInfo connection.TurnInfo
		if err := json.Unmarshal(raw, &turnInfo); err != nil {
			println("TurnInfo unmarshal error:", err)
			return
		}

		g.HandleTurnInfoMessage(turnInfo)

	default:
		println("Message with unexpected type encountered:", msg.Type)
		return
	}

	g.debugPrintln(fmt.Sprintf("Type: %v\n", msg.Type))
	g.debugPrintln(fmt.Sprintf("Data: %v\n", msg.Data))

}
