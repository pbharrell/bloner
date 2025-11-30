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
	g.debugPrintln("Handled lobby assign message!")
}

func (g *Game) HandleStateRequestMessage() {
	g.SendStateResponse()
}

func (g *Game) HandleStateResponseMessage(data connection.StateResponse) {
	g.mode = GameActive
	g.DecodeGameState(data)
}

func (g *Game) HandleTurnInfoMessage(data connection.TurnInfo) {
	g.DecodeTurnInfo(data)
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
			println("StateResponse unmarshal error:", err)
			return
		}

		g.HandleStateResponseMessage(stateResponse)
		break

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
