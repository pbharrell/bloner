package main

import (
	_ "image/png"

	"github.com/pbharrell/bloner-server/connection"
)

func pickUpTrump(g *Game) {
	if len(g.trick.Pile) < 1 {
		println("Should not be here - picked trump with an empty pile!!")
		return
	}

	topCard := g.trick.Pile[len(g.trick.Pile)-1]
	g.trick.Pile = g.trick.Pile[:len(g.trick.Pile)-1]
	g.trumpSuit = &topCard.Suit

	player := g.GetPlayer(g.activePlayer)
	player.Cards = append(player.Cards, topCard)
	player.ArrangeHand(player.Id)

	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = -1
	g.SendTurnInfo()
	g.EndTurn()
}

func passTrump(g *Game) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPass
	g.SendTurnInfo()
	g.EndTurn()
}

func heartsTrump(g *Game) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = int8(Hearts)
	g.SendTurnInfo()
	g.EndTurn()
}

func diamondsTrump(g *Game) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = int8(Diamonds)
	g.SendTurnInfo()
	g.EndTurn()
}

func clubsTrump(g *Game) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = int8(Clubs)
	g.SendTurnInfo()
	g.EndTurn()
}

func spadesTrump(g *Game) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = int8(Spades)
	g.SendTurnInfo()
	g.EndTurn()
}

func GetRelPos(clientAbsPos PlayPos, absPos PlayPos) PlayPos {
	return (absPos - clientAbsPos) % 4
}
