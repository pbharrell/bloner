package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"net"
	"slices"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/pbharrell/bloner-server/connection"

	"github.com/pbharrell/bloner/graphics"
)

var (
	overlayImage *ebiten.Image

	buttonConfirmImage        *ebiten.Image
	buttonConfirmAlpha        *image.Alpha
	buttonPressedConfirmImage *ebiten.Image
	buttonPressedConfirmAlpha *image.Alpha

	buttonCancelImage        *ebiten.Image
	buttonCancelAlpha        *image.Alpha
	buttonPressedCancelImage *ebiten.Image
	buttonPressedCancelAlpha *image.Alpha

	ebitenImage *ebiten.Image
	cardImage   *ebiten.Image
)

func init() {
	overlayImage = ebiten.NewImage(3, 3)
	overlayImage.Fill(color.RGBA{0, 0, 0, 200})

	buttonConfirmImage, buttonConfirmAlpha = graphics.LoadImageFromFile("./assets/confirm_button.png")
	buttonPressedConfirmImage, buttonPressedConfirmAlpha = graphics.LoadImageFromFile("./assets/confirm_button_pressed.png")

	buttonCancelImage, buttonCancelAlpha = graphics.LoadImageFromFile("./assets/cancel_button.png")
	buttonPressedCancelImage, buttonPressedCancelAlpha = graphics.LoadImageFromFile("./assets/cancel_button_pressed.png")

	ebitenImage = graphics.LoadImage(&images.Ebiten_png)

	// Read the file into a byte array
	var err error
	cardImage, _, err = ebitenutil.NewImageFromFile("./assets/ace_of_spades.png")
	if err != nil {
		log.Fatal(err)
	}

	initCardImages()
}

const (
	screenWidth  = 720
	screenHeight = 540
	maxAngle     = 360
)

type mode uint8

const (
	LobbyWait mode = iota
	LobbyAssigned
	GameActive
)

type turnType uint8

const (
	TrumpChoice turnType = iota
	CardPlay
)

type Server struct {
	connected bool
	server    connection.Server
}

type TurnInfo struct {
	inited   bool
	turnInfo connection.TurnInfo
}

type Game struct {
	inited        bool
	debug         bool
	server        Server
	id            int
	mode          mode
	lobbyId       int
	fontSource    *text.GoTextFaceSource
	trumpSuit     *Suit
	teams         [2]Team
	activePlayer  int
	touchIDs      []ebiten.TouchID
	buttonConfirm Button
	buttonCancel  Button
	overlay       graphics.Shape
	drawPile      DrawPile
	trick         Trick
	turnInfo      TurnInfo
}

