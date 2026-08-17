package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/coder/websocket"

	"github.com/pbharrell/bloner-server/connection"

	"github.com/pbharrell/bloner/graphics"
)

var (
	overlayImage     *ebiten.Image
	txtInputBoxImage *ebiten.Image
)

//go:embed assets/*
var content embed.FS

func init() {
	overlayImage = ebiten.NewImage(3, 3)
	overlayImage.Fill(color.RGBA{0, 0, 0, 200})

	txtInputBoxImage = ebiten.NewImage(3, 3)
	txtInputBoxImage.Fill(color.RGBA{0, 0, 0, 255})

	initCardImages()
}

const (
	screenWidth  = 720
	screenHeight = 540
	maxAngle     = 360
)

type mode uint8

const (
	LobbyTypeSelect mode = iota
	LobbyRequestInput
	LobbyRequested
	LobbyAssigned
	GameActive
)

type turnType uint8

const (
	TrumpChoice turnType = iota
	CardPlay
)

type server struct {
	server    connection.Server
	connected bool
	lobbyId   int
}

type Page interface {
	Update()
	Draw(screen *ebiten.Image)
}

type TurnInfo struct {
	inited   bool
	turnInfo connection.TurnInfo
}

type GameState struct {
	id              int
	activePlayer    PlayPos
	trumpDrawPlayer PlayPos
	trumpSuit       *Suit
	teams           [2]Team
	drawPile        DrawPile
	trick           Trick
	turnInfo        TurnInfo
	passCounter     int
}

func CreateGameState(id int) GameState {
	s := GameState{}
	s.id = id
	s.drawPile.Sprite = *graphics.CreateSprite(blankCardImage, blankCardAlphaImage, .1, screenWidth/2, screenHeight/2, 0, 0, 0, 0)
	s.drawPile.Sprite.Visible = false
	s.drawPile.Sprite.X = screenWidth/2 - s.drawPile.Sprite.ImageWidth - 20
	s.drawPile.Sprite.Y = screenHeight/2 - s.drawPile.Sprite.ImageHeight/2
	s.drawPile.shuffleDrawPile()

	s.teams[Black].teamColor = Black
	s.teams[Red].teamColor = Red
	s.teams[Black].points = 0
	s.teams[Red].points = 0
	s.teams[Black].players[0] = CreatePlayer(0, Black, 5, Bottom, .1, &s.drawPile /* faceDown */, false, 0)
	s.teams[Red].players[0] = CreatePlayer(1, Red, 5, Left, .1, &s.drawPile /* faceDown */, false, 0)
	s.teams[Black].players[1] = CreatePlayer(2, Black, 5, Top, .1, &s.drawPile /* faceDown */, false, 0)
	s.teams[Red].players[1] = CreatePlayer(3, Red, 5, Right, .1, &s.drawPile /* faceDown */, false, 0)

	s.trick.X = screenWidth/2 + 20
	s.trick.Y = screenHeight/2 - s.drawPile.Sprite.ImageHeight/2
	s.trick.playCard(s.drawPile.drawCard(.1, screenWidth/2+20, 0, 0 /*faceDown */, false))

	s.trumpDrawPlayer = 0
	s.activePlayer = s.GetNextPlayerByAbsPos(s.trumpDrawPlayer).AbsPos

	s.turnInfo.inited = false
	return s
}

type Game struct {
	inited     bool
	debug      bool
	server     server
	mode       mode
	id         int
	lobbyId    int
	fontSource *text.GoTextFaceSource
	pageStack  []Page
	toast      *Toast
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

	g.mode = LobbyTypeSelect

	g.id = -1
	g.lobbyId = -1

	g.fontSource = fontSource

	g.server.connected = true
	g.server.lobbyId = -1

	g.pageStack = append(g.pageStack, CreateLobbyTypeSelectPage(g))

	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:9000/ws", nil)

	if err != nil {
		println("Failed to initialize connection! No server found.")
	}

	g.server.server = connection.Server{
		Ctx:  ctx,
		WS:   conn,
		Data: make(chan connection.Message),
	}
	go g.server.server.Listen()
}

