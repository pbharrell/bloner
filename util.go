package main

import (
	_ "image/png"
)

func confirmTrump(g *Game) {
	if len(g.trick.Pile) < 1 {
		println("Should not be here - picked trump with an empty pile!!")
		return
	}

	topCard := g.trick.Pile[len(g.trick.Pile)-1]
	g.trick.Pile = g.trick.Pile[:len(g.trick.Pile)-1]
	g.trumpSuit = &topCard.Suit
	g.player.Cards = append(g.player.Cards, topCard)
	g.player.ArrangeHand()
}

func cancelTrump(g *Game) {
	g.EndTurn()
}
