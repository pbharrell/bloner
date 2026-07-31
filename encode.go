package main

import (
	"github.com/pbharrell/bloner-server/connection"
)

func (g *Game) EncodeGameState() connection.GameState {
	gameActivePage, ok := g.pageStack[len(g.pageStack)-1].(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}
	
	intTrumpSuit := -1
	if gameActivePage.state.trumpSuit != nil {
		intTrumpSuit = int(*gameActivePage.state.trumpSuit)
	}

	encDrawPile := make([]connection.Card, len(gameActivePage.state.drawPile.Pile))
	for i, cardInt := range gameActivePage.state.drawPile.Pile {
		encDrawPile[i] = CreateCard(Suit(cardInt/6), Number(cardInt%6), -1, 0, 0, 0, 0 /*faceDown*/, true).Encode()
	}

	encPlayPile := gameActivePage.state.trick.Encode()

	teamState := [2]connection.TeamState{
		gameActivePage.state.teams[Black].Encode(),
		gameActivePage.state.teams[Red].Encode(),
	}

	return connection.GameState{
		PlayerId:        gameActivePage.state.id,
		ActivePlayer:    int(gameActivePage.state.activePlayer),
		TrumpDrawPlayer: int(gameActivePage.state.trumpDrawPlayer),
		TrumpSuit:       intTrumpSuit,
		DrawPile:        encDrawPile,
		PlayPile:        encPlayPile,
		TeamState:       teamState,
	}
}
