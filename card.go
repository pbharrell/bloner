package main

import (
	"fmt"
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner-server/connection"
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
	PrimBauer
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

	TrumpVal = map[Number]uint8{
		Nine:      0,
		Ten:       1,
		Queen:     2,
		King:      3,
		Ace:       4,
		AltBauer:  5,
		PrimBauer: 6,
	}
)

type Card struct {
	Sprite   graphics.Sprite
	Suit     Suit
	Number   Number
	PlayerId int
	FaceDown bool
}

var (
	cardImages          [][]*ebiten.Image // Think of each row as the suit, each col as the num.
	cardAlphaImages     [][]*image.Alpha
	blankCardImage      *ebiten.Image
	blankCardAlphaImage *image.Alpha
	cardImageFilenames  [][]string
	blankImageFilename  string
	// TODO: Blank side image
)

func SuitToString(suit Suit) string {
	switch suit {
	case Spades:
		return "spades"
	case Clubs:
		return "clubs"
	case Hearts:
		return "hearts"
	case Diamonds:
		return "diamonds"
	default:
		return fmt.Sprintf("Found unsupported suit with value %d", suit)
	}
}

func NumberToString(num Number) string {
	switch num {
	case Nine:
		return "nine"
	case Ten:
		return "ten"
	case Jack:
		return "jack"
	case Queen:
		return "queen"
	case King:
		return "king"
	case Ace:
		return "ace"
	default:
		return fmt.Sprintf("Found unsupported suit with value %d", num)
	}
}

func initCardImageFiles() {
	// One image for each card + blank side
	allowedImageFiles := []string{"assets/ace_of_spades.png", "assets/ten_of_clubs.png", "assets/jack_of_clubs.png"}

	cardImageFilenames = make([][]string, 4) // <-- the number of suits in play + 1 for blank side
	for i := range cardImageFilenames {
		cardImageFilenames[i] = make([]string, 6) // <-- the number of distinct nums in play

		for j := range cardImageFilenames[i] {
			cardImageFilenames[i][j] = "assets/" + NumberToString(Number(j)) + "_of_" + SuitToString(Suit(i)) + ".png"

			// TODO: Change the image overriding when other images are in place
			if !slices.Contains(allowedImageFiles, cardImageFilenames[i][j]) {
				cardImageFilenames[i][j] = allowedImageFiles[0]
			}
		}
	}

	// FIXME: THIS IS FACE-DOWN CARD IMAGE SOMEDAY
	blankImageFilename = "assets/blank.png"
}

func initCardImages() {
	initCardImageFiles()
	cardImages = make([][]*ebiten.Image, 4) // <-- the number of suits in play
	cardAlphaImages = make([][]*image.Alpha, len(cardImages))
	for i := range cardAlphaImages {
		cardImages[i] = make([]*ebiten.Image, 6) // <-- the number of distinct nums in play
		cardAlphaImages[i] = make([]*image.Alpha, len(cardImages[0]))

		for j := range cardAlphaImages[i] {
			// Read the file into a byte array
			cardImage, cardAlphaImage := graphics.LoadImageFromFile(&content, cardImageFilenames[i][j])
			cardImages[i][j] = cardImage
			cardAlphaImages[i][j] = cardAlphaImage
		}
	}

	blankCardImage, blankCardAlphaImage = graphics.LoadImageFromFile(&content, blankImageFilename)
}

func CreateCard(suit Suit, number Number, playerId int, scale float64, x int, y int, angle int, faceDown bool) *Card {
	if len(cardImages) < 4 || len(cardImages[0]) < 6 {
		panic("`cardImages` of unexpected size! Please call `InitCardImages()` first!")
	}

	c := &Card{
		Sprite:   *graphics.CreateSprite(cardImages[suit][number], cardAlphaImages[suit][number], scale, x, y, angle, 0, 0, 0),
		PlayerId: playerId,
		Suit:     suit,
		Number:   number,
		FaceDown: faceDown,
	}

	c.UpdateSprite()
	return c
}

func (c *Card) Encode() connection.Card {
	return connection.Card{
		Suit:   uint8(c.Suit),
		Number: uint8(c.Number),
	}
}

func (c *Card) UpdateSprite() {
	if c.FaceDown {
		c.Sprite.Image = blankCardImage
		c.Sprite.AlphaImage = blankCardAlphaImage
	} else {
		c.Sprite.Image = cardImages[c.Suit][c.Number]
		c.Sprite.AlphaImage = cardAlphaImages[c.Suit][c.Number]
	}
}

func (c *Card) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	c.Sprite.Draw(screen, op)
}
