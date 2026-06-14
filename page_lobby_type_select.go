package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pbharrell/bloner/graphics"
)

type LobbyTypeSelectPage struct {
	game            *Game
	buttonNewLobby  *Button
	buttonJoinLobby *Button
	overlay         graphics.Shape
}

func CreateLobbyTypeSelectPage(game *Game) *LobbyTypeSelectPage {
	overlay := *graphics.CreateRectangle(overlayImage, screenWidth, screenHeight, 0, 0, 0, 0, 0, 0)

	p := &LobbyTypeSelectPage{
		game:    game,
		overlay: overlay,
	}

	p.buttonNewLobby = CreateButton(p, game, newLobby, "assets/new_lobby_button.png", "assets/new_lobby_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	newLobbyWidth := p.buttonNewLobby.sprite.ImageWidth
	newLobbyX := screenWidth/2 - newLobbyWidth/2 - 80
	p.buttonNewLobby.SetLoc(newLobbyX, p.buttonNewLobby.sprite.Y)

	p.buttonJoinLobby = CreateButton(p, game, joinLobby, "assets/join_lobby_button.png", "assets/join_lobby_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	joinLobbyWidth := p.buttonJoinLobby.sprite.ImageWidth
	joinLobbyX := screenWidth/2 - joinLobbyWidth/2 + 80
	p.buttonJoinLobby.SetLoc(joinLobbyX, p.buttonJoinLobby.sprite.Y)

	return p
}

func (p *LobbyTypeSelectPage) Update() {
	x, y := ebiten.CursorPosition()
	mouseButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	p.buttonNewLobby.Update(x, y, mouseButtonPressed)
	p.buttonJoinLobby.Update(x, y, mouseButtonPressed)
}

func (p *LobbyTypeSelectPage) Draw(screen *ebiten.Image) {
	p.overlay.Draw(screen)

	var (
		lobbyWaitText = "Create new game or join existing?"
		txtOp         = text.DrawOptions{}
		op            = ebiten.DrawImageOptions{}
	)

	fontFace := &text.GoTextFace{
		Source: p.game.GetFontSource(),
		Size:   32,
	}

	txtW, txtH := text.Measure(lobbyWaitText, fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2)
	text.Draw(screen, lobbyWaitText, fontFace, &txtOp)

	p.buttonNewLobby.Draw(screen, op)
	p.buttonJoinLobby.Draw(screen, op)
}
