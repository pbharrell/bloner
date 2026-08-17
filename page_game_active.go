package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pbharrell/bloner-server/connection"
	"github.com/pbharrell/bloner/graphics"
)

type GameActivePage struct {
	game              *Game
	ready             bool
	pickingTrumpReady bool
	state             GameState
	fontFace          *text.GoTextFace
	overlay           graphics.Shape
	previewCardInited bool
	previewCard       Card
	buttonPass        Button
	buttonConfirm     Button
	buttonHearts      Button
	buttonDiamonds    Button
	buttonClubs       Button
	buttonSpades      Button
	textHearts        Text
	textDiamonds      Text
	textClubs         Text
	textSpades        Text
}

func CreateGameActivePage(g *Game) *GameActivePage {
	overlay := *graphics.CreateRectangle(overlayImage, screenWidth, screenHeight, 0, 0, 0, 0, 0, 0)

	fontFace := &text.GoTextFace{
		Source: g.fontSource,
		Size:   24,
	}

	page := &GameActivePage{
		game:              g,
		ready:             false,
		pickingTrumpReady: false,
		state:             CreateGameState(g.id),
		fontFace:          fontFace,
		previewCardInited: false,
		overlay:           overlay,
	}

	buttonConfirm := *CreateButton(page, g, confirmTrump, "assets/confirm_button.png", "assets/confirm_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	confirmWidth := buttonConfirm.sprite.ImageWidth
	confirmX := screenWidth/2 - confirmWidth/2 + 80
	buttonConfirm.SetLoc(confirmX, buttonConfirm.sprite.Y)

	buttonPass := *CreateButton(page, g, passTrump, "assets/pass_button.png", "assets/pass_button.png", 5, 0, screenHeight/2+80, 0)
	passWidth := buttonPass.sprite.ImageWidth
	passX := screenWidth/2 - passWidth/2 - 80
	buttonPass.SetLoc(passX, buttonPass.sprite.Y)

	buttonHearts := *CreateButton(page, g, heartsTrump, "assets/hearts_button.png", "assets/hearts_button_pressed.png", 4, 0, screenHeight/2-140, 0)
	heartsWidth := buttonHearts.sprite.ImageWidth
	heartsX := screenWidth/2 - heartsWidth/2 - 140
	buttonHearts.SetLoc(heartsX, buttonHearts.sprite.Y)

	buttonDiamonds := *CreateButton(page, g, diamondsTrump, "assets/diamonds_button.png", "assets/diamonds_button_pressed.png", 4, 0, screenHeight/2-140, 0)
	diamondsWidth := buttonDiamonds.sprite.ImageWidth
	diamondsX := screenWidth/2 - diamondsWidth/2 + 140
	buttonDiamonds.SetLoc(diamondsX, buttonDiamonds.sprite.Y)

	buttonClubs := *CreateButton(page, g, clubsTrump, "assets/clubs_button.png", "assets/clubs_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	buttonClubs.SetLoc(heartsX, buttonClubs.sprite.Y)

	buttonSpades := *CreateButton(page, g, spadesTrump, "assets/spades_button.png", "assets/spades_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	buttonSpades.SetLoc(diamondsX, buttonSpades.sprite.Y)

	page.buttonPass = buttonPass
	page.buttonConfirm = buttonConfirm
	page.buttonHearts = buttonHearts
	page.buttonDiamonds = buttonDiamonds
	page.buttonClubs = buttonClubs
	page.buttonSpades = buttonSpades
	return page
}

func (p *GameActivePage) HandleHandFinish() {
	toastTxt := ""

	teamBlackTricks := 0
	teamRedTricks := 0
	for i := range p.state.teams {
		for j, player := range p.state.teams[i].players {
			switch i {
			case 0:
				teamBlackTricks += len(player.wonTricks)
			case 1:
				teamRedTricks += len(player.wonTricks)
			}

			p.state.teams[i].players[j].wonTricks = []*Card{}
		}
	}

	if teamBlackTricks > teamRedTricks {
		if teamBlackTricks == 5 {
			p.state.teams[Black].points += 2
			toastTxt += "\nBlack team awarded 2 points!"
		} else {
			p.state.teams[Black].points++
			toastTxt += "\nBlack team awarded 1 point!"
		}
	} else if teamBlackTricks < teamRedTricks {
		if teamRedTricks == 5 {
			p.state.teams[Red].points += 2
			toastTxt += "\nRed team awarded 2 points!"
		} else {
			p.state.teams[Red].points++
			toastTxt += "\nRed team awarded 1 point!"
		}
	}

	p.game.ShowToast(toastTxt, 5*time.Second)

	p.state.DealCards()

	p.previewCardInited = false
	p.state.trick.playCard(p.state.drawPile.drawCard(.1, screenWidth/2+20, 0, 0 /*faceDown */, false))
	p.state.trumpDrawPlayer = (p.state.trumpDrawPlayer + 1) % 4
	p.state.activePlayer += p.state.trumpDrawPlayer

	if p.state.activePlayer == p.state.GetClient().AbsPos {
		p.game.SendStateResponse()
	}
}

func (p *GameActivePage) Update() {
	if !p.state.turnInfo.inited {
		p.state.turnInfo.turnInfo = connection.TurnInfo{
			PlayerId: p.state.id,
		}
		p.state.turnInfo.inited = true

	}

	if p.ready && !p.previewCardInited && len(p.state.trick.Pile) > 0 {
		p.previewCard = *p.state.trick.Pile[len(p.state.trick.Pile)-1]
		p.previewCard.Sprite.ImageScale = .22
		p.previewCard.Sprite.SyncSpriteDimensions()
		p.previewCard.Sprite.X = screenWidth/2 - p.previewCard.Sprite.ImageWidth/2
		p.previewCard.Sprite.Y = screenHeight/2 - p.previewCard.Sprite.ImageHeight/2 - 50
		p.previewCardInited = true
	}

	client := p.state.GetClient()
	if p.state.activePlayer == client.AbsPos {
		p.UpdateClientTurn()
	}

	if len(p.state.trick.Pile) >= 4 {
		highestCard := GetHighestCardFromPile(p.state.trick.Pile, p.state.trick.LeadSuit, *p.state.trumpSuit)
		p.state.GetPlayerById(highestCard.PlayerId).WinTrick(p.state.id)
		p.state.trick.clear()
	}

	outOfCards := true
	for _, team := range p.state.teams {
		for _, player := range team.players {
			if len(player.Cards) > 0 {
				outOfCards = false
			}
		}
	}

	if outOfCards {
		p.HandleHandFinish()
	}
}

func (p *GameActivePage) UpdateClientDiscard() {
	x, y := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Look through sprites in reverse order since a card on the right is on top
		for i, card := range slices.Backward(p.state.GetClient().Cards) {
			if card.Sprite.In(x, y) {
				discarded := card
				p.state.drawPile.discard(discarded)
				if p.state.GetClient().Discard(i, p.state.id) != discarded {
					println("Failed to discard card from hand!! Should not be here.")
				}

				p.SendTurnTrumpDiscard(discarded)
				break
			}
		}
	}
}

func (p *GameActivePage) InitPickingTrumpButtons() {
	type loc struct {
		X       int
		Y       int
		OffsetY int
	}

	trumpLocs := [3]loc{
		{X: screenWidth / 2, Y: screenHeight/2 - 140, OffsetY: -20},
		{X: screenWidth/2 + 140, Y: screenHeight/2 + 80, OffsetY: 85},
		{X: screenWidth/2 - 140, Y: screenHeight/2 + 80, OffsetY: 85},
	}

	type suitText struct {
		Suit    Suit
		SuitStr string
		Button  *Button
		Text    *Text
	}

	suitTexts := [4]suitText{
		{Suit: Hearts, SuitStr: "Hearts", Button: &p.buttonHearts, Text: &p.textHearts},
		{Suit: Diamonds, SuitStr: "Diamonds", Button: &p.buttonDiamonds, Text: &p.textDiamonds},
		{Suit: Clubs, SuitStr: "Clubs", Button: &p.buttonClubs, Text: &p.textClubs},
		{Suit: Spades, SuitStr: "Spades", Button: &p.buttonSpades, Text: &p.textSpades},
	}

	trumpLoc := 0
	for _, suitText := range suitTexts {
		if suitText.Suit == p.previewCard.Suit {
			suitText.Button.sprite.Visible = false
			suitText.Button.pressedSprite.Visible = false
			continue
		}

		suitText.Button.SetLoc(
			trumpLocs[trumpLoc].X-suitText.Button.sprite.ImageWidth/2,
			trumpLocs[trumpLoc].Y-suitText.Button.sprite.ImageHeight/2,
		)

		txtOp := text.DrawOptions{}
		txtW, txtH := text.Measure(suitText.SuitStr, p.fontFace, 0)

		*suitText.Text = *CreateText(p, p.game, p.fontFace, suitText.SuitStr, txtOp, int(suitText.Button.sprite.X+suitText.Button.sprite.ImageWidth/2-int(txtW/2)), int(suitText.Button.sprite.Y-int(txtH/2)+trumpLocs[trumpLoc].OffsetY))

		trumpLoc++
	}

	p.buttonPass.SetLoc(screenWidth/2-p.buttonPass.sprite.ImageWidth/2, screenHeight/2-p.buttonPass.sprite.ImageHeight/2)

	p.pickingTrumpReady = true
}

func (p *GameActivePage) UpdatePickingTrump() {
	x, y := ebiten.CursorPosition()
	mouseButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	if p.state.passCounter < 4 {
		p.buttonConfirm.Update(x, y, mouseButtonPressed)
		p.buttonPass.Update(x, y, mouseButtonPressed)
	} else {
		if !p.pickingTrumpReady {
			p.InitPickingTrumpButtons()
		}
		p.buttonHearts.Update(x, y, mouseButtonPressed)
		p.buttonDiamonds.Update(x, y, mouseButtonPressed)
		p.buttonClubs.Update(x, y, mouseButtonPressed)
		p.buttonSpades.Update(x, y, mouseButtonPressed)
		if p.state.passCounter < 7 {
			p.buttonPass.Update(x, y, mouseButtonPressed)
		}
	}
}

func (p *GameActivePage) UpdateClientPlay() {
	x, y := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		// Look through sprites in reverse order since a card on the right is on top
		for i := len(p.state.GetClient().Cards) - 1; i >= 0; i-- {
			card := p.state.GetClient().Cards[i]
			if card.Sprite.In(x, y) {
				p.state.PlayCard(p.state.id, i)
				p.SendTurnCardPlay(card)
				break
			}
		}
	}
}

func (p *GameActivePage) UpdateClientTurn() {
	if len(p.state.GetClient().Cards) > 5 {
		p.UpdateClientDiscard()

	} else if p.state.IsPickingTrump() {
		p.UpdatePickingTrump()

	} else {
		p.UpdateClientPlay()
	}

	p.state.drawPile.Update()
}

func (p *GameActivePage) Draw(screen *ebiten.Image) {
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.p.state. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkp.state.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	op := ebiten.DrawImageOptions{}

	if !p.state.IsPickingTrump() {
		p.state.GetClient().Draw(screen, op)
	}

	teamScoreText := "Team %v score: %v"

	team1ScoreText := fmt.Sprintf(teamScoreText, 1, p.state.teams[Black].points)
	txtOp := text.DrawOptions{}
	txtW, txtH := text.Measure(team1ScoreText, p.fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2-txtW-30, screenHeight/2-txtH/2-110)
	text.Draw(screen, team1ScoreText, p.fontFace, &txtOp)

	team2ScoreText := fmt.Sprintf(teamScoreText, 2, p.state.teams[Red].points)
	txtOp = text.DrawOptions{}
	txtW, txtH = text.Measure(team2ScoreText, p.fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2+30, screenHeight/2-txtH/2-110)
	text.Draw(screen, team2ScoreText, p.fontFace, &txtOp)

	p.state.drawPile.Draw(screen, op)
	p.state.trick.Draw(screen, op)

	for _, team := range p.state.teams {
		for _, player := range team.players {
			if player.Id != p.state.id {
				// Simply draw the other players (non-client)
				player.Draw(screen, op)
			}
		}
	}

	// We've got some work to do for the client
	if len(p.state.GetClient().Cards) > 5 {
		p.overlay.Draw(screen)

		var (
			discardText = "Click a card to discard"
			txtOp       = text.DrawOptions{}
		)

		txtW, txtH := text.Measure(discardText, p.fontFace, 0)
		txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
		text.Draw(screen, discardText, p.fontFace, &txtOp)
		p.state.GetClient().Draw(screen, op)

	} else if p.state.trumpSuit == nil {
		p.overlay.Draw(screen)

		if p.previewCardInited && p.state.passCounter < 4 {
			p.previewCard.Draw(screen, op)
		}

		// **Everything on top of fade overlay start here**

		if p.state.activePlayer == p.state.GetClient().AbsPos {
			p.state.GetClient().Draw(screen, op)

			if p.state.passCounter < 4 {
				p.buttonConfirm.Draw(screen, op)
				p.buttonPass.Draw(screen, op)
			} else {
				type SuitText struct {
					Suit    string
					OffsetX float64
					OffsetY float64
				}

				p.buttonHearts.Draw(screen, op)
				if p.buttonHearts.sprite.Visible {
					p.textHearts.Draw(screen, op)
				}
				p.buttonDiamonds.Draw(screen, op)
				if p.buttonDiamonds.sprite.Visible {
					p.textDiamonds.Draw(screen, op)
				}
				p.buttonClubs.Draw(screen, op)
				if p.buttonClubs.sprite.Visible {
					p.textClubs.Draw(screen, op)
				}
				p.buttonSpades.Draw(screen, op)
				if p.buttonSpades.sprite.Visible {
					p.textSpades.Draw(screen, op)
				}

				if p.state.passCounter < 7 {
					p.buttonPass.Draw(screen, op)
				}
			}

		} else {
			var (
				waitingText = fmt.Sprintf("Waiting on player %v to choose...", p.state.GetPlayerByAbsPos(p.state.activePlayer).Id)
				txtOp       = text.DrawOptions{}
			)

			txtW, txtH := text.Measure(waitingText, p.fontFace, 0)
			txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
			text.Draw(screen, waitingText, p.fontFace, &txtOp)
		}
	} else if p.state.trumpSuit != nil && p.state.activePlayer != p.state.GetClient().AbsPos {
		p.overlay.Draw(screen)

		// **Everything on top of fade overlay start here**
		var (
			waitingText = fmt.Sprintf("Player %v's turn...", p.state.activePlayer)
			txtOp       = text.DrawOptions{}
		)

		txtW, txtH := text.Measure(waitingText, p.fontFace, 0)
		txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
		text.Draw(screen, waitingText, p.fontFace, &txtOp)
	}

}

func (p *GameActivePage) SendTurnInfo() {
	if p.game.server.connected && p.state.turnInfo.inited {
		println("Sending turn info for player id:", p.state.turnInfo.turnInfo.PlayerId)
		p.game.server.server.Send(connection.Message{
			Type: "turn_info",
			Data: p.state.turnInfo.turnInfo,
		})
	} else {
		println("turn_info not sent since no server is connected.")
	}
	p.state.turnInfo.inited = false
}

func (p *GameActivePage) SendTurnCardPlay(card *Card) {
	p.state.turnInfo.turnInfo.TurnType = connection.CardPlay
	p.state.turnInfo.turnInfo.CardPlay = card.Encode()
	p.SendTurnInfo()
	p.state.EndTurn()
}

func (p *GameActivePage) SendTurnTrumpDiscard(card *Card) {
	p.state.turnInfo.turnInfo.TurnType = connection.TrumpDiscard
	p.state.turnInfo.turnInfo.TrumpDiscard = card.Encode()
	p.SendTurnInfo()
	p.state.EndTurn()
}

func (p *GameActivePage) SendTurnTrumpPick(suit int8) {
	p.state.turnInfo.turnInfo.TurnType = connection.TrumpPick
	p.state.turnInfo.turnInfo.TrumpPick = suit
	p.SendTurnInfo()
}

func (p *GameActivePage) SendTurnTrumpPass() {
	p.state.turnInfo.turnInfo.TurnType = connection.TrumpPass
	p.SendTurnInfo()
	p.state.EndTurn()
}
