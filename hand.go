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

var (
	maxSpan = [5]int{0, 25, 40, 50, 60}
)

type Hand struct {
	Cards           []*Card
	PlayPos         PlayPos
	SideLen         int
	fullPercentSpan int
	perpAxisPos     int
	tricksWon       int
}

func CreateHand(handSize int, playPos PlayPos, scale float64, drawPile *DrawPile) *Hand {
	if handSize == 0 {
		return nil
	}

	cards := make([]*Card, handSize)
	for i := range cards {
		if drawPile != nil {
			cards[i] = drawPile.drawCard(scale, 0, 0, 0)
		} else {
			cards[i] = CreateCard(Spades, Ace, .35, 0, 0, 0)
		}
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

	hand := &Hand{
		Cards:           cards,
		PlayPos:         playPos,
		SideLen:         sideLen,
		fullPercentSpan: percentHandSpan,
		perpAxisPos:     perpAxisPos,
	}

	// Calculate and set the card position
	hand.ArrangeHand()

	return hand
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

func (h *Hand) ArrangeHand() {
	cards := h.Cards
	sideLen := h.SideLen

	// Assume that all the cards are of the same width
	numCards := len(cards)
	if numCards == 0 {
		return
	}

	var percentHandSpan int
	if len(cards) <= len(maxSpan) {
		percentHandSpan = min(h.fullPercentSpan, maxSpan[len(cards)-1])
	} else {
		percentHandSpan = h.fullPercentSpan
	}

	cardWidth := cards[0].Sprite.ImageWidth

	handSpan := int(float32(percentHandSpan) * .01 * float32(sideLen))
	cardMargin := int(float32(handSpan-numCards*cardWidth) / float32(numCards))

	handStart := int(float32(sideLen-int(len(cards)*cardWidth+max(0, len(cards)-1)*cardMargin)) / 2)

	for cardInd := range cards {
		playAxisPos := handStart + (cardInd * (cardWidth + cardMargin))
		switch h.PlayPos {
		case Bottom:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = h.perpAxisPos
			cards[cardInd].Sprite.Angle = 0
		case Left:
			cards[cardInd].Sprite.X = h.perpAxisPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 90
		case Top:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = h.perpAxisPos
			cards[cardInd].Sprite.Angle = 180
		case Right:
			cards[cardInd].Sprite.X = h.perpAxisPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 270
		}

	}
}
