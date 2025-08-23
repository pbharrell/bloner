package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pbharrell/bloner/graphics"
	"slices"
)

var (
	ebitenImage *ebiten.Image
	cardImage   *ebiten.Image
)

func init() {
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

type Game struct {
	inited   bool
	touchIDs []ebiten.TouchID
	hand     *Hand
	oppHands [3]*Hand
	drawPile DrawPile
	trick    Trick
}

func (g *Game) getOppHand(pos PlayPos) *Hand {
	return g.oppHands[pos-1]
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	g.drawPile.Sprite = *graphics.CreateSpriteFromFile("./assets/ace_of_spades.png", .35, screenWidth/2, screenHeight/2, 0, 0, 0, 0)
	g.drawPile.Sprite.X = screenWidth/2 - g.drawPile.Sprite.ImageWidth - 20
	g.drawPile.Sprite.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
	g.drawPile.shuffleDrawPile()

	g.hand = CreateHand(5, Bottom, &g.drawPile)
	g.oppHands[0] = CreateHand(5, Left, &g.drawPile)
	g.oppHands[1] = CreateHand(5, Top, &g.drawPile)
	g.oppHands[2] = CreateHand(5, Right, &g.drawPile)

	g.trick.Pile = append(g.trick.Pile, g.drawPile.drawCard())
	g.trick.X = screenWidth/2 + 20
	g.trick.Y = screenHeight/2 - g.drawPile.Sprite.ImageHeight/2
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

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
			card := g.drawPile.drawCard()
			if card != nil {
				g.hand.Cards = append(g.hand.Cards, card)
				g.hand.ArrangeHand()
			}
		}
	}

	g.drawPile.Update()

	// Add a card to the mix.
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// println(g.hand.Cards[0])
		// card := CreateCard(Spades, Ace, .25, 50, 50, 0)
		// graphics.Sprite{
		// 	Image:       cardImage,
		// 	ImageWidth:  int(float64(cardImage.Bounds().Dx()) * .25),
		// 	ImageHeight: int(float64(cardImage.Bounds().Dy()) * .25),
		// 	ImageScale:  .25,
		// 	X:           20,
		// 	Y:           20,
		// 	Vx:          vx,
		// 	Vy:          vy,
		// 	Angle:       a,
		// }
		// g.cards = append(g.cards, card)
	}

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

	g.drawPile.Draw(screen, op)
	g.trick.Draw(screen, op)
	g.hand.Draw(screen, op)

	for _, hand := range g.oppHands {
		hand.Draw(screen, op)
	}
	// msg := fmt.Sprintf(`TPS: %0.2f
	// FPS: %0.2f
	// Num of sprites: %d
	// Press <- or -> to change the number of sprites`, ebiten.ActualTPS(), ebiten.ActualFPS(), len(g.hand.Draw))
	// ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
