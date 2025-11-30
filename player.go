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
	cardHeight        int
	SideLen           int
	fullPercentSpan   int
	perpAxisHandPos   int
	perpAxisTricksPos int
}

type Player struct {
	Id        int
	Cards     []*Card
	AbsPos    PlayPos
	RelPos    PlayPos
	PosInfo   PosInfo
	wonTricks []*Card
}

func GetPosInfoFromPos(relPos PlayPos, cardHeight int) PosInfo {
	var (
		sideLen           int
		percentHandSpan   int
		perpAxisHandPos   int
		perpAxisTricksPos int
	)

	switch relPos {
	case Bottom:
		sideLen = screenWidth
		percentHandSpan = 60
		perpAxisHandPos = screenHeight - cardHeight - 20
		perpAxisTricksPos = screenHeight - cardHeight - 20

	case Left:
		sideLen = screenHeight
		percentHandSpan = 20
		perpAxisHandPos = 20
		perpAxisTricksPos = 20

	case Top:
		sideLen = screenWidth
		percentHandSpan = 20
		perpAxisHandPos = 20
		perpAxisTricksPos = 20

	case Right:
		sideLen = screenHeight
		percentHandSpan = 20
		perpAxisHandPos = screenWidth - cardHeight - 20
		perpAxisTricksPos = screenWidth - cardHeight - 20
	}

	return PosInfo{
		cardHeight:        cardHeight,
		SideLen:           sideLen,
		fullPercentSpan:   percentHandSpan,
		perpAxisHandPos:   perpAxisHandPos,
		perpAxisTricksPos: perpAxisTricksPos,
	}
}

func CreatePlayer(id int, team teamColor, handSize int, relPos PlayPos, scale float64, drawPile *DrawPile, faceDown bool, tricksWon uint8) Player {
	if handSize == 0 {
		return Player{}
	}

	tricks := make([]*Card, tricksWon)
	for i := range tricks {
		// Card suit and number are placholders
		tricks[i] = CreateCard(Spades, Ace, id, .35, 0, 0, 0 /* faceDown */, true)
	}

	player := Player{
		Id:        id,
		AbsPos:    relPos,
		RelPos:    relPos,
		wonTricks: tricks,
	}

	player.DealHand(scale, drawPile, handSize, faceDown)

	var imageHeight int
	if len(player.Cards) > 0 {
		imageHeight = player.Cards[0].Sprite.ImageHeight
	} else {
		// No cards, the best we can do is use a random image
		imageHeight = CreateCard(Spades, Ace, 0, scale, 0, 0, 0, false).Sprite.ImageWidth
	}

	player.PosInfo = GetPosInfoFromPos(player.RelPos, imageHeight)

	// Calculate and set the card position
	player.Arrange(0, relPos)

	return player
}

func (p *Player) DealHand(scale float64, drawPile *DrawPile, handSize int, faceDown bool) {
	p.Cards = make([]*Card, handSize)
	for i := range p.Cards {
		if drawPile != nil {
			p.Cards[i] = drawPile.drawCard(scale, 0, 0, 0, faceDown)
			if p.Cards[i] == nil {
				println("Drew card but somehow wound up with nil")
			} else {
				println("Settings drawn card id to:", p.Id)
				p.Cards[i].PlayerId = p.Id
			}
		} else {
			p.Cards[i] = CreateCard(Spades, Ace, p.Id, .35, 0, 0, 0, faceDown)
		}
	}
}

func (p *Player) Arrange(clientId int, clientPos PlayPos) {
	p.RelPos = PlayPos((uint8(p.AbsPos) - uint8(clientPos)) % 4)
	p.PosInfo = GetPosInfoFromPos(p.RelPos, p.PosInfo.cardHeight)
	p.ArrangeHand(clientId)
	p.ArrangeTricks()
}

func (p *Player) Decode(teamColor teamColor, playerNum uint8, playerState connection.PlayerState) {
	// Face down value should be overridden by `ArrangeHand` later
	p.Id = playerState.PlayerId
	p.Cards = DecodeCardPile(playerState.Cards, playerState.PlayerId, .35 /*faceDown*/, true)
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

func (p *Player) Discard(card int, clientId int) *Card {
	if card >= len(p.Cards) {
		return nil
	}

	discarded := p.Cards[card]
	p.Cards = slices.Delete(p.Cards, card, card+1)
	p.ArrangeHand(clientId)
	return discarded
}

func (p *Player) DiscardEncoded(card connection.Card, clientId int) *Card {
	for i, c := range p.Cards {
		if c.Suit == Suit(card.Suit) && c.Number == Number(card.Number) {
			return p.Discard(i, clientId)
		}
	}

	return nil
}

func (p *Player) WinTrick(clientId int) {
	p.wonTricks = append(p.wonTricks, CreateCard(Spades, Ace, p.Id, .35, 0, 0, 0, true))
	p.ArrangeTricks()
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

	for j := range p.wonTricks {
		p.wonTricks[j].Draw(screen, op)
	}
}

func (p *Player) ArrangeHand(clientId int) {
	cards := p.Cards
	sideLen := p.PosInfo.SideLen

	// // Assume that all the cards are of the same width
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
			cards[cardInd].Sprite.Y = p.PosInfo.perpAxisHandPos
			cards[cardInd].Sprite.Angle = 0
		case Left:
			cards[cardInd].Sprite.X = p.PosInfo.perpAxisHandPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 90
		case Top:
			cards[cardInd].Sprite.X = playAxisPos
			cards[cardInd].Sprite.Y = p.PosInfo.perpAxisHandPos
			cards[cardInd].Sprite.Angle = 180
		case Right:
			cards[cardInd].Sprite.X = p.PosInfo.perpAxisHandPos
			cards[cardInd].Sprite.Y = playAxisPos
			cards[cardInd].Sprite.Angle = 270
		}

	}
}

func (p *Player) ArrangeTricks() {
	tricks := p.wonTricks
	if len(tricks) == 0 {
		return
	}

	trickWidth := tricks[0].Sprite.ImageWidth
	for trickInd := range tricks {
		tricks[trickInd].FaceDown = true

		trickAxisOffset := trickInd * (trickWidth - 40)
		switch p.RelPos {
		case Bottom:
			tricks[trickInd].Sprite.X = 20 + trickAxisOffset
			tricks[trickInd].Sprite.Y = p.PosInfo.perpAxisTricksPos
			tricks[trickInd].Sprite.Angle = 0
		case Left:
			tricks[trickInd].Sprite.X = p.PosInfo.perpAxisTricksPos
			tricks[trickInd].Sprite.Y = 20 + trickAxisOffset
			tricks[trickInd].Sprite.Angle = 90
		case Top:
			tricks[trickInd].Sprite.X = screenWidth - trickWidth - 20 - trickAxisOffset
			tricks[trickInd].Sprite.Y = p.PosInfo.perpAxisTricksPos
			tricks[trickInd].Sprite.Angle = 180
		case Right:
			tricks[trickInd].Sprite.X = p.PosInfo.perpAxisTricksPos
			tricks[trickInd].Sprite.Y = screenHeight - p.PosInfo.cardHeight - 20 - trickAxisOffset
			tricks[trickInd].Sprite.Angle = 270
		}
	}
}
