package graphics

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxAngle = 360
)

type Sprite struct {
	Image       *ebiten.Image
	AlphaImage  *image.Alpha
	ImageWidth  int
	ImageHeight int
	ImageScale  float64
	X           int
	Y           int
	Angle       int
	Vx          int
	Vy          int
	Vangle      int
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

func CreateSprite(image *ebiten.Image, alphaImage *image.Alpha, scale float64, x int, y int, angle int, vx int, vy int, vangle int) *Sprite {
	rawW, rawH := image.Bounds().Dx(), image.Bounds().Dy()
	scaledW, scaledH := int(float64(rawW)*scale), int(float64(rawH)*scale)

	return &Sprite{
		Image:       image,
		AlphaImage:  alphaImage,
		ImageWidth:  scaledW,
		ImageHeight: scaledH,
		ImageScale:  scale,
		X:           x,
		Y:           y,
		Angle:       angle,
		Vx:          vx,
		Vy:          vy,
		Vangle:      vangle,
	}
}

func CreateSpriteFromFile(imageFile string, scale float64, x int, y int, angle int, vx int, vy int, vangle int) *Sprite {
	reader, err := os.Open(imageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	spriteImage := ebiten.NewImageFromImage(img)

	b := img.Bounds()
	spriteAlphaImage := image.NewAlpha(b)
	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			spriteAlphaImage.Set(i, j, img.At(i, j))
		}
	}

	return CreateSprite(spriteImage, spriteAlphaImage, scale, x, y, angle, vx, vy, vangle)
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
	s.Angle += s.Vangle
	if s.Angle == 360 {
		s.Angle = 0
	}
}

func (s *Sprite) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	w, h := s.ImageWidth, s.ImageHeight
	op.GeoM.Scale(s.ImageScale, s.ImageScale)
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(2 * math.Pi * float64(s.Angle) / maxAngle)
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	screen.DrawImage(s.Image, &op)
}

func (s *Sprite) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to users.
	//
	// Use alphaImage (*image.Alpha) instead of image (*ebiten.Image) here.
	// It is because (*ebiten.Image).At is very slow as this reads pixels from GPU,
	// and should be avoided whenever possible.
	return s.AlphaImage.At(int(float64(x-s.X)/s.ImageScale), int(float64(y-s.Y)/s.ImageScale)).(color.Alpha).A > 0
}
