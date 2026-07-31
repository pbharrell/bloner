package main

import (
	"strconv"
)

func newLobby(b *Button, _ Page, g *Game) {
	g.PushPage(CreateLobbyRequestedPage(g))
	g.SendServerMessage("lobby_req", -1)
}

func joinLobby(b *Button, _ Page, g *Game) {
	g.mode = LobbyRequestInput
	g.PushPage(CreateLobbyRequestInputPage(g))
}

func backToLobbyTypeSelect(b *Button, _ Page, g *Game) {
	g.mode = LobbyTypeSelect
	g.PopPage()
}

func joinSpecifiedLobby(_ *Button, p Page, g *Game) {
	page := p.(*LobbyRequestInputPage)
	lobbyReqId, err := strconv.ParseInt(page.LobbyRequestStr, 10, 16)
	if err != nil {
		panic("Invalid int passed for lobby request!")
	}
	g.PushPage(CreateLobbyRequestedPage(g))
	g.SendServerMessage("lobby_req", int16(lobbyReqId))
}

func confirmTrump(_ *Button, p Page, _ *Game) {
	gameActivePage, ok := p.(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	if len(gameActivePage.state.trick.Pile) < 1 {
		println("Should not be here - picked trump with an empty pile!!")
		return
	}

	gameActivePage.state.PickUpTrump(gameActivePage.state.GetPlayerByAbsPos(gameActivePage.state.trumpDrawPlayer))
	gameActivePage.state.SetActiveAbsPos(gameActivePage.state.trumpDrawPlayer)
	gameActivePage.SendTurnTrumpPick(-1)
}

func passTrump(_ *Button, p Page, g *Game) {
	gameActivePage, ok := p.(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	gameActivePage.SendTurnTrumpPass()
}

func heartsTrump(_ *Button, p Page, g *Game) {
	gameActivePage, ok := p.(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	gameActivePage.SendTurnTrumpPick(int8(Hearts))
	gameActivePage.state.trick.clear()
}

func diamondsTrump(_ *Button, p Page, g *Game) {
	gameActivePage, ok := p.(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	gameActivePage.SendTurnTrumpPick(int8(Diamonds))
	gameActivePage.state.trick.clear()
}

func clubsTrump(_ *Button, p Page, g *Game) {
	gameActivePage, ok := p.(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	gameActivePage.SendTurnTrumpPick(int8(Clubs))
	gameActivePage.state.trick.clear()
}

func spadesTrump(_ *Button, p Page, g *Game) {
	gameActivePage, ok := p.(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	gameActivePage.SendTurnTrumpPick(int8(Spades))
	gameActivePage.state.trick.clear()
}
