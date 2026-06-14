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

type Game struct {
	inited          bool
	debug           bool
	server          server
	id              int
	mode            mode
	lobbyId         int
	fontSource      *text.GoTextFaceSource
	trumpSuit       *Suit
	passCounter     int
	teams           [2]Team
	activePlayer    PlayPos
	trumpDrawPlayer PlayPos
	buttonNewLobby  Button
	buttonJoinLobby Button
	lobbyRequestStr string
	overlay         graphics.Shape
	drawPile        DrawPile
	trick           Trick
	turnInfo        TurnInfo
	pageStack       []Page
}

func (g *Game) initOverlay() {
	g.overlay = *graphics.CreateRectangle(overlayImage, screenWidth, screenHeight, 0, 0, 0, 0, 0, 0)
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

	g.lobbyId = -1

	g.fontSource = fontSource

	g.drawPile.Sprite = *graphics.CreateSprite(blankCardImage, blankCardAlphaImage, .1, screenWidth/2, screenHeight/2, 0, 0, 0, 0)
	g.drawPile.Sprite.X = screenWidth/2 - g.drawPile.Sprite.ImageWidth - 20
	g.drawPile.Sprite.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.drawPile.shuffleDrawPile()

	g.teams[Black].teamColor = Black
	g.teams[Red].teamColor = Red
	g.teams[Black].points = 0
	g.teams[Red].points = 0
	g.teams[Black].players[0] = CreatePlayer(0, Black, 5, Bottom, .1, &g.drawPile /* faceDown */, false, 0)
	g.teams[Red].players[0] = CreatePlayer(1, Red, 5, Left, .1, &g.drawPile /* faceDown */, false, 0)
	g.teams[Black].players[1] = CreatePlayer(2, Black, 5, Top, .1, &g.drawPile /* faceDown */, false, 0)
	g.teams[Red].players[1] = CreatePlayer(3, Red, 5, Right, .1, &g.drawPile /* faceDown */, false, 0)

	g.trick.X = screenWidth/2 + 20
	g.trick.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.trick.playCard(g.drawPile.drawCard(.1, screenWidth/2+20, 0, 0 /*faceDown */, false))

	g.trumpDrawPlayer = 0
	g.activePlayer = g.GetNextPlayerByAbsPos(g.trumpDrawPlayer).AbsPos

	g.initOverlay()

	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:9000/ws", nil)

	if err != nil {
		return
	}

	g.server.server = connection.Server{
		Ctx:  ctx,
		WS:   conn,
		Data: make(chan connection.Message),
	}
	g.server.connected = true
	g.server.lobbyId = -1

	g.turnInfo.inited = false

	g.pageStack = append(g.pageStack, CreateLobbyTypeSelectPage(g))

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

func (g *Game) PopPage(p Page) {
	g.pageStack = g.pageStack[:len(g.pageStack)-2]
}

func (g *Game) GetPlayerByAbsPos(absPos PlayPos) *Player {
	for i := range g.teams {
		for j := range g.teams[i].players {
			if absPos == g.teams[i].players[j].AbsPos {
				return &g.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetPlayer(%v)` with no player matching that id present in player list", absPos)
	return nil
}

func (g *Game) GetPlayerById(id int) *Player {
	for i := range g.teams {
		for j := range g.teams[i].players {
			if id == g.teams[i].players[j].Id {
				return &g.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetPlayer(%v)` with no player matching that id present in player list", id)
	return nil
}

func (g *Game) GetNextPlayerByAbsPos(absPos PlayPos) *Player {
	nextPlayerPos := (absPos + 1) % 4

	for i := range g.teams {
		for j := range g.teams[i].players {
			if nextPlayerPos == g.teams[i].players[j].AbsPos {
				return &g.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetNextPlayerByAbsPos(%v)` and somehow could not find the player at the adjascent position.", absPos)
	return nil
}

func (g *Game) GetNextPlayerById(id int) *Player {
	prevPlayerPos := g.GetPlayerById(id).AbsPos
	nextPlayerPos := (prevPlayerPos + 1) % 4

	for i := range g.teams {
		for j := range g.teams[i].players {
			if nextPlayerPos == g.teams[i].players[j].AbsPos {
				return &g.teams[i].players[j]
			}
		}
	}

	fmt.Printf("ERROR: Should not be here! Called `game.GetNextPlayerById(%v)` and somehow could not find the player at the adjascent position.", id)
	return nil
}

func (g *Game) GetClient() *Player {
	return g.GetPlayerById(g.id)
}

func (g *Game) GetActivePlayer() *Player {
	return g.GetPlayerByAbsPos(g.activePlayer)
}

func (g *Game) SetActiveAbsPos(absPos PlayPos) {
	g.activePlayer = absPos
}

func (g *Game) SetActivePlayer(player *Player) {
	g.activePlayer = player.AbsPos
}

func (g *Game) GetTeam(player *Player) *Team {
	for i := range g.teams {
		for j := range g.teams[i].players {
			if &g.teams[i].players[j] == player {
				return &g.teams[i]
			}
		}
	}
	return nil
}

func (g *Game) DealCards() {
	g.trick.clear()

	g.drawPile.shuffleDrawPile()

	for i := range g.teams {
		for j := range g.teams[i].players {
			if len(g.teams[i].players[j].Cards) <= 0 {
				faceDown := g.id != g.teams[i].players[j].Id
				g.teams[i].players[j].DealHand(.1, &g.drawPile, 5, faceDown)
			}
		}
	}
	g.trumpSuit = nil
	g.passCounter = 0

	g.ArrangeTeams()
}

func (g *Game) ArrangeTeams() {
	client := g.GetClient()
	g.teams[Black].Arrange(client.Id, client.AbsPos)
	g.teams[Red].Arrange(client.Id, client.AbsPos)
}

func (g *Game) SendStateResponse() {
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

	g.debugPrintln("Handled state request message!")
}

func (g *Game) SendTurnInfo() {
	if g.server.connected && g.turnInfo.inited {
		println("Sending turn info for player id:", g.turnInfo.turnInfo.PlayerId)
		g.server.server.Send(connection.Message{
			Type: "turn_info",
			Data: g.turnInfo.turnInfo,
		})
	} else {
		println("turn_info not sent since no server is connected.")
	}
	g.turnInfo.inited = false
}

func (g *Game) SendTurnCardPlay(card *Card) {
	g.turnInfo.turnInfo.TurnType = connection.CardPlay
	g.turnInfo.turnInfo.CardPlay = card.Encode()
	g.SendTurnInfo()
	g.EndTurn()
}

func (g *Game) SendTurnTrumpDiscard(card *Card) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpDiscard
	g.turnInfo.turnInfo.TrumpDiscard = card.Encode()
	g.SendTurnInfo()
	g.EndTurn()
}

func (g *Game) SendTurnTrumpPick(suit int8) {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPick
	g.turnInfo.turnInfo.TrumpPick = suit
	g.SendTurnInfo()
}

func (g *Game) SendTurnTrumpPass() {
	g.turnInfo.turnInfo.TurnType = connection.TrumpPass
	g.SendTurnInfo()
	g.EndTurn()
}

func (g *Game) PickUpTrump(player *Player) {
	topCard := g.trick.Pile[len(g.trick.Pile)-1]
	topCard.PlayerId = player.Id
	g.trick.Pile = g.trick.Pile[:len(g.trick.Pile)-1]
	g.trumpSuit = &topCard.Suit

	player.Cards = append(player.Cards, topCard)
	player.ArrangeHand(g.GetClient().Id)
}

func (g *Game) PlayCard(playerId int, cardInd int) {
	player := g.GetPlayerById(playerId)
	playedCard := player.Cards[cardInd]

	if len(g.trick.Pile) == 0 {
		g.trick.LeadSuit = playedCard.Suit
	}
	g.trick.playCard(playedCard)
	println("player", playerId, "just played card", SuitToString(playedCard.Suit), "of", NumberToString(playedCard.Number), "owned by", playedCard.PlayerId)
	player.Cards = slices.Delete(player.Cards, cardInd, cardInd+1)
	player.ArrangeHand(g.GetClient().Id)
}

func (g *Game) EndTurn() {
	g.activePlayer = (g.activePlayer + 1) % 4
}

func (g *Game) IsPickingTrump() bool {
	return g.activePlayer == g.GetClient().AbsPos && g.trumpSuit == nil
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	if !g.turnInfo.inited {
		g.turnInfo.turnInfo = connection.TurnInfo{
			PlayerId: g.id,
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

	if len(g.pageStack) == 0 {
		panic("Ran `Game.Update()` with an empty page stack!")
	}
	g.pageStack[len(g.pageStack)-1].Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{161, 191, 123, 1})

	if len(g.pageStack) == 0 {
		panic("Ran `Game.Draw()` with an empty page stack!")
	}
	g.pageStack[len(g.pageStack)-1].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// *Game getters

func (g *Game) GetFontSource() *text.GoTextFaceSource {
	return g.fontSource
}

func (g *Game) GetLobbyId() int {
	return g.lobbyId
}

func (g *Game) GetActivePlayerAbsPos() PlayPos {
	return g.activePlayer
}

func (g *Game) GetTrumpSuit() *Suit {
	return g.trumpSuit
}

func (g *Game) GetPassCounter() int {
	return g.passCounter
}

func (g *Game) GetDebug() bool {
	return g.debug
}

func (g *Game) GetId() int {
	return g.id
}

func (g *Game) GetTrick() *Trick {
	return &g.trick
}

func (g *Game) GetDrawPile() *DrawPile {
	return &g.drawPile
}

func (g *Game) GetTeams() *[2]Team {
	return &g.teams
}

// func (g *Game) SetMode(m ) {
// 	g.mode = m
// }
//
// func (g *Game) GetMode() Mode {
// 	return g.mode
// }

func (g *Game) SetActivePlayerAbsPos(pos PlayPos) {
	g.activePlayer = pos
}

func (g *Game) SetTrumpDrawPlayer(pos PlayPos) {
	g.trumpDrawPlayer = pos
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
