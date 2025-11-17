package main

import (
	_ "image/png"

	"github.com/pbharrell/bloner-server/connection"
)

func confirmTrump(g *Game) {
	if len(g.trick.Pile) < 1 {
		println("Should not be here - picked trump with an empty pile!!")
		return
	}

	topCard := g.trick.Pile[len(g.trick.Pile)-1]
	g.trick.Pile = g.trick.Pile[:len(g.trick.Pile)-1]
	g.trumpSuit = &topCard.Suit
	g.GetClient().Cards = append(g.GetClient().Cards, topCard)
	g.GetClient().ArrangeHand(g.GetClient().Id)

	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = topCard.Encode()
	g.SendTurnInfo()
}

func cancelTrump(g *Game) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPass
	g.SendTurnInfo()
	g.EndTurn()
}

func GetRelPos(clientAbsPos PlayPos, absPos PlayPos) PlayPos {
	return (absPos - clientAbsPos) % 4
}