func (g *Game) initOverlay() {
	var vertices []ebiten.Vertex
	vertices = append(vertices, ebiten.Vertex{
		DstX:   0,
		DstY:   0,
		SrcX:   float32(0),
		SrcY:   float32(0),
		ColorR: float32(1),
		ColorG: float32(1),
		ColorB: float32(1),
		ColorA: 1,
	})
	vertices = append(vertices, ebiten.Vertex{
		DstX:   screenWidth,
		DstY:   0,
		SrcX:   float32(1),
		SrcY:   float32(0),
		ColorR: float32(1),
		ColorG: float32(1),
		ColorB: float32(1),
		ColorA: 1,
	})
	vertices = append(vertices, ebiten.Vertex{
		DstX:   screenWidth,
		DstY:   screenHeight,
		SrcX:   float32(1),
		SrcY:   float32(1),
		ColorR: float32(1),
		ColorG: float32(1),
		ColorB: float32(1),
		ColorA: 1,
	})
	vertices = append(vertices, ebiten.Vertex{
		DstX:   0,
		DstY:   screenHeight,
		SrcX:   float32(0),
		SrcY:   float32(1),
		ColorR: float32(1),
		ColorG: float32(1),
		ColorB: float32(1),
		ColorA: 1,
	})

	indices := graphics.GenIndices(len(vertices))
	g.overlay = *graphics.CreateShape(overlayImage, vertices, indices, 0, 0, 0, 0, 0, 0)
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	// FIXME: Remove debug when appropriate
	g.debug = true

	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	g.mode = LobbyWait

	g.lobbyId = -1

	g.fontSource = fontSource

	g.drawPile.Sprite = *graphics.CreateSprite(blankCardImage, blankCardAlphaImage, .35, screenWidth/2, screenHeight/2, 0, 0, 0, 0)
	g.drawPile.Sprite.X = screenWidth/2 - g.drawPile.Sprite.ImageWidth - 20
	g.drawPile.Sprite.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.drawPile.shuffleDrawPile()

	g.teams[Black].teamColor = Black
	g.teams[Red].teamColor = Red
	g.teams[Black].tricksWon = 0
	g.teams[Red].tricksWon = 0
	g.teams[Black].players[0] = CreatePlayer(0, Black, 5, Bottom, .35, &g.drawPile /* faceDown */, false)
	g.teams[Red].players[0] = CreatePlayer(1, Red, 5, Left, .35, &g.drawPile /* faceDown */, false)
	g.teams[Black].players[1] = CreatePlayer(2, Black, 5, Top, .35, &g.drawPile /* faceDown */, false)
	g.teams[Red].players[1] = CreatePlayer(3, Red, 5, Right, .35, &g.drawPile /* faceDown */, false)

	g.trick.X = screenWidth/2 + 20
	g.trick.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.trick.playCard(g.drawPile.drawCard(.35, screenWidth/2+20, 0, 0 /*faceDown */, false))

	g.activePlayer = 0

	g.buttonConfirm = *CreateButton(g, confirmTrump, buttonConfirmImage, buttonConfirmAlpha, buttonPressedConfirmImage, buttonPressedConfirmAlpha, 4, 0, screenHeight/2+80, 0)
	confirmWidth := g.buttonConfirm.sprite.ImageWidth
	confirmX := screenWidth/2 - confirmWidth/2 + 80
	g.buttonConfirm.SetLoc(confirmX, g.buttonConfirm.sprite.Y)

	g.buttonCancel = *CreateButton(g, cancelTrump, buttonCancelImage, buttonCancelAlpha, buttonPressedCancelImage, buttonPressedCancelAlpha, 4, 0, screenHeight/2+80, 0)
	cancelWidth := g.buttonCancel.sprite.ImageWidth
	cancelX := screenWidth/2 - cancelWidth/2 - 80
	g.buttonCancel.SetLoc(cancelX, g.buttonCancel.sprite.Y)

	g.initOverlay()

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Printf("Error connecting to server: `%v`\nCan debug in offline mode, but don't expect to join a game anytime soon.", err)
		return
	}

	g.server.server = connection.Server{
		Conn: conn,
		Data: make(chan connection.Message),
	}
	g.server.connected = true

	go g.server.server.Listen()
}

func (g *Game) debugPrintln(msg string) {
	if g.debug {
		println(msg)
	}
}

