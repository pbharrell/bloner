package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner/graphics"
)

type DrawPile struct {
	Sprite graphics.Sprite
	Pile   []int
}

func (p *DrawPile) shuffleDrawPile() {
	p.Pile = make([]int, 4*6)

	for i := range p.Pile {
		p.Pile[i] = i
	}

	for i := range p.Pile {
		j := rand.Intn(i + 1)
		p.Pile[i], p.Pile[j] = p.Pile[j], p.Pile[i]
	}
}

func (p *DrawPile) setDrawPile(pile []int) {
	p.Pile = pile
}

func (p *DrawPile) drawCard(scale float64, x int, y int, angle int) *Card {
	if len(p.Pile) == 0 {
		return nil
	}

	num := p.Pile[len(p.Pile)-1]
	suit := Suit(num / 6)
	number := Number(num % 6)

	card := CreateCard(suit, number, scale, x, y, angle)

	// Drop the last card off the pile
	p.Pile = p.Pile[:len(p.Pile)-1]

	return card
}

func (p *DrawPile) Update() {
	if len(p.Pile) == 0 {
		p.Sprite.Visible = false
	} else if !p.Sprite.Visible {
		p.Sprite.Visible = true
	}
}

func (p *DrawPile) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	p.Sprite.Draw(screen, op)
}
