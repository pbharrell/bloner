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

// func GenShapeVertices(numSides int, x int, y int, width int, height int, angle int) []ebiten.Vertex {
// 	var vertices []ebiten.Vertex
// 	// fAngle := float64(angle)
// 	// fX := float64(x)
// 	// fY := float64(y)
// 	for i := range numSides {
// 		// num := float64(i) / float64(numSides) // - 1/float64(numSides*2)
//
// 		// fWidth := float64(width)
// 		// fHeight := float64(height)
//
// 		vertices = append(vertices, ebiten.Vertex{
// 			DstX:   float32(100 * (i % 2)),       //(fWidth)*math.Cos(2*math.Pi*fAngle/360) - (fHeight)*math.Sin(2*math.Pi*fAngle/360)),
// 			DstY:   float32(100 * ((i + 1) % 2)), //(fWidth)*math.Sin(2*math.Pi*fAngle/360) + (fHeight)*math.Cos(2*math.Pi*fAngle/360)),
// 			SrcX:   float32(i),
// 			SrcY:   float32(i%2 + 1),
// 			ColorR: float32(1),
// 			ColorG: float32(1),
// 			ColorB: float32(1),
// 			ColorA: 1,
// 		})
// 	}
//
// 	return vertices
// }

func GenIndices(numVertices int) []uint16 {
	indices := []uint16{}
	numVertices = numVertices - 1
	for i := range numVertices - 1 {
		indices = append(indices, uint16(i), uint16(i+1)%uint16(numVertices), uint16(numVertices))
		// println(uint16(i), uint16(i+1)%uint16(numVertices), uint16(numVertices))
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
