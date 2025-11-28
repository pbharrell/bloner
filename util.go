package main

import (
	_ "image/png"
)

func confirmTrump(g *Game) {
	if len(g.trick.Pile) < 1 {
		println("Should not be here - picked trump with an empty pile!!")
		return
	}

	g.PickUpTrump(g.GetActivePlayer())
	g.SendTurnTrumpPick(-1)
}

func passTrump(g *Game) {
	g.SendTurnTrumpPass()
}

func heartsTrump(g *Game) {
	g.SendTurnTrumpPick(int8(Hearts))
	g.trick.clear()
}

func diamondsTrump(g *Game) {
	g.SendTurnTrumpPick(int8(Diamonds))
	g.trick.clear()
}

func clubsTrump(g *Game) {
	g.SendTurnTrumpPick(int8(Clubs))
	g.trick.clear()
}

func spadesTrump(g *Game) {
	g.SendTurnTrumpPick(int8(Spades))
	g.trick.clear()
}

func GetRelPos(clientAbsPos PlayPos, absPos PlayPos) PlayPos {
	return (absPos - clientAbsPos) % 4
}

func GetHighestCardFromPile(cards []*Card, lead Suit, trump Suit) *Card {
	var highestCard *Card = nil
	for _, card := range cards {
		if highestCard == nil || GetHighestCard(highestCard, card, lead, trump) != highestCard {
			highestCard = card
		}
	}

	return highestCard
}

func GetHighestCard(card1 *Card, card2 *Card, lead Suit, trump Suit) *Card {
	// Correct for alt bauer suit
	// if card1.Suit == trump {
	// } if card2.Suit == trump {
	// }
	type compCard struct {
		correctedSuit   Suit
		correctedNumber Number
	}

	getCompCard := func(card *Card) compCard {
		correctedSuit := card.Suit
		correctedNumber := card.Number
		switch trump {
		case Spades:
			if card.Suit == Clubs && card.Number == Jack {
				correctedSuit = Spades
				correctedNumber = AltBauer
			}
		case Clubs:
			if card.Suit == Spades && card.Number == Jack {
				correctedSuit = Clubs
				correctedNumber = AltBauer
			}
		case Hearts:
			if card.Suit == Diamonds && card.Number == Jack {
				correctedSuit = Hearts
				correctedNumber = AltBauer
			}
		case Diamonds:
			if card.Suit == Hearts && card.Number == Jack {
				correctedSuit = Diamonds
				correctedNumber = AltBauer
			}
		}

		return compCard{
			correctedSuit:   correctedSuit,
			correctedNumber: correctedNumber,
		}
	}

	compCard1 := getCompCard(card1)
	compCard2 := getCompCard(card2)

	if compCard1.correctedSuit != compCard2.correctedSuit {
		if compCard1.correctedSuit == trump {
			return card1
		}
		if compCard2.correctedSuit == trump {
			return card2
		}
		if compCard1.correctedSuit == lead {
			return card1
		}
		if compCard2.correctedSuit == lead {
			return card2
		}
	}

	// card1's and card2's suits match - check for both trump
	if compCard1.correctedSuit == trump {
		if TrumpVal[compCard1.correctedNumber] > TrumpVal[compCard2.correctedNumber] {
			return card1
		} else if TrumpVal[compCard1.correctedNumber] < TrumpVal[compCard2.correctedNumber] {
			return card2
		} else {
			println("Comparing 2 of the same value compCards of trump suit. Returned card1.")
			return card1
		}
	}

	// regardless of whether suits match, just return higher number
	if OffVal[compCard1.correctedNumber] > OffVal[compCard2.correctedNumber] {
		return card1
	} else if OffVal[compCard1.correctedNumber] < OffVal[compCard2.correctedNumber] {
		return card2
	} else {
		println("Comparing 2 of the same value compCards numbers - could be same or different suits. Returned card1.")
		return card1
	}
}
