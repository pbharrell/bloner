package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner/graphics"
)

type Trick struct {
	Pile   []*Card
	Sprite graphics.Sprite
	X      int
	Y      int
}

func (t *Trick) playCard(card *Card) {
	t.Pile = append(t.Pile, card)
	t.Sprite = t.Pile[len(t.Pile)-1].Sprite
	t.Sprite.X = t.X
	t.Sprite.Y = t.Y
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
