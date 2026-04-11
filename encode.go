package main

import (
	"github.com/pbharrell/bloner-server/connection"
)

func (g *Game) EncodeGameState() connection.GameState {
	intTrumpSuit := -1
	if g.trumpSuit != nil {
		intTrumpSuit = int(*g.trumpSuit)
	}

	encDrawPile := make([]connection.Card, len(g.drawPile.Pile))
	for i, cardInt := range g.drawPile.Pile {
		encDrawPile[i] = CreateCard(Suit(cardInt/6), Number(cardInt%6), -1, 0, 0, 0, 0 /*faceDown*/, true).Encode()
	}

	encPlayPile := g.trick.Encode()

	teamState := [2]connection.TeamState{
		g.teams[Black].Encode(),
		g.teams[Red].Encode(),
	}

	return connection.GameState{
		PlayerId:        g.id,
		ActivePlayer:    int(g.activePlayer),
		TrumpDrawPlayer: int(g.trumpDrawPlayer),
		TrumpSuit:       intTrumpSuit,
		DrawPile:        encDrawPile,
		PlayPile:        encPlayPile,
		TeamState:       teamState,
	}
}