func (g *Game) debugPrintln(msg string) {
	if g.debug {
		println(msg)
	}
}

func (g *Game) PushPage(p Page) {
	g.pageStack = append(g.pageStack, p)
}

func (g *Game) PopPage() {
	g.pageStack = g.pageStack[:len(g.pageStack)-1]
}

func (g *Game) ShowToast(message string, duration time.Duration) {
	if duration > 0 {
		g.toast = &Toast{message: message, expires: time.Now().Add(duration)}
	}
}

func (s *GameState) GetPlayerByAbsPos(absPos PlayPos) *Player {
	for i := range s.teams {
		for j := range s.teams[i].players {
			if absPos == s.teams[i].players[j].AbsPos {
				return &s.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetPlayer(%v)` with no player matching that id present in player list", absPos)
	return nil
}

func (s *GameState) GetPlayerById(id int) *Player {
	for i := range s.teams {
		for j := range s.teams[i].players {
			if id == s.teams[i].players[j].Id {
				return &s.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetPlayer(%v)` with no player matching that id present in player list", id)
	return nil
}

func (s *GameState) GetNextPlayerByAbsPos(absPos PlayPos) *Player {
	nextPlayerPos := (absPos + 1) % 4

	for i := range s.teams {
		for j := range s.teams[i].players {
			if nextPlayerPos == s.teams[i].players[j].AbsPos {
				return &s.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetNextPlayerByAbsPos(%v)` and somehow could not find the player at the adjascent position.", absPos)
	return nil
}

func (s *GameState) GetNextPlayerById(id int) *Player {
	prevPlayerPos := s.GetPlayerById(id).AbsPos
	nextPlayerPos := (prevPlayerPos + 1) % 4

	for i := range s.teams {
		for j := range s.teams[i].players {
			if nextPlayerPos == s.teams[i].players[j].AbsPos {
				return &s.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetNextPlayerById(%v)` and somehow could not find the player at the adjascent position.", id)
	return nil
}

func (s *GameState) GetClient() *Player {
	return s.GetPlayerById(s.id)
}

func (s *GameState) GetActivePlayer() *Player {
	return s.GetPlayerByAbsPos(s.activePlayer)
}

func (s *GameState) SetActiveAbsPos(absPos PlayPos) {
	s.activePlayer = absPos
}

func (s *GameState) SetActivePlayer(player *Player) {
	s.activePlayer = player.AbsPos
}

func (s *GameState) GetTeam(player *Player) *Team {
	for i := range s.teams {
		for j := range s.teams[i].players {
			if &s.teams[i].players[j] == player {
				return &s.teams[i]
			}
		}
	}
	return nil
}

func (s *GameState) DealCards() {
	s.trick.clear()

	s.drawPile.shuffleDrawPile()

	for i := range s.teams {
		for j := range s.teams[i].players {
			if len(s.teams[i].players[j].Cards) <= 0 {
				faceDown := s.id != s.teams[i].players[j].Id
				s.teams[i].players[j].DealHand(.1, &s.drawPile, 5, faceDown)
			}
		}
	}
	s.trumpSuit = nil
	s.passCounter = 0

	s.ArrangeTeams()
}

func (s *GameState) ArrangeTeams() {
	client := s.GetClient()
	s.teams[Black].Arrange(client.Id, client.AbsPos)
	s.teams[Red].Arrange(client.Id, client.AbsPos)
}

func (g *Game) SendStateResponse() {

	gameActivePage, ok := g.pageStack[len(g.pageStack)-1].(*GameActivePage)
	if !ok {
		panic("Page is not `GameActivePage`!")
	}

	var gameState connection.StateResponse
	gameState = g.EncodeGameState()
	fmt.Printf("Player id: %v\n", gameState.PlayerId)
	if g.server.connected {
		g.server.server.Send(connection.Message{
			Type: "state_res",
			Data: gameState,
		})
	} else {
		println("state_res not sent since no server is connected")
	}

	gameActivePage.ready = true

	g.debugPrintln("Handled state request message!")
}

func (g *GameState) PickUpTrump(player *Player) {
	topCard := g.trick.Pile[len(g.trick.Pile)-1]
	topCard.PlayerId = player.Id
	g.trick.Pile = g.trick.Pile[:len(g.trick.Pile)-1]
	g.trumpSuit = &topCard.Suit

	player.Cards = append(player.Cards, topCard)
	player.ArrangeHand(g.GetClient().Id)
}

func (s *GameState) PlayCard(playerId int, cardInd int) {
	player := s.GetPlayerById(playerId)
	playedCard := player.Cards[cardInd]

	if len(s.trick.Pile) == 0 {
		s.trick.LeadSuit = playedCard.Suit
	}
	s.trick.playCard(playedCard)
	println("player", playerId, "just played card", SuitToString(playedCard.Suit), "of", NumberToString(playedCard.Number), "owned by", playedCard.PlayerId)
	player.Cards = slices.Delete(player.Cards, cardInd, cardInd+1)
	player.ArrangeHand(s.GetClient().Id)
}

func (s *GameState) EndTurn() {
	s.activePlayer = (s.activePlayer + 1) % 4
}

func (s *GameState) IsPickingTrump() bool {
	return s.activePlayer == s.GetClient().AbsPos && s.trumpSuit == nil
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
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

	if len(g.pageStack) == 0 {
		panic("Ran `Game.Update()` with an empty page stack!")
	}
	g.pageStack[len(g.pageStack)-1].Update()

	return nil
}

func (g *Game) drawToast(screen *ebiten.Image) {
	if g.toast == nil {
		return
	}
	if time.Now().After(g.toast.expires) {
		g.toast = nil
		return
	}

	face := &text.GoTextFace{Source: g.fontSource, Size: 22}
	lines := strings.Split(g.toast.message, "\n")
	lineWidths := make([]float64, len(lines))
	lineHeights := make([]float64, len(lines))
	var width, height float64
	for i, line := range lines {
		lineWidths[i], lineHeights[i] = text.Measure(line, face, 0)
		width = max(width, lineWidths[i])
		height += lineHeights[i]
	}
	const lineSpacing = 4
	height += float64(max(0, len(lines)-1) * lineSpacing)

	toast := ebiten.NewImage(int(width)+32, int(height)+20)
	toast.Fill(color.RGBA{0, 0, 0, 220})

	toastOp := ebiten.DrawImageOptions{}
	toastX := float64(screenWidth-int(width)-32) / 2
	toastY := float64(screenHeight - int(height) - 60)
	toastOp.GeoM.Translate(toastX, toastY)
	screen.DrawImage(toast, &toastOp)

	lineY := toastY + 10
	for i, line := range lines {
		textOp := text.DrawOptions{}
		textOp.GeoM.Translate(float64(screenWidth-int(lineWidths[i]))/2, lineY)
		text.Draw(screen, line, face, &textOp)
		lineY += lineHeights[i] + lineSpacing
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{161, 191, 123, 1})

	if len(g.pageStack) == 0 {
		panic("Ran `Game.Draw()` with an empty page stack!")
	}
	g.pageStack[len(g.pageStack)-1].Draw(screen)
	g.drawToast(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// *Game getters

func (g *Game) GetFontSource() *text.GoTextFaceSource {
	return g.fontSource
}

func (s *GameState) SetActivePlayerAbsPos(pos PlayPos) {
	s.activePlayer = pos
}

func (s *GameState) SetTrumpDrawPlayer(pos PlayPos) {
	s.trumpDrawPlayer = pos
}

func (g *Game) SendServerMessage(msgType string, data interface{}) {
	g.server.server.Send(connection.Message{
		Type: msgType,
		Data: data,
	})
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

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
