package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pbharrell/bloner/graphics"
)

type LobbyAssignedPage struct {
	game     *Game
	fontFace *text.GoTextFace
	overlay  graphics.Shape
}

func CreateLobbyAssignedPage(g *Game) *LobbyAssignedPage {
	overlay := *graphics.CreateRectangle(overlayImage, screenWidth, screenHeight, 0, 0, 0, 0, 0, 0)

	fontFace := &text.GoTextFace{
		Source: g.fontSource,
		Size:   32,
	}

	return &LobbyAssignedPage{
		game:     g,
		fontFace: fontFace,
		overlay:  overlay,
	}
}

func (p *LobbyAssignedPage) Update() {
}

func (p *LobbyAssignedPage) Draw(screen *ebiten.Image) {
	p.overlay.Draw(screen)

	var (
		lobbyAssignedText  = fmt.Sprintf("Lobby found with id: %v!", p.game.lobbyId)
		lobbyWaitingText   = fmt.Sprintf("Waiting on more players...")
		lobbyAssignedTxtOp = text.DrawOptions{}
		lobbyWaitingTxtOp  = text.DrawOptions{}
	)

	lobbyAssignedTxtW, lobbyAssignedTxtH := text.Measure(lobbyAssignedText, p.fontFace, 0)
	lobbyWaitingTxtW, lobbyWaitingTxtH := text.Measure(lobbyWaitingText, p.fontFace, 0)
	lobbyAssignedTxtOp.GeoM.Translate(screenWidth/2-lobbyAssignedTxtW/2, screenHeight/2-lobbyAssignedTxtH/2)
	text.Draw(screen, lobbyAssignedText, p.fontFace, &lobbyAssignedTxtOp)

	lobbyWaitingTxtOp.GeoM.Translate(screenWidth/2-lobbyWaitingTxtW/2, screenHeight/2+lobbyWaitingTxtH/2)
	text.Draw(screen, lobbyWaitingText, p.fontFace, &lobbyWaitingTxtOp)
}
