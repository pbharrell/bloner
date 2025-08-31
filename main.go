package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"

	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type player uint8

const (
	Main player = iota
	LeftOpp
	TopOpp
	RightOpp
)

type Game struct {
	inited        bool
	trumpSuit     *Suit
	activePlayer  player
	touchIDs      []ebiten.TouchID
	buttonConfirm Button
	buttonCancel  Button
	overlay       graphics.Shape
	hand          *Hand
	oppHands      [3]*Hand
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

	g.drawPile.Sprite = *graphics.CreateSpriteFromFile("./assets/ace_of_spades.png", .35, screenWidth/2, screenHeight/2, 0, 0, 0, 0)
	g.drawPile.Sprite.X = screenWidth/2 - g.drawPile.Sprite.ImageWidth - 20
	g.drawPile.Sprite.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.drawPile.shuffleDrawPile()

	g.hand = CreateHand(5, Bottom, .35, &g.drawPile)
	g.oppHands[0] = CreateHand(5, Left, .35, &g.drawPile)
	g.oppHands[1] = CreateHand(5, Top, .35, &g.drawPile)
	g.oppHands[2] = CreateHand(5, Right, .35, &g.drawPile)

	g.trick.Pile = append(g.trick.Pile, g.drawPile.drawCard(.35, 0, 0, 0))
	g.trick.X = screenWidth/2 + 20
	g.trick.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2

	g.activePlayer = LeftOpp

	g.buttonConfirm = *CreateButton(buttonConfirmImage, buttonConfirmAlpha, buttonPressedConfirmImage, buttonPressedConfirmAlpha, 4, 0, screenHeight/2+80, 0)
	confirmWidth := g.buttonConfirm.sprite.ImageWidth
	confirmX := screenWidth/2 - confirmWidth/2 + 80
	g.buttonConfirm.SetLoc(confirmX, g.buttonConfirm.sprite.Y)

	g.buttonCancel = *CreateButton(buttonCancelImage, buttonCancelAlpha, buttonPressedCancelImage, buttonPressedCancelAlpha, 4, 0, screenHeight/2+80, 0)
	cancelWidth := g.buttonCancel.sprite.ImageWidth
	cancelX := screenWidth/2 - cancelWidth/2 - 80
	g.buttonCancel.SetLoc(cancelX, g.buttonCancel.sprite.Y)

	g.initOverlay()
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	if g.trumpSuit == nil {
		x, y := ebiten.CursorPosition()
		mouseButtonPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

		g.buttonConfirm.Update(x, y, mouseButtonPressed)
		g.buttonCancel.Update(x, y, mouseButtonPressed)

	} else {
		x, y := ebiten.CursorPosition()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

			// Look through sprites in reverse order since a card on the right is on top
			for i := len(g.hand.Cards) - 1; i >= 0; i-- {
				card := g.hand.Cards[i]
				if card.Sprite.In(x, y) {
					g.trick.playCard(g.hand.Cards[i])
					g.hand.Cards = slices.Delete(g.hand.Cards, i, i+1)
					g.hand.ArrangeHand()
					break
				}
			}

			if g.drawPile.Sprite.In(x, y) && len(g.hand.Cards) < 5 {
				card := g.drawPile.drawCard(.35, 0, 0, 0)
				if card != nil {
					g.hand.Cards = append(g.hand.Cards, card)
					g.hand.ArrangeHand()
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

	if g.trumpSuit != nil {
		g.hand.Draw(screen, op)
	}

	g.drawPile.Draw(screen, op)
	g.trick.Draw(screen, op)

	for _, hand := range g.oppHands {
		hand.Draw(screen, op)
	}

	// TODO: AND ACCOUNT FOR DROPPING OVERLAY WHEN NOT YOUR TURN
	if g.trumpSuit == nil {
		g.overlay.Draw(screen)

		// **Everything on top of fade overlay start here**

		g.buttonConfirm.Draw(screen, op)
		g.buttonCancel.Draw(screen, op)
		g.hand.Draw(screen, op)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("bloner")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
