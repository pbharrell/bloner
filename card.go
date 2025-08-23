package main

import (
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner/graphics"
)

type Suit uint8
type Number uint8

const (
	Spades Suit = iota
	Clubs
	Hearts
	Diamonds
)

const (
	Nine Number = iota
	Ten
	Jack
	Queen
	King
	Ace
	AltBauer
)

var (
	OffVal = map[Number]uint8{
		Nine:  0,
		Ten:   1,
		Jack:  2,
		Queen: 3,
		King:  4,
		Ace:   5,
	}

	TrumpColorVal = map[Number]uint8{
		Nine:  0,
		Ten:   1,
		Queen: 2,
		King:  3,
		Ace:   4,
	}

	TrumpVal = map[Number]uint8{
		Nine:     0,
		Ten:      1,
		Queen:    2,
		King:     3,
		Ace:      4,
		AltBauer: 5,
		Jack:     6,
	}
)

type Card struct {
	Sprite graphics.Sprite
	Suit   Suit
	Number Number
}

var (
	cardImages      [][]*ebiten.Image // Think of each row as the suit, each col as the num.
	cardAlphaImages [][]*image.Alpha
)

func initCardImages() {
	cardImages = make([][]*ebiten.Image, 4) // <-- the number of suits in play
	cardAlphaImages = make([][]*image.Alpha, len(cardImages))
	for i := range cardAlphaImages {
		cardImages[i] = make([]*ebiten.Image, 6) // <-- the number of distinct nums in play
		cardAlphaImages[i] = make([]*image.Alpha, len(cardImages[0]))

		for j := range cardAlphaImages[i] {
			// Read the file into a byte array
			var (
				cardImage      *ebiten.Image
				cardAlphaImage *image.Alpha
				err            error
			)

			reader, err := os.Open("./assets/ace_of_spades.png")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			img, _, err := image.Decode(reader)
			if err != nil {
				log.Fatal(err)
			}
			cardImage = ebiten.NewImageFromImage(img)

			b := img.Bounds()
			cardAlphaImage = image.NewAlpha(b)
			for j := b.Min.Y; j < b.Max.Y; j++ {
				for i := b.Min.X; i < b.Max.X; i++ {
					cardAlphaImage.Set(i, j, img.At(i, j))
				}
			}

			cardImages[i][j] = cardImage
			cardAlphaImages[i][j] = cardAlphaImage
		}
	}
}

func CreateCard(suit Suit, number Number, scale float64, x int, y int, angle int) *Card {
	if len(cardImages) < 4 || len(cardImages[0]) < 6 {
		panic("`cardImages` of unexpected size! Please call `InitCardImages()` first!")
	}

	return &Card{
		// Sprite: *graphics.CreateSpriteFromFile("./assets/ace_of_spades.png", .35, screenWidth/2, screenHeight/2, 0, 0, 0, 0),
		Sprite: *graphics.CreateSprite(cardImages[suit][number], cardAlphaImages[suit][number], scale, x, y, angle, 0, 0, 0),
		Suit:   suit,
		Number: number,
	}
}

func (c *Card) Update() {
	c.Sprite.Update()
}

func (c *Card) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	c.Sprite.Draw(screen, op)
}
