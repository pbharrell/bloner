package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pbharrell/bloner/graphics"
)

type LobbyRequestedPage struct {
	game     *Game
	fontFace *text.GoTextFace
	overlay  graphics.Shape
}

func CreateLobbyRequestedPage(g *Game) *LobbyRequestedPage {
	overlay := *graphics.CreateRectangle(OverlayImage, ScreenWidth, ScreenHeight, 0, 0, 0, 0, 0, 0)

	fontFace := &text.GoTextFace{
		Source: g.GetFontSource(),
		Size:   32,
	}

	return &LobbyRequestedPage{
		game:     g,
		fontFace: fontFace,
		overlay:  overlay,
	}
}

func (p *LobbyRequestedPage) Update() {
}

func (p *LobbyRequestedPage) Draw(screen *ebiten.Image) {
	p.overlay.Draw(screen)

	var (
		lobbyRequestedText = "Lobby requested! Waiting for response..."
		txtOp              = text.DrawOptions{}
	)

	txtW, txtH := text.Measure(lobbyRequestedText, p.fontFace, 0)
	txtOp.GeoM.Translate(ScreenWidth/2-txtW/2, ScreenHeight/2-txtH/2)
	text.Draw(screen, lobbyRequestedText, p.fontFace, &txtOp)
}
