package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	x        int
	y        int
	page     Page
	game     *Game
	fontFace *text.GoTextFace
	text     string
	textOp   text.DrawOptions
}

func CreateText(page Page, game *Game, fontFace *text.GoTextFace, text string, textOp text.DrawOptions, x int, y int) *Text {
	textOp.GeoM.Translate(float64(x), float64(y))
	return &Text{
		page:     page,
		game:     game,
		fontFace: fontFace,
		text:     text,
		textOp:   textOp,
		x:        x,
		y:        y,
	}
}

func (t *Text) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	text.Draw(screen, t.text, t.fontFace, &t.textOp)
}
