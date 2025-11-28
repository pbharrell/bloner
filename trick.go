package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner-server/connection"
	"github.com/pbharrell/bloner/graphics"
)

type Trick struct {
	Pile     []*Card
	LeadSuit Suit
	Sprite   graphics.Sprite
	X        int
	Y        int
}

func (t *Trick) SyncSprite() {
	t.Sprite = t.Pile[len(t.Pile)-1].Sprite
	t.Sprite.X = t.X
	t.Sprite.Y = t.Y
}

func (t *Trick) Decode(pile []connection.Card) {
	t.Pile = DecodeCardPile(pile, t.Sprite.ImageScale, false)
	t.SyncSprite()
}

func (t *Trick) Encode() []connection.Card {
	encTrick := make([]connection.Card, len(t.Pile))
	for i, card := range t.Pile {
		encTrick[i] = card.Encode()
	}

	return encTrick
}

func (t *Trick) playCard(card *Card) {
	card.Sprite.Angle = t.Sprite.Angle
	card.FaceDown = false
	card.UpdateSprite()
	t.Pile = append(t.Pile, card)
	t.SyncSprite()
}

func (t *Trick) clear() {
	t.Pile = t.Pile[:0]
}

func (t *Trick) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if len(t.Pile) == 0 {
		return
	}

	t.Sprite.Draw(screen, op)
}
