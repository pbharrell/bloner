package graphics

import (
	"bytes"
	"image"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

const (
	maxAngle = 360
)

type Sprite struct {
	Image       *ebiten.Image
	ImageWidth  int
	ImageHeight int
	ImageScale  float64
	X           int
	Y           int
	Vx          int
	Vy          int
	Angle       int
}

func LoadImage(i *[]byte) *ebiten.Image {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(*i))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	s := origEbitenImage.Bounds().Size()
	ei := ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.5)

	ei.DrawImage(origEbitenImage, op)
	return ei
}

func CreateSprite(image *ebiten.Image, scale float64, x int, y int, vx int, vy int, angle int) *Sprite {
	rawW, rawH := image.Bounds().Dx(), image.Bounds().Dy()
	scaledW, scaledH := int(float64(rawW)*scale), int(float64(rawH)*scale)

	return &Sprite{
		Image:       image,
		ImageWidth:  scaledW,
		ImageHeight: scaledH,
		ImageScale:  scale,
		X:           x,
		Y:           y,
		Vx:          vx,
		Vy:          vy,
		Angle:       angle,
	}
}

func DefaultSprite() Sprite {
	ebitenImage := LoadImage(&images.Ebiten_png)
	w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	x, y := rand.IntN(320-w), rand.IntN(240-h)
	vx, vy := 1, 1
	a := rand.IntN(maxAngle)

	return Sprite{
		Image:       ebitenImage,
		ImageWidth:  w,
		ImageHeight: h,
		ImageScale:  1,
		X:           x,
		Y:           y,
		Vx:          vx,
		Vy:          vy,
		Angle:       a,
	}
}

func (s *Sprite) Update() {
	s.X += s.Vx
	s.Y += s.Vy
	if s.X < 0 {
		s.X = -s.X
		s.Vx = -s.Vx
	} else if mx := 320 - s.ImageWidth; mx <= s.X {
		s.X = 2*mx - s.X
		s.Vx = -s.Vx
	}
	if s.Y < 0 {
		s.Y = -s.Y
		s.Vy = -s.Vy
	} else if my := 240 - s.ImageHeight; my <= s.Y {
		s.Y = 2*my - s.Y
		s.Vy = -s.Vy
	}
	s.Angle++
	if s.Angle == 360 {
		s.Angle = 0
	}
}

func (s *Sprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	w, h := s.ImageWidth, s.ImageHeight
	op.GeoM.Scale(s.ImageScale, s.ImageScale)
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(2 * math.Pi * float64(s.Angle) / maxAngle)
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	screen.DrawImage(s.Image, op)
}

type Card struct {
	Sprite
}
