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
	gameActivePage, ok := g.pageStack[len(g.pageStack)-1].(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	gameActivePage.state.SetActiveAbsPos(PlayPos(state.ActivePlayer))
	gameActivePage.state.trumpDrawPlayer = PlayPos(state.TrumpDrawPlayer)

	if state.TrumpSuit < 0 {
		gameActivePage.state.trumpSuit = nil
	} else {
		*gameActivePage.state.trumpSuit = Suit(state.TrumpSuit)
	}

	gameActivePage.state.drawPile.Pile = make([]int, len(state.DrawPile))
	for i, card := range state.DrawPile {
		gameActivePage.state.drawPile.Pile[i] = int(card.Suit)*6 + int(card.Number)
	}

	gameActivePage.state.trick.Decode(state.PlayPile)

	gameActivePage.state.teams[Black].Decode(Black, state.TeamState[Black])
	gameActivePage.state.teams[Red].Decode(Red, state.TeamState[Red])

	gameActivePage.state.ArrangeTeams()

	// Need to recreate the turn info data if already created
	gameActivePage.state.turnInfo.inited = false
}

func (g *Game) DecodeTurnInfo(turnInfo connection.TurnInfo) {
	gameActivePage, ok := g.pageStack[len(g.pageStack)-1].(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}
	
	switch turnInfo.TurnType {
	case connection.TrumpPass:
		gameActivePage.state.passCounter++
		gameActivePage.state.SetActiveAbsPos(gameActivePage.state.GetNextPlayerById(turnInfo.PlayerId).AbsPos)
		break
	case connection.TrumpPick:
		if turnInfo.TrumpPick < 0 {
			// Don't want to repeat picking up trump, since client already did it
			if turnInfo.PlayerId != gameActivePage.state.id {
				gameActivePage.state.PickUpTrump(gameActivePage.state.GetPlayerByAbsPos(gameActivePage.state.trumpDrawPlayer))
				gameActivePage.state.SetActiveAbsPos(gameActivePage.state.trumpDrawPlayer)
			}
		} else {
			trumpSuit := Suit(turnInfo.TrumpPick)
			gameActivePage.state.trumpSuit = &trumpSuit
			gameActivePage.state.activePlayer = gameActivePage.state.trumpDrawPlayer
			gameActivePage.state.trick.clear()
		}
		break
	case connection.TrumpDiscard:
		gameActivePage.state.GetPlayerById(turnInfo.PlayerId).DiscardEncoded(turnInfo.TrumpDiscard, gameActivePage.state.id)
		gameActivePage.state.activePlayer = gameActivePage.state.trumpDrawPlayer
		break
	case connection.CardPlay:
		if turnInfo.PlayerId != gameActivePage.state.id {
			turnPlayer := gameActivePage.state.GetPlayerById(turnInfo.PlayerId)
			cardPlayed := CreateCard(Suit(turnInfo.CardPlay.Suit), Number(turnInfo.CardPlay.Number), turnInfo.PlayerId, .1, 0, 0, 0, true)
			gameActivePage.state.PlayCard(turnPlayer.Id, turnPlayer.GetCardInd(cardPlayed))
			gameActivePage.state.activePlayer = gameActivePage.state.GetNextPlayerById(turnInfo.PlayerId).AbsPos
		}
		break
	}

}
