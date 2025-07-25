package graphics

import ()

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
}

func (s *Sprite) Update() {
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if mx := 50 - s.imageWidth; mx <= s.x {
		s.x = 2*mx - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if my := 50 - s.imageHeight; my <= s.y {
		s.y = 2*my - s.y
		s.vy = -s.vy
	}
	s.angle++
	if s.angle == 360 {
		s.angle = 0
	}
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for i := 0; i < s.num; i++ {
		s.sprites[i].Update()
	}
}
