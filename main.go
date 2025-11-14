package main

import (
	"bytes"
	"encoding/json"
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

type phase uint8

const (
	LobbyWait phase = iota
	GameStart
)

type Game struct {
	inited        bool
	debug         bool
	server        connection.Server
	id            int
	phase         phase
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

	g.fontSource = fontSource

	g.drawPile.Sprite = *graphics.CreateSpriteFromFile("./assets/ace_of_spades.png", .35, screenWidth/2, screenHeight/2, 0, 0, 0, 0)
	g.drawPile.Sprite.X = screenWidth/2 - g.drawPile.Sprite.ImageWidth - 20
	g.drawPile.Sprite.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.drawPile.shuffleDrawPile()

	g.teams[uint8(Black)].teamColor = Black
	g.teams[uint8(Red)].teamColor = Red
	g.teams[uint8(Black)].tricksWon = 0
	g.teams[uint8(Red)].tricksWon = 0
	g.teams[uint8(Black)].players[0] = CreatePlayer(0, Black, 5, Bottom, .35, &g.drawPile)
	g.teams[uint8(Red)].players[0] = CreatePlayer(1, Red, 5, Left, .35, &g.drawPile)
	g.teams[uint8(Black)].players[1] = CreatePlayer(2, Black, 5, Top, .35, &g.drawPile)
	g.teams[uint8(Red)].players[1] = CreatePlayer(3, Red, 5, Right, .35, &g.drawPile)

	g.trick.X = screenWidth/2 + 20
	g.trick.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.trick.playCard(g.drawPile.drawCard(.35, screenWidth/2+20, 0, 0))

	// g.activePlayer = LeftOpp
	g.activePlayer = 0 // FIXME:

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

	g.server = connection.Server{
		Conn: conn,
		Data: make(chan connection.Message),
	}

	go g.server.Listen()
}

func (g *Game) debugPrintln(msg string) {
	if g.debug {
		println(msg)
	}
}

func (g *Game) GetClient() *Player {
	for i := range g.teams {
		for j := range g.teams[i].players {
			if g.id == g.teams[i].players[j].Id {
				return &g.teams[i].players[j]
			}
		}
	}

	println("ERROR: Should not be here! Called `game.GetClient()` with no client present in player list")
	return nil
}

func DecodeCardPile(encPile []connection.Card, scale float64) []*Card {
	pile := make([]*Card, len(encPile))
	for i, encCard := range encPile {
		pile[i] = (CreateCard(Suit(encCard.Suit), Number(encCard.Number), scale, 0, 0, 0))
	}

	return pile
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
		encDrawPile[i] = CreateCard(Suit(cardInt/6), Number(cardInt%6), 0, 0, 0, 0).Encode()
	}

	encPlayPile := g.trick.Encode()

	teamState := [2]connection.TeamState{
		g.teams[uint8(Black)].Encode(),
		g.teams[uint8(Red)].Encode(),
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
	g.activePlayer = state.ActivePlayer

	if state.TrumpSuit < 0 {
		g.trumpSuit = nil
	} else {
		*g.trumpSuit = Suit(state.TrumpSuit)
	}

	g.drawPile.Pile = make([]int, len(state.DrawPile))
	for i, card := range state.DrawPile {
		g.drawPile.Pile[i] = int(card.Suit)*6 + int(card.Number)
	}

	g.trick.Pile = DecodeCardPile(state.PlayPile, g.trick.Sprite.ImageScale)

	g.teams[uint8(Black)].Decode(state.TeamState[uint8(Black)])
	g.teams[uint8(Red)].Decode(state.TeamState[uint8(Red)])
}

func (g *Game) HandleLobbyAssignMessage(data connection.LobbyAssign) {
	println("Player with id:", data.PlayerId)
	println("Lobby with id:", data.LobbyId)

	g.debugPrintln("Handled lobby assign message!")
}

func (g *Game) HandleStateRequestMessage() {
	gameState := g.EncodeGameState()
	fmt.Printf("%v", gameState)
	g.server.Send(connection.Message{
		Type: "state_res",
		Data: gameState,
	})

	g.debugPrintln("Handled state request message!")
}

func (g *Game) HandleStateResponseMessage(data connection.GameState) {
	g.DecodeGameState(data)
}

func (g *Game) HandleMessage(msg connection.Message) {
	// Marshal Data back into JSON bytes
	raw, err := json.Marshal(msg.Data)
	if err != nil {
		println("marshal error:", err)
		return
	}

	switch msg.Type {
	case "lobby_assign":
		var lobbyAssign connection.LobbyAssign
		if err := json.Unmarshal(raw, &lobbyAssign); err != nil {
			println("LobbyAssign unmarshal error:", err)
			return
		}

		g.HandleLobbyAssignMessage(lobbyAssign)
		break
	case "state_req":
		g.HandleStateRequestMessage()
		break

	case "state_res":
		var stateResponse connection.StateResponse
		if err := json.Unmarshal(raw, &stateResponse); err != nil {
			println("LobbyAssign unmarshal error:", err)
			return
		}

		g.HandleStateResponseMessage(stateResponse)
		break

	default:
		println("Message with unexpected type encountered:", msg.Type)
		return
	}

	g.debugPrintln(fmt.Sprintf("Type:", msg.Type))
	g.debugPrintln(fmt.Sprintf("Data:", msg.Data))

}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	select {
	case msg := <-g.server.Data:
		g.HandleMessage(msg)
		break
	default:
		break
	}

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
					client.ArrangeHand()
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
					client.ArrangeHand()
					break
				}
			}

			if g.drawPile.Sprite.In(x, y) && len(client.Cards) < 5 {
				card := g.drawPile.drawCard(.35, 0, 0, 0)
				if card != nil {
					client.Cards = append(client.Cards, card)
					client.ArrangeHand()
				}
			}
		}
	}

	g.drawPile.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{161, 191, 123, 1})

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