func (g *Game) GetPlayer(id int) *Player {
	for i := range g.teams {
		for j := range g.teams[i].players {
			if g.id == g.teams[i].players[j].Id {
				return &g.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetPlayer(%v)` with no player matching that id present in player list", id)
	return nil
}

func (g *Game) GetClient() *Player {
	return g.GetPlayer(g.id)
}

func (g *Game) SetActiveId(id int) {
	g.activePlayer = id
}

func (g *Game) SetActivePlayer(player *Player) {
	g.activePlayer = player.Id
}

func DecodeCardPile(encPile []connection.Card, scale float64, faceDown bool) []*Card {
	pile := make([]*Card, len(encPile))
	for i, encCard := range encPile {
		pile[i] = (CreateCard(Suit(encCard.Suit), Number(encCard.Number), scale, 0, 0, 0, faceDown))
	}

	return pile
}

func (g *Game) ArrangeTeams() {
	client := g.GetClient()
	g.teams[0].Arrange(client.Id, client.AbsPos)
	g.teams[1].Arrange(client.Id, client.AbsPos)
	println("consider me: arranged")
}

func (g *Game) SendTurnInfo() {
	if g.server.connected {
		g.server.server.Send(connection.Message{
			Type: "turn_info",
			Data: g.turnInfo.turnInfo,
		})
		g.turnInfo.inited = false
	} else {
		println("turn_info not sent since no server is connected.")
	}
}

func (g *Game) EndTurn() {
	g.activePlayer = (g.activePlayer + 1) % 4
}

func (g *Game) IsPickingTrump() bool {
	return g.activePlayer == 0 && g.trumpSuit == nil
}

func (g *Game) EncodeGameState() connection.GameState {
	intTrumpSuit := -1
	if g.trumpSuit != nil {
		intTrumpSuit = int(*g.trumpSuit)
	}

	encDrawPile := make([]connection.Card, len(g.drawPile.Pile))
	for i, cardInt := range g.drawPile.Pile {
		encDrawPile[i] = CreateCard(Suit(cardInt/6), Number(cardInt%6), 0, 0, 0, 0 /*faceDown*/, true).Encode()
	}

	encPlayPile := g.trick.Encode()

	teamState := [2]connection.TeamState{
		g.teams[Black].Encode(),
		g.teams[Red].Encode(),
	}

	return connection.GameState{
		PlayerId:     g.id,
		ActivePlayer: g.activePlayer,
		TrumpSuit:    intTrumpSuit,
		DrawPile:     encDrawPile,
		PlayPile:     encPlayPile,
		TeamState:    teamState,
	}
}

func (g *Game) DecodeGameState(state connection.GameState) {
	g.SetActiveId(state.ActivePlayer)

	if state.TrumpSuit < 0 {
		g.trumpSuit = nil
	} else {
		*g.trumpSuit = Suit(state.TrumpSuit)
	}

	g.drawPile.Pile = make([]int, len(state.DrawPile))
	for i, card := range state.DrawPile {
		g.drawPile.Pile[i] = int(card.Suit)*6 + int(card.Number)
	}

	g.trick.Decode(state.PlayPile)

	g.teams[Black].Decode(Black, state.TeamState[Black])
	g.teams[Red].Decode(Red, state.TeamState[Red])

	g.ArrangeTeams()
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	if !g.turnInfo.inited {
		g.turnInfo.turnInfo = connection.TurnInfo{
			PlayerId: g.GetClient().Id,
		}
		g.turnInfo.inited = true
	}

	select {
	case msg := <-g.server.server.Data:
		g.HandleMessage(msg)
		break
	default:
		break
	}

	if g.debug {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.mode != GameActive {
			g.mode++
		}
	}

	switch g.mode {
	case LobbyWait:
		break
	case LobbyAssigned:
		break
	case GameActive:
		g.UpdateGameActive()
		break
	}

	return nil
}

func (g *Game) UpdateGameActive() {
	client := g.GetClient()
	if g.activePlayer == client.Id {
		g.UpdateClientTurn()
	}
}

func (g *Game) UpdateClientTurn() {
	client := g.GetClient()
	if len(client.Cards) > 5 {
		x, y := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Look through sprites in reverse order since a card on the right is on top
			for i := len(client.Cards) - 1; i >= 0; i-- {
				card := client.Cards[i]
				if card.Sprite.In(x, y) {
					// TODO: Update sprite here to be blank side
					discarded := client.Cards[i]
					g.drawPile.discard(discarded)
					client.Cards = slices.Delete(client.Cards, i, i+1)
					client.ArrangeHand(client.Id)

					g.turnInfo.turnInfo.TurnType = connection.TrumpDiscard
					g.turnInfo.turnInfo.TrumpDiscard = discarded.Encode()
					g.SendTurnInfo()
					g.EndTurn()
					break
				}
			}
		}

	} else if g.IsPickingTrump() {
		x, y := ebiten.CursorPosition()
		mouseButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

		g.buttonConfirm.Update(x, y, mouseButtonPressed)
		g.buttonCancel.Update(x, y, mouseButtonPressed)

	} else {
		x, y := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

			// Look through sprites in reverse order since a card on the right is on top
			for i := len(client.Cards) - 1; i >= 0; i-- {
				card := client.Cards[i]
				if card.Sprite.In(x, y) {
					g.trick.playCard(client.Cards[i])
					client.Cards = slices.Delete(client.Cards, i, i+1)
					client.ArrangeHand(client.Id)
					break
				}
			}

			// Only want to add a card to hand from draw pile if debugging
			if g.debug {
				if g.drawPile.Sprite.In(x, y) && len(client.Cards) < 5 {
					card := g.drawPile.drawCard(.35, 0, 0, 0 /* faceDown */, false)
					if card != nil {
						client.Cards = append(client.Cards, card)
						client.ArrangeHand(client.Id)
					}
				}
			}
		}
	}

	g.drawPile.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{161, 191, 123, 1})

	switch g.mode {
	case LobbyWait:
		g.DrawLobbyWait(screen)
		break
	case LobbyAssigned:
		g.DrawLobbyAssigned(screen)
		break
	case GameActive:
		g.DrawGameActive(screen)
		break
	}
}

func (g *Game) DrawLobbyWait(screen *ebiten.Image) {
	g.overlay.Draw(screen)

	var (
		discardText = "Searching for a lobby..."
		txtOp       = text.DrawOptions{}
	)

	// Create font faces with different sizes as needed
	fontFace := &text.GoTextFace{
		Source: g.fontSource,
		Size:   32,
	}

	txtW, txtH := text.Measure(discardText, fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2)
	text.Draw(screen, discardText, fontFace, &txtOp)
}

func (g *Game) DrawLobbyAssigned(screen *ebiten.Image) {
	g.overlay.Draw(screen)

	var (
		discardText = fmt.Sprintf("Lobby found with id: %v!", g.lobbyId)
		txtOp       = text.DrawOptions{}
	)

	// Create font faces with different sizes as needed
	fontFace := &text.GoTextFace{
		Source: g.fontSource,
		Size:   32,
	}

	txtW, txtH := text.Measure(discardText, fontFace, 0)
	txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2)
	text.Draw(screen, discardText, fontFace, &txtOp)
}

func (g *Game) DrawGameActive(screen *ebiten.Image) {
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	op := ebiten.DrawImageOptions{}

	if !g.IsPickingTrump() {
		g.GetClient().Draw(screen, op)
	}

	g.drawPile.Draw(screen, op)
	g.trick.Draw(screen, op)

	for _, team := range g.teams {
		for _, player := range team.players {
			if player.Id != g.GetClient().Id {
				// Simply draw the other players (non-client)
				player.Draw(screen, op)
			}
		}
	}

	// We've got some work to do for the client
	if len(g.GetClient().Cards) > 5 {
		g.overlay.Draw(screen)

		var (
			discardText = "Click a card to discard"
			txtOp       = text.DrawOptions{}
		)

		// Create font faces with different sizes as needed
		fontFace := &text.GoTextFace{
			Source: g.fontSource,
			Size:   24,
		}

		txtW, txtH := text.Measure(discardText, fontFace, 0)
		txtOp.GeoM.Translate(screenWidth/2-txtW/2, screenHeight/2-txtH/2+110)
		text.Draw(screen, discardText, fontFace, &txtOp)
		g.GetClient().Draw(screen, op)

	} else if g.activePlayer == 0 && g.trumpSuit == nil {
		g.overlay.Draw(screen)

		// **Everything on top of fade overlay start here**

		g.buttonConfirm.Draw(screen, op)
		g.buttonCancel.Draw(screen, op)
		g.GetClient().Draw(screen, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func GetEncodedCard(c *Card) connection.Card {
	return connection.Card{
		Suit:   uint8(c.Suit),
		Number: uint8(c.Number),
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("bloner")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// gob.Register(state.GameState{})

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
