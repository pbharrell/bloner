package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pbharrell/bloner/graphics"
)

type GameActivePage struct {
	game           *Game
	fontFace       *text.GoTextFace
	overlay        graphics.Shape
	buttonConfirm  Button
	buttonCancel   Button
	buttonPass     Button
	buttonHearts   Button
	buttonDiamonds Button
	buttonClubs    Button
	buttonSpades   Button
}

func CreateGameActivePage(g *Game) *GameActivePage {
	overlay := *graphics.CreateRectangle(overlayImage, screenWidth, screenHeight, 0, 0, 0, 0, 0, 0)

	fontFace := &text.GoTextFace{
		Source: g.fontSource,
		Size:   24,
	}

	p := &GameActivePage{
		game:     g,
		fontFace: fontFace,
		overlay:  overlay,
	}

	p.buttonConfirm = *CreateButton(p, g, confirmTrump, "assets/confirm_button.png", "assets/confirm_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	confirmWidth := p.buttonConfirm.sprite.ImageWidth
	confirmX := screenWidth/2 - confirmWidth/2 + 80
	p.buttonConfirm.SetLoc(confirmX, p.buttonConfirm.sprite.Y)

	p.buttonCancel = *CreateButton(p, g, passTrump, "assets/cancel_button.png", "assets/cancel_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	cancelWidth := p.buttonCancel.sprite.ImageWidth
	cancelX := screenWidth/2 - cancelWidth/2 - 80
	p.buttonCancel.SetLoc(cancelX, p.buttonCancel.sprite.Y)

	p.buttonPass = *CreateButton(p, g, passTrump, "assets/pass_button.png", "assets/pass_button.png", 5, 0, 0, 0)
	passWidth := p.buttonPass.sprite.ImageWidth
	passHeight := p.buttonPass.sprite.ImageHeight
	passX := screenWidth/2 - passWidth/2
	passY := screenHeight/2 - passHeight/2
	p.buttonPass.SetLoc(passX, passY)

	p.buttonHearts = *CreateButton(p, g, heartsTrump, "assets/hearts_button.png", "assets/hearts_button_pressed.png", 4, 0, screenHeight/2-140, 0)
	heartsWidth := p.buttonHearts.sprite.ImageWidth
	heartsX := screenWidth/2 - heartsWidth/2 - 140
	p.buttonHearts.SetLoc(heartsX, p.buttonHearts.sprite.Y)

	p.buttonDiamonds = *CreateButton(p, g, diamondsTrump, "assets/diamonds_button.png", "assets/diamonds_button_pressed.png", 4, 0, screenHeight/2-140, 0)
	diamondsWidth := p.buttonDiamonds.sprite.ImageWidth
	diamondsX := screenWidth/2 - diamondsWidth/2 + 140
	p.buttonDiamonds.SetLoc(diamondsX, p.buttonDiamonds.sprite.Y)

	p.buttonClubs = *CreateButton(p, g, clubsTrump, "assets/clubs_button.png", "assets/clubs_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	p.buttonClubs.SetLoc(heartsX, p.buttonClubs.sprite.Y)

	p.buttonSpades = *CreateButton(p, g, spadesTrump, "assets/spades_button.png", "assets/spades_button_pressed.png", 4, 0, screenHeight/2+80, 0)
	p.buttonSpades.SetLoc(diamondsX, p.buttonSpades.sprite.Y)

	return p
}

func (p *GameActivePage) Update() {
	client := p.game.GetClient()
	if p.game.activePlayer == client.AbsPos {
		p.UpdateClientTurn()
	}

	if len(p.game.trick.Pile) >= 4 {
		highestCard := GetHighestCardFromPile(p.game.trick.Pile, p.game.trick.LeadSuit, *p.game.trumpSuit)
		println("Highest card returned from pile:", SuitToString(highestCard.Suit), NumberToString(highestCard.Number))
		println("Highest card player id:", highestCard.PlayerId)
		p.game.GetPlayerById(highestCard.PlayerId).WinTrick(p.game.id)
		p.game.trick.clear()
	}

	outOfCards := true
	for _, team := range p.game.teams {
		for _, player := range team.players {
			if len(player.Cards) > 0 {
				outOfCards = false
			}
		}
	}

	if outOfCards {
		teamBlackTricks := 0
		teamRedTricks := 0
		for i := range p.game.teams {
			for j, player := range p.game.teams[i].players {
				if i == 0 {
					teamBlackTricks += len(player.wonTricks)
				} else if i == 1 {
					teamRedTricks += len(player.wonTricks)
				}

				p.game.teams[i].players[j].wonTricks = []*Card{}
			}
		}

		if teamBlackTricks > teamRedTricks {
			if teamBlackTricks == 5 {
				p.game.teams[Black].points += 2
			} else {
				p.game.teams[Black].points++
			}
		} else if teamBlackTricks < teamRedTricks {
			if teamRedTricks == 5 {
				p.game.teams[Red].points += 2
			} else {
				p.game.teams[Red].points++
			}
		}

		p.game.DealCards()

		p.game.trick.playCard(p.game.drawPile.drawCard(.1, screenWidth/2+20, 0, 0 /*faceDown */, false))
		p.game.trumpDrawPlayer = (p.game.trumpDrawPlayer + 1) % 4
		p.game.activePlayer += p.game.trumpDrawPlayer

		if p.game.activePlayer == client.AbsPos {
			p.game.SendStateResponse()
		}
	}
}

func (p *GameActivePage) Draw(screen *ebiten.Image) {
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.p.game. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkp.game.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	op := ebiten.DrawImageOptions{}

	if !p.game.IsPickingTrump() {
		p.game.GetClient().Draw(screen, op)
	}

	teamScoreText := "Team %v score: %v"

	team1ScoreText := fmt.Sprintf(teamScoreText, 1, p.game.teams[Black].points)
	txtOp := text.DrawOptions{}
	txtW, txtH := text.Measure(team1ScoreText, p.fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2-txtW-30, screenHeight/2-txtH/2-110)
	text.Draw(screen, team1ScoreText, p.fontFace, &txtOp)

	team2ScoreText := fmt.Sprintf(teamScoreText, 2, p.game.teams[Red].points)
	txtOp = text.DrawOptions{}
	txtW, txtH = text.Measure(team2ScoreText, p.fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2+30, screenHeight/2-txtH/2-110)
	text.Draw(screen, team2ScoreText, p.fontFace, &txtOp)

	p.game.drawPile.Draw(screen, op)
	p.game.trick.Draw(screen, op)

	for _, team := range p.game.teams {
		for _, player := range team.players {
			if player.Id != p.game.id {
				// Simply draw the other players (non-client)
				player.Draw(screen, op)
			}
		}
	}

	// We've got some work to do for the client
	if len(p.game.GetClient().Cards) > 5 {
		p.overlay.Draw(screen)

		var (
			discardText = "Click a card to discard"
			txtOp       = text.DrawOptions{}
		)

		// Create font faces with different sizes as needed
		fontFace := &text.GoTextFace{
			Source: p.game.fontSource,
			Size:   24,
		}

		txtW, txtH := text.Measure(discardText, fontFace, 0)
		txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
		text.Draw(screen, discardText, fontFace, &txtOp)
		p.game.GetClient().Draw(screen, op)

	} else if p.game.trumpSuit == nil {
		p.overlay.Draw(screen)

		// **Everything on top of fade overlay start here**

		if p.game.activePlayer == p.game.GetClient().AbsPos {
			p.game.GetClient().Draw(screen, op)

			if p.game.passCounter < 4 {
				p.buttonConfirm.Draw(screen, op)
				p.buttonCancel.Draw(screen, op)
			} else {
				// Create font faces with different sizes as needed
				fontFace := &text.GoTextFace{
					Source: p.game.fontSource,
					Size:   24,
				}

				type SuitText struct {
					Suit    string
					OffsetX float64
					OffsetY float64
				}
				suitTexts := []SuitText{
					{Suit: "Hearts", OffsetX: -140, OffsetY: -50},
					{Suit: "Diamonds", OffsetX: +140, OffsetY: -50},
					{Suit: "Clubs", OffsetX: -140, OffsetY: +60},
					{Suit: "Spades", OffsetX: +140, OffsetY: +60},
				}

				for _, suitText := range suitTexts {
					txtOp := text.DrawOptions{}
					txtW, txtH := text.Measure(suitText.Suit, fontFace, 0)
					centeredX, centeredY := screenWidth/2-txtW/2, screenHeight/2-txtH/2
					txtOp.GeoM.Translate(centeredX+suitText.OffsetX, centeredY+suitText.OffsetY)
					text.Draw(screen, suitText.Suit, fontFace, &txtOp)
				}

				p.buttonHearts.Draw(screen, op)
				p.buttonDiamonds.Draw(screen, op)
				p.buttonClubs.Draw(screen, op)
				p.buttonSpades.Draw(screen, op)

				if p.game.passCounter < 7 {
					p.buttonPass.Draw(screen, op)
				}
			}

		} else {
			var (
				waitingText = fmt.Sprintf("Waiting on player %v to choose...", p.game.GetActivePlayer().Id)
				txtOp       = text.DrawOptions{}
			)

			// Create font faces with different sizes as needed
			fontFace := &text.GoTextFace{
				Source: p.game.fontSource,
				Size:   24,
			}

			txtW, txtH := text.Measure(waitingText, fontFace, 0)
			txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
			text.Draw(screen, waitingText, fontFace, &txtOp)
		}
	} else if p.game.trumpSuit != nil && p.game.activePlayer != p.game.GetClient().AbsPos {
		p.overlay.Draw(screen)

		// **Everything on top of fade overlay start here**
		var (
			waitingText = fmt.Sprintf("Player %v's turn...", p.game.activePlayer)
			txtOp       = text.DrawOptions{}
		)

		// Create font faces with different sizes as needed

		txtW, txtH := text.Measure(waitingText, p.fontFace, 0)
		txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
		text.Draw(screen, waitingText, p.fontFace, &txtOp)
	}
}

func (p *GameActivePage) UpdateClientTurn() {
	client := p.game.GetClient()
	if len(client.Cards) > 5 {
		x, y := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Look through sprites in reverse order since a card on the right is on top
			for i := len(client.Cards) - 1; i >= 0; i-- {
				card := client.Cards[i]
				if card.Sprite.In(x, y) {
					discarded := client.Cards[i]
					p.game.drawPile.discard(discarded)
					if client.Discard(i, p.game.id) != discarded {
						println("Failed to discard card from hand!! Should not be here.")
					}

					p.game.SendTurnTrumpDiscard(discarded)
					break
				}
			}
		}

	} else if p.game.IsPickingTrump() {
		x, y := ebiten.CursorPosition()
		mouseButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

		if p.game.passCounter < 4 {
			p.buttonConfirm.Update(x, y, mouseButtonPressed)
			p.buttonCancel.Update(x, y, mouseButtonPressed)
		} else {
			p.buttonHearts.Update(x, y, mouseButtonPressed)
			p.buttonDiamonds.Update(x, y, mouseButtonPressed)
			p.buttonClubs.Update(x, y, mouseButtonPressed)
			p.buttonSpades.Update(x, y, mouseButtonPressed)
			if p.game.passCounter < 7 {
				p.buttonPass.Update(x, y, mouseButtonPressed)
			}
		}

	} else {
		x, y := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

			// Look through sprites in reverse order since a card on the right is on top
			for i := len(client.Cards) - 1; i >= 0; i-- {
				card := client.Cards[i]
				if card.Sprite.In(x, y) {
					p.game.PlayCard(p.game.id, i)
					p.game.SendTurnCardPlay(card)
					break
				}
			}

			// Only want to add a card to hand from draw pile if debugging
			if false || p.game.debug {
				if p.game.drawPile.Sprite.In(x, y) && len(client.Cards) < 5 {
					card := p.game.drawPile.drawCard(.1, 0, 0, 0 /* faceDown */, false)
					if card != nil {
						card.PlayerId = client.Id
						client.Cards = append(client.Cards, card)
						client.ArrangeHand(client.Id)
					}
				}
			}
		}
	}

	p.game.drawPile.Update()
}
