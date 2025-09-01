package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner/graphics"
)

type ButtonPressCallback func(g *Game)

type Button struct {
	game          *Game
	pressCallback ButtonPressCallback
	sprite        graphics.Sprite
	pressedSprite graphics.Sprite
	isHovered     bool
	isPressed     bool
}

func CreateButton(game *Game, pressCallback ButtonPressCallback, image *ebiten.Image, alphaImage *image.Alpha, pressedImage *ebiten.Image, pressedAlphaImage *image.Alpha, scale float64, x int, y int, angle int) *Button {
	return &Button{
		game:          game,
		pressCallback: pressCallback,
		sprite:        *graphics.CreateSprite(image, alphaImage, scale, x, y, angle, 0, 0, 0),
		pressedSprite: *graphics.CreateSprite(pressedImage, pressedAlphaImage, scale, x, y, angle, 0, 0, 0),
		isHovered:     false,
		isPressed:     false,
	}
}

func (b *Button) SetLoc(x int, y int) {
	b.sprite.X = x
	b.pressedSprite.X = x
	b.sprite.Y = y
	b.pressedSprite.Y = y
}

func (b *Button) Update(x int, y int, isMouseClick bool) {
	// Check if mouse is over button
	b.isHovered = b.sprite.In(x, y)

	// Check if button is clicked
	if b.isHovered && isMouseClick {
		b.isPressed = true
		b.pressCallback(b.game)
		println("Button clicked!")
	} else if b.isPressed {
		b.isPressed = false
	}
}

func (b *Button) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if !b.isHovered {
		b.sprite.Draw(screen, op)
	} else {
		b.pressedSprite.Draw(screen, op)
	}
}
