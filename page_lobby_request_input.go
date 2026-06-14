package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pbharrell/bloner/graphics"
)

type LobbyRequestInputPage struct {
	game            *Game
	buttonBack      *Button
	buttonJoinLobby *Button
	LobbyRequestStr string
	txtInputBox     *graphics.Shape
	txtOp           text.DrawOptions
	fontFace        text.GoTextFace
	overlay         graphics.Shape
}

func CreateLobbyRequestInputPage(game *Game) *LobbyRequestInputPage {
	lobbyRequestStr := "Id of lobby you'd like to join:"

	overlay := *graphics.CreateRectangle(overlayImage, screenWidth, screenHeight, 0, 0, 0, 0, 0, 0)

	fontFace := text.GoTextFace{
		Source: game.GetFontSource(),
		Size:   32,
	}

	txtOp := text.DrawOptions{}

	lobbyRequestInputTxtW, lobbyRequestInputTxtH := text.Measure(lobbyRequestStr, &fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2-lobbyRequestInputTxtW/2, screenHeight/2-lobbyRequestInputTxtH)

	txtInputBox := graphics.CreateRectangle(txtInputBoxImage, int(lobbyRequestInputTxtW)+10, int(lobbyRequestInputTxtH)+10, int(screenWidth/2-lobbyRequestInputTxtW/2-5), int(screenHeight/2+lobbyRequestInputTxtH-5), 0, 0, 0, 0)

	p := &LobbyRequestInputPage{
		game:            game,
		LobbyRequestStr: lobbyRequestStr,
		txtInputBox:     txtInputBox,
		txtOp:           txtOp,
		fontFace:        fontFace,
		overlay:         overlay,
	}

	p.buttonBack = CreateButton(p, game, newLobby, "assets/new_lobby_button.png", "assets/new_lobby_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	newLobbyWidth := p.buttonBack.sprite.ImageWidth
	newLobbyX := screenWidth/2 - newLobbyWidth/2 - 80
	p.buttonBack.SetLoc(newLobbyX, p.buttonBack.sprite.Y)

	p.buttonJoinLobby = CreateButton(p, game, joinLobby, "assets/join_lobby_button.png", "assets/join_lobby_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	joinLobbyWidth := p.buttonJoinLobby.sprite.ImageWidth
	joinLobbyX := screenWidth/2 - joinLobbyWidth/2 + 80
	p.buttonJoinLobby.SetLoc(joinLobbyX, p.buttonJoinLobby.sprite.Y)

	return p
}

func (p *LobbyRequestInputPage) Update() {
	x, y := ebiten.CursorPosition()
	mouseButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	p.buttonBack.Update(x, y, mouseButtonPressed)
	p.buttonJoinLobby.Update(x, y, mouseButtonPressed)

	if len(p.LobbyRequestStr) < 22 {
		if inpututil.IsKeyJustPressed(ebiten.Key0) {
			p.LobbyRequestStr += "0"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			p.LobbyRequestStr += "1"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			p.LobbyRequestStr += "2"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			p.LobbyRequestStr += "3"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			p.LobbyRequestStr += "4"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key5) {
			p.LobbyRequestStr += "5"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key6) {
			p.LobbyRequestStr += "6"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key7) {
			p.LobbyRequestStr += "7"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key8) {
			p.LobbyRequestStr += "8"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key9) {
			p.LobbyRequestStr += "9"
		}
	}

	if len(p.LobbyRequestStr) > 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			p.LobbyRequestStr = p.LobbyRequestStr[:len(p.LobbyRequestStr)-1]
		}
	}
}

func (p *LobbyRequestInputPage) Draw(screen *ebiten.Image) {
	p.overlay.Draw(screen)

	var (
		op                       = ebiten.DrawImageOptions{}
		lobbyRequestInputTxt     = "Id of lobby you'd like to join:"
		lobbyRequestInputTxtOp   = text.DrawOptions{}
		lobbyRequstInputTxtBoxOp = text.DrawOptions{}
	)

	fontFace := &text.GoTextFace{
		Source: p.game.GetFontSource(),
		Size:   32,
	}

	lobbyRequestInputTxtW, lobbyRequestInputTxtH := text.Measure(lobbyRequestInputTxt, fontFace, 0)
	lobbyRequestInputTxtOp.GeoM.Translate(screenWidth/2-lobbyRequestInputTxtW/2, screenHeight/2-lobbyRequestInputTxtH)
	text.Draw(screen, lobbyRequestInputTxt, fontFace, &lobbyRequestInputTxtOp)

	if p.txtInputBox == nil {
		p.txtInputBox = graphics.CreateRectangle(txtInputBoxImage, int(lobbyRequestInputTxtW)+10, int(lobbyRequestInputTxtH)+10, int(screenWidth/2-lobbyRequestInputTxtW/2-5), int(screenHeight/2+lobbyRequestInputTxtH-5), 0, 0, 0, 0)
		println(p.txtInputBox.X, p.txtInputBox.Y)
	}
	p.txtInputBox.Draw(screen)

	lobbyRequstInputTxtBoxOp.GeoM.Translate(float64(p.txtInputBox.X+5), float64(p.txtInputBox.Y+5))
	text.Draw(screen, p.LobbyRequestStr, fontFace, &lobbyRequstInputTxtBoxOp)

	p.buttonBack.Draw(screen, op)
	p.buttonJoinLobby.Draw(screen, op)
}
