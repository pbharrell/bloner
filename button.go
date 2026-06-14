package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbharrell/bloner/graphics"
)

type ButtonPressCallback func(b *Button, p Page, g *Game)

type Button struct {
	page          Page
	game          *Game
	pressCallback ButtonPressCallback
	sprite        graphics.Sprite
	pressedSprite graphics.Sprite
	isHovered     bool
	isPressed     bool
}

func CreateButton(page Page, game *Game, pressCallback ButtonPressCallback, spritePath string, pressedSpritePath string, scale float64, x int, y int, angle int) *Button {
	return &Button{
		page:          page,
		game:          game,
		pressCallback: pressCallback,
		sprite:        *CreateSpriteFromAssetPath(spritePath, scale, x, y, angle, 0, 0, 0),
		pressedSprite: *CreateSpriteFromAssetPath(pressedSpritePath, scale, x, y, angle, 0, 0, 0),
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
	b.isHovered = b.sprite.In(x, y)

	if b.isHovered && isMouseClick {
		b.isPressed = true
		b.pressCallback(b, b.page, b.game)
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
