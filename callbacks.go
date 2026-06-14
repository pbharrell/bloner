package main

import (
	"strconv"
)

func newLobby(b *Button, _ Page, g *Game) {
	g.SendServerMessage("lobby_req", -1)
}

func joinLobby(b *Button, _ Page, g *Game) {
	g.mode = LobbyRequestInput
	g.PushPage(CreateLobbyRequestInputPage(g))
}

func backToLobbyTypeSelect(b *Button, _ Page, g *Game) {
	g.mode = LobbyTypeSelect
	g.PushPage(CreateLobbyTypeSelectPage(g))
}

func joinSpecifiedLobby(_ *Button, p Page, g *Game) {
	page := p.(*LobbyRequestInputPage)
	lobbyReqId, err := strconv.ParseInt(page.LobbyRequestStr, 10, 16)
	if err != nil {
		panic("Invalid int passed for lobby request!")
	}
	g.SendServerMessage("lobby_req", int16(lobbyReqId))
}

func confirmTrump(_ *Button, _ Page, g *Game) {
	if len(g.GetTrick().Pile) < 1 {
		println("Should not be here - picked trump with an empty pile!!")
		return
	}

	g.PickUpTrump(g.GetPlayerByAbsPos(g.trumpDrawPlayer))
	g.SetActiveAbsPos(g.trumpDrawPlayer)
	g.SendTurnTrumpPick(-1)
}

func passTrump(_ *Button, _ Page, g *Game) {
	g.SendTurnTrumpPass()
}

func heartsTrump(_ *Button, _ Page, g *Game) {
	g.SendTurnTrumpPick(int8(Hearts))
	g.trick.clear()
}

func diamondsTrump(_ *Button, _ Page, g *Game) {
	g.SendTurnTrumpPick(int8(Diamonds))
	g.trick.clear()
}

func clubsTrump(_ *Button, _ Page, g *Game) {
	g.SendTurnTrumpPick(int8(Clubs))
	g.trick.clear()
}

func spadesTrump(_ *Button, _ Page, g *Game) {
	g.SendTurnTrumpPick(int8(Spades))
	g.trick.clear()
}
