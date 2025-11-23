package main

import (
	"slices"

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

func CreatePlayer(id int, team teamColor, handSize int, relPos PlayPos, scale float64, drawPile *DrawPile, faceDown bool) Player {
	if handSize == 0 {
		return Player{}
	}

	cards := make([]*Card, handSize)
	for i := range cards {
		if drawPile != nil {
			cards[i] = drawPile.drawCard(scale, 0, 0, 0, faceDown)
		} else {
			cards[i] = CreateCard(Spades, Ace, .35, 0, 0, 0, faceDown)
		}
	}

	player := Player{
		Id:      id,
		Cards:   cards,
		AbsPos:  relPos,
		RelPos:  relPos,
		PosInfo: GetPosInfoFromPos(relPos, cards[0].Sprite.ImageHeight),
	}

	// Calculate and set the card position
	player.Arrange(0, Bottom)

	return player
}

func (p *Player) Arrange(clientId int, clientPos PlayPos) {
	p.RelPos = PlayPos((uint8(p.AbsPos) - uint8(clientPos)) % 4)
	p.PosInfo = GetPosInfoFromPos(p.RelPos, p.PosInfo.cardHeight)
	p.ArrangeHand(clientId)
}

func (p *Player) Decode(teamColor teamColor, playerNum uint8, playerState connection.PlayerState) {
	// Face down value should be overridden by `ArrangeHand` later
	p.Id = playerState.PlayerId
	p.Cards = DecodeCardPile(playerState.Cards, .35 /*faceDown*/, true)
	p.AbsPos = PlayPos(uint8(teamColor) + playerNum*2)
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

func (p *Player) GetCardInd(card *Card) int {
	for i, c := range p.Cards {
		if c.Suit == card.Suit && c.Number == card.Number {
			return i
		}
	}

	println("`GetCardInd()` called with no matching card for that hand!")
	return -1
}

func (p *Player) Discard(card int) *Card {
	if card >= len(p.Cards) {
		return nil
	}

	discarded := p.Cards[card]
	p.Cards = slices.Delete(p.Cards, card, card+1)
	p.ArrangeHand(p.Id)
	return discarded
}

func (p *Player) DiscardEncoded(card connection.Card) *Card {
	for i, c := range p.Cards {
		if c.Suit == Suit(card.Suit) && c.Number == Number(card.Number) {
			return p.Discard(i)
		}
	}

	return nil
}

func (p *Player) Update() {
	for i := range p.Cards {
		p.Cards[i].UpdateSprite()
	}
}

func (p *Player) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	for i := range p.Cards {
		p.Cards[i].Draw(screen, op)
	}
}

func (p *Player) ArrangeHand(clientId int) {
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
		cards[cardInd].FaceDown = clientId != p.Id
		cards[cardInd].UpdateSprite()

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
