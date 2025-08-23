package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayPos uint8

const (
	Bottom PlayPos = iota
	Left
	Top
	Right
)

type Hand struct {
	Cards           []*Card
	PlayPos         PlayPos
	SideLen         int
	fullPercentSpan int
}

func CreateHand(handSize int, playPos PlayPos) *Hand {
	if handSize == 0 {
		return nil
	}

	cards := make([]*Card, handSize)
	for i := range cards {
		// Create card at placeholder position of 0,0
		cards[i] = CreateCard(Spades, Ace, .35, 0, 0, 0)
		//                    ^^^^^^^^^^^ replace these with DrawPile.GetCard()
	}

	var (
		sideLen         int
		percentHandSpan int
		perpAxisPos     int
	)

	switch playPos {
	case Bottom:
		sideLen = screenWidth
		percentHandSpan = 60
		perpAxisPos = screenHeight - cards[0].Sprite.ImageHeight - 20

	case Left:
		sideLen = screenHeight
		percentHandSpan = 20
		perpAxisPos = 20

	case Top:
		sideLen = screenWidth
		percentHandSpan = 20
		perpAxisPos = 20

	case Right:
		sideLen = screenHeight
		percentHandSpan = 20
		perpAxisPos = screenWidth - cards[0].Sprite.ImageHeight - 20
	}

	// Calculate and set the card position
	ArrangeHand(cards, playPos, sideLen, percentHandSpan, perpAxisPos)

	return &Hand{
		Cards:           cards,
		PlayPos:         playPos,
		SideLen:         sideLen,
		fullPercentSpan: percentHandSpan,
	}
}

func (h *Hand) Update() {
	for i := range h.Cards {
		h.Cards[i].Update()
	}
}

func (h *Hand) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	for i := range h.Cards {
		h.Cards[i].Draw(screen, op)
	}
}

func ArrangeHand(cards []*Card, playPos PlayPos, sideLen int,
	percentHandSpan int, perpAxisHeight int) {
	// Assume that all the cards are of the same width
	numCards := len(cards)
	cardWidth := cards[0].Sprite.ImageWidth

	handSpan := int(float32(percentHandSpan) * .01 * float32(sideLen))
	cardMargin := int(float32(handSpan-numCards*cardWidth) / float32(numCards))

	handStart := int(float32(sideLen-int(len(cards)*cardWidth+max(0, len(cards)-1)*cardMargin)) / 2)

	for cardInd := range cards {
		playAxisPos := handStart + (cardInd * (cardWidth + cardMargin))
		switch playPos {
		case Bottom:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = perpAxisHeight
			cards[cardInd].Sprite.Angle = 0
		case Left:
			cards[cardInd].Sprite.X = perpAxisHeight
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 90
		case Top:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = perpAxisHeight
			cards[cardInd].Sprite.Angle = 180
		case Right:
			cards[cardInd].Sprite.X = perpAxisHeight
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 270
		}

	}
}
