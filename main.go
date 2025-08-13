package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/pbharrell/bloner/graphics"
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
}

const (
	screenWidth  = 320
	screenHeight = 240
	maxAngle     = 360
)

type Game struct {
	touchIDs []ebiten.TouchID
	sprites  []*graphics.Sprite
	inited   bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	// g.sprites = make([]*graphics.Sprite, 0)
}

func (g *Game) leftTouched() bool {
	for _, id := range g.touchIDs {
		x, _ := ebiten.TouchPosition(id)
		if x < screenWidth/2 {
			return true
		}
	}
	return false
}

func (g *Game) rightTouched() bool {
	for _, id := range g.touchIDs {
		x, _ := ebiten.TouchPosition(id)
		if x >= screenWidth/2 {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	// Decrease the number of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || g.leftTouched() {
		if len(g.sprites) > 0 {
			g.sprites = g.sprites[:len(g.sprites)-1]
		}
	}

	// Increase the number of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || g.rightTouched() {
		defaultSprite := graphics.DefaultSprite()
		g.sprites = append(g.sprites, &defaultSprite)
	}

	// Add a card to the mix.
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
		// x, y := rand.IntN(screenWidth-w), rand.IntN(screenHeight-h)
		vx, vy := 1, 1
		a := rand.IntN(maxAngle)
		// CreateSprite(cardImage,
		card := graphics.Sprite{
			Image:       cardImage,
			ImageWidth:  int(float64(cardImage.Bounds().Dx()) * .25),
			ImageHeight: int(float64(cardImage.Bounds().Dy()) * .25),
			ImageScale:  .25,
			X:           20,
			Y:           20,
			Vx:          vx,
			Vy:          vy,
			Angle:       a,
		}
		g.sprites = append(g.sprites, &card)
	}

	for _, sprite := range g.sprites {
		sprite.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	// w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	for i := range g.sprites {
		op := ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(0.5)
		g.sprites[i].Draw(screen, &op)
	}
	msg := fmt.Sprintf(`TPS: %0.2f
	FPS: %0.2f
	Num of sprites: %d
	Press <- or -> to change the number of sprites`, ebiten.ActualTPS(), ebiten.ActualFPS(), len(g.sprites))
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
