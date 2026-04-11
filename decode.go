package main

import (
	"github.com/pbharrell/bloner-server/connection"
)

func DecodeCardPile(encPile []connection.Card, playerId int, scale float64, faceDown bool) []*Card {
	pile := make([]*Card, len(encPile))
	for i, encCard := range encPile {
		pile[i] = (CreateCard(Suit(encCard.Suit), Number(encCard.Number), playerId, scale, 0, 0, 0, faceDown))
	}

	return pile
}

func (g *Game) DecodeGameState(state connection.GameState) {
	g.SetActiveAbsPos(PlayPos(state.ActivePlayer))
	g.trumpDrawPlayer = PlayPos(state.TrumpDrawPlayer)

	if state.TrumpSuit < 0 {
		g.trumpSuit = nil
	} else {
		*g.trumpSuit = Suit(state.TrumpSuit)
	}

	g.drawPile.Pile = make([]int, len(state.DrawPile))
	for i, card := range state.DrawPile {
		g.drawPile.Pile[i] = int(card.Suit)*6 + int(card.Number)
	}

	g.trick.Decode(state.PlayPile)

	g.teams[Black].Decode(Black, state.TeamState[Black])
	g.teams[Red].Decode(Red, state.TeamState[Red])

	g.ArrangeTeams()

	// Need to recreate the turn info data if already created
	g.turnInfo.inited = false
}

func (g *Game) DecodeTurnInfo(turnInfo connection.TurnInfo) {
	switch turnInfo.TurnType {
	case connection.TrumpPass:
		g.passCounter++
		g.SetActiveAbsPos(g.GetNextPlayerById(turnInfo.PlayerId).AbsPos)
		break
	case connection.TrumpPick:
		if turnInfo.TrumpPick < 0 {
			// Don't want to repeat picking up trump, since client already did it
			if turnInfo.PlayerId != g.id {
				g.PickUpTrump(g.GetPlayerByAbsPos(g.trumpDrawPlayer))
				g.SetActiveAbsPos(g.trumpDrawPlayer)
			}
		} else {
			trumpSuit := Suit(turnInfo.TrumpPick)
			g.trumpSuit = &trumpSuit
			g.activePlayer = g.trumpDrawPlayer
			g.trick.clear()
		}
		break
	case connection.TrumpDiscard:
		g.GetPlayerById(turnInfo.PlayerId).DiscardEncoded(turnInfo.TrumpDiscard, g.id)
		g.activePlayer = g.trumpDrawPlayer
		break
	case connection.CardPlay:
		if turnInfo.PlayerId != g.id {
			turnPlayer := g.GetPlayerById(turnInfo.PlayerId)
			cardPlayed := CreateCard(Suit(turnInfo.CardPlay.Suit), Number(turnInfo.CardPlay.Number), turnInfo.PlayerId, .35, 0, 0, 0, true)
			g.PlayCard(turnPlayer.Id, turnPlayer.GetCardInd(cardPlayed))
			g.activePlayer = g.GetNextPlayerById(turnInfo.PlayerId).AbsPos
		}
		break
	}

}
