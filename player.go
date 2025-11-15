package main

import (
	"fmt"

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

type PosInfo struct {
	cardHeight      int
	SideLen         int
	fullPercentSpan int
	perpAxisPos     int
}

type Player struct {
	Id        int
	Cards     []*Card
	AbsPos    PlayPos
	RelPos    PlayPos
	PosInfo   PosInfo
	tricksWon int
}

func GetPosInfoFromPos(relPos PlayPos, cardHeight int) PosInfo {
	var (
		sideLen         int
		percentHandSpan int
		perpAxisPos     int
	)

	switch relPos {
	case Bottom:
		sideLen = screenWidth
		percentHandSpan = 60
		perpAxisPos = screenHeight - cardHeight - 20

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
		perpAxisPos = screenWidth - cardHeight - 20
	}

	return PosInfo{
		cardHeight:      cardHeight,
		SideLen:         sideLen,
		fullPercentSpan: percentHandSpan,
		perpAxisPos:     perpAxisPos,
	}
}

func CreatePlayer(id int, team teamColor, handSize int, relPos PlayPos, scale float64, drawPile *DrawPile) Player {
	if handSize == 0 {
		return Player{}
	}

	cards := make([]*Card, handSize)
	for i := range cards {
		if drawPile != nil {
			cards[i] = drawPile.drawCard(scale, 0, 0, 0)
			// println("player with id:", id, "drew:", NumberToString(cards[i].Number), "of", SuitToString(cards[i].Suit))
		} else {
			cards[i] = CreateCard(Spades, Ace, .35, 0, 0, 0)
		}
	}

	hand := Player{
		Id:      id,
		Cards:   cards,
		AbsPos:  relPos,
		RelPos:  relPos,
		PosInfo: GetPosInfoFromPos(relPos, cards[0].Sprite.ImageHeight),
	}

	// Calculate and set the card position
	hand.Arrange(Bottom)

	return hand
}

func (p *Player) Arrange(clientPos PlayPos) {
	p.RelPos = PlayPos((uint8(p.AbsPos) - uint8(clientPos)) % 4)
	p.PosInfo = GetPosInfoFromPos(p.RelPos, p.PosInfo.cardHeight)
	p.ArrangeHand()

	fmt.Printf("Pos calc for player w/ id: %v\n", p.Id)
	fmt.Printf("Abs pos: %v\n", p.AbsPos)
	fmt.Printf("Client pos: %v\n", clientPos)
	fmt.Printf("Resultant rel pos: %v\n\n", p.RelPos)
}

func (p *Player) Decode(teamColor teamColor, playerNum uint8, playerState connection.PlayerState) {
	p.Cards = DecodeCardPile(playerState.Cards, .35)
	p.AbsPos = PlayPos(uint8(teamColor)*2 + playerNum)
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
	sideLen := p.PosInfo.SideLen

	// Assume that all the cards are of the same width
	numCards := len(cards)
	if numCards == 0 {
		return
	}

	var percentHandSpan int
	if len(cards) <= len(maxSpan) {
		percentHandSpan = min(p.PosInfo.fullPercentSpan, maxSpan[len(cards)-1])
	} else {
		percentHandSpan = p.PosInfo.fullPercentSpan
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
			cards[cardInd].Sprite.Y = p.PosInfo.perpAxisPos
			cards[cardInd].Sprite.Angle = 0
		case Left:
			cards[cardInd].Sprite.X = p.PosInfo.perpAxisPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 90
		case Top:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = p.PosInfo.perpAxisPos
			cards[cardInd].Sprite.Angle = 180
		case Right:
			cards[cardInd].Sprite.X = p.PosInfo.perpAxisPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 270
		}

	}
}
