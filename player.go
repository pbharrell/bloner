package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/pbharrell/bloner-server/connection"
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

type Player struct {
	Id              int
	Cards           []*Card
	AbsPos          PlayPos
	RelPos          PlayPos
	SideLen         int
	fullPercentSpan int
	perpAxisPos     int
	tricksWon       int
}

func CreatePlayer(id int, team teamColor, handSize int, relPos PlayPos, scale float64, drawPile *DrawPile) Player {
	if handSize == 0 {
		return Player{}
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

	switch relPos {
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

	hand := Player{
		Id:              id,
		Cards:           cards,
		AbsPos:          relPos, // FIXME:
		RelPos:          relPos,
		SideLen:         sideLen,
		fullPercentSpan: percentHandSpan,
		perpAxisPos:     perpAxisPos,
	}

	// Calculate and set the card position
	hand.ArrangeHand()

	return hand
}

func (p *Player) Decode(playerState connection.PlayerState) {
	p.Cards = DecodeCardPile(playerState.Cards, .35)
}

func (p *Player) Encode() connection.PlayerState {
	playerId := p.Id
	cards := make([]connection.Card, len(p.Cards))
	for k := range cards {
		cards[k] = GetEncodedCard(p.Cards[k])
	}

	return connection.PlayerState{
		PlayerId: playerId,
		Cards:    cards,
	}
}

func (p *Player) GetTeam() teamColor {
	return teamColor(uint8(p.AbsPos) % 2)
}

func (p *Player) Update() {
	for i := range p.Cards {
		p.Cards[i].Update()
	}
}

func (p *Player) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	for i := range p.Cards {
		p.Cards[i].Draw(screen, op)
	}
}

func (p *Player) ArrangeHand() {
	cards := p.Cards
	sideLen := p.SideLen

	// Assume that all the cards are of the same width
	numCards := len(cards)
	if numCards == 0 {
		return
	}

	var percentHandSpan int
	if len(cards) <= len(maxSpan) {
		percentHandSpan = min(p.fullPercentSpan, maxSpan[len(cards)-1])
	} else {
		percentHandSpan = p.fullPercentSpan
	}

	cardWidth := cards[0].Sprite.ImageWidth

	handSpan := int(float32(percentHandSpan) * .01 * float32(sideLen))
	cardMargin := int(float32(handSpan-numCards*cardWidth) / float32(numCards))

	handStart := int(float32(sideLen-int(len(cards)*cardWidth+max(0, len(cards)-1)*cardMargin)) / 2)

	for cardInd := range cards {
		playAxisPos := handStart + (cardInd * (cardWidth + cardMargin))
		switch p.RelPos {
		case Bottom:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = p.perpAxisPos
			cards[cardInd].Sprite.Angle = 0
		case Left:
			cards[cardInd].Sprite.X = p.perpAxisPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 90
		case Top:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = p.perpAxisPos
			cards[cardInd].Sprite.Angle = 180
		case Right:
			cards[cardInd].Sprite.X = p.perpAxisPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 270
		}

	}
}
