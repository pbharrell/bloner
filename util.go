package main

import (
	_ "image/png"
	"github.com/pbharrell/bloner/graphics"
)

func CreateSpriteFromAssetPath(assetPath string, scale float64, x int, y int, angle int, vx int, vy int, vangle int) *graphics.Sprite {
	image, alphaImage := graphics.LoadImageFromFile(&content, assetPath)
	return graphics.CreateSprite(image, alphaImage, scale, x, y, angle, 0, 0, 0)
}

func GetHighestCardFromPile(cards []*Card, lead Suit, trump Suit) *Card {
	println("Trump suit:", SuitToString(trump))
	var highestCard *Card = nil
	for _, card := range cards {
		if highestCard == nil || GetHighestCard(highestCard, card, lead, trump) != highestCard {
			highestCard = card
			println("New highest card unlocked:", SuitToString(highestCard.Suit), NumberToString(highestCard.Number))
		}
	}

	return highestCard
}

func GetHighestCard(card1 *Card, card2 *Card, lead Suit, trump Suit) *Card {

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

	println("Comparing cards card1:", SuitToString(card1.Suit), NumberToString(card1.Number), "from player:", card1.PlayerId)
	println("Comparing cards card2:", SuitToString(card2.Suit), NumberToString(card2.Number), "from player:", card2.PlayerId)

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
