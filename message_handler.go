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
	g.debugPrintln("Handled lobby assign message!")
}

func (g *Game) HandleStateRequestMessage() {
	gameState := g.EncodeGameState()
	fmt.Printf("%v", gameState)
	g.server.Send(connection.Message{
		Type: "state_res",
		Data: gameState,
	})

	g.debugPrintln("Handled state request message!")
}

func (g *Game) HandleStateResponseMessage(data connection.GameState) {
	g.DecodeGameState(data)
}

func (g *Game) HandleMessage(msg connection.Message) {
	// Marshal Data back into JSON bytes
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
		break
	case "state_req":
		g.HandleStateRequestMessage()
		break

	case "state_res":
		var stateResponse connection.StateResponse
		if err := json.Unmarshal(raw, &stateResponse); err != nil {
			println("LobbyAssign unmarshal error:", err)
			return
		}

		g.HandleStateResponseMessage(stateResponse)
		break

	default:
		println("Message with unexpected type encountered:", msg.Type)
		return
	}

	g.debugPrintln(fmt.Sprintf("Type: %v\n", msg.Type))
	g.debugPrintln(fmt.Sprintf("Data: %v\n", msg.Data))

}
