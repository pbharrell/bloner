package graphics

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// whiteImage = ebiten.NewImage(3, 3)
	inited = false
)

type Shape struct {
	Image    *ebiten.Image
	Vertices []ebiten.Vertex
	Indices  []uint16
	X        int
	Y        int
	Angle    int
	Vx       int
	Vy       int
	Vangle   int
	Visible  bool
}

func GenIndices(numVertices int) []uint16 {
	indices := []uint16{}
	numVertices = numVertices - 1
	for i := range numVertices - 1 {
		indices = append(indices, uint16(i), uint16(i+1)%uint16(numVertices), uint16(numVertices))
	}
	return indices
}

func CreateShape(image *ebiten.Image, vertices []ebiten.Vertex, indices []uint16, x int, y int, angle int, vx int, vy int, vangle int) *Shape {
	return &Shape{
		Image:    image,
		Vertices: vertices,
		Indices:  indices,
		X:        x,
		Y:        y,
		Angle:    angle,
		Vx:       vx,
		Vy:       vy,
		Vangle:   vangle,
		Visible:  true,
	}
}

func (s *Shape) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe

	screen.DrawTriangles(s.Vertices, s.Indices, s.Image.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}
