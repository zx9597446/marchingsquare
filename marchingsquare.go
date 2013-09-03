package marchingsquare

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

var _ = png.Decode
var _ = jpeg.Decode

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	directionNone = iota
	diretionUp
	directionLeft
	directionDown
	directionRight
)

type Point struct {
	X, Y int
}

type TestFunc func(r, g, b, a uint32) bool

type marchingSquare struct {
	img          image.Image
	previousStep int
	nextStep     int
	result       []Point
	test         TestFunc
}

func (m *marchingSquare) findStartPoint() (int, int) {
	b := m.img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := m.img.At(x, y).RGBA()
			if m.test(r, g, b, a) {
				return x, y
			}
		}
	}
	return 0, 0
}

func (m *marchingSquare) isPixelSolid(x, y int) bool {
	rect := m.img.Bounds()
	if x < rect.Min.X || x > rect.Max.X || y < rect.Min.Y || y > rect.Max.Y {
		return false
	}
	r, g, b, a := m.img.At(x, y).RGBA()
	if m.test(r, g, b, a) {
		return true
	}
	return false
}

func (m *marchingSquare) step(x, y int) {
	upLeft := m.isPixelSolid(x-1, y-1)
	upRight := m.isPixelSolid(x, y-1)
	downLeft := m.isPixelSolid(x-1, y)
	downRight := m.isPixelSolid(x, y)

	m.previousStep = m.nextStep

	state := 0
	if upLeft {
		state |= 1
	}
	if upRight {
		state |= 2
	}
	if downLeft {
		state |= 4
	}
	if downRight {
		state |= 8
	}

	switch state {
	case 1:
		m.nextStep = diretionUp
	case 2:
		m.nextStep = directionRight
	case 3:
		m.nextStep = directionRight
	case 4:
		m.nextStep = directionLeft
	case 5:
		m.nextStep = diretionUp
	case 6:
		if m.previousStep == diretionUp {
			m.nextStep = directionLeft
		} else {
			m.nextStep = directionRight
		}
	case 7:
		m.nextStep = directionRight
	case 8:
		m.nextStep = directionDown
	case 9:
		if m.previousStep == directionRight {
			m.nextStep = diretionUp
		} else {
			m.nextStep = directionDown
		}
	case 10:
		m.nextStep = directionDown
	case 11:
		m.nextStep = directionDown
	case 12:
		m.nextStep = directionLeft
	case 13:
		m.nextStep = diretionUp
	case 14:
		m.nextStep = directionLeft
	default:
		m.nextStep = directionNone
	}
}

func (m *marchingSquare) walk(startX, startY int) {
	b := m.img.Bounds()
	if startX < b.Min.X {
		startX = b.Min.X
	}
	if startX > b.Max.X {
		startX = b.Max.X
	}
	if startY < b.Min.Y {
		startY = b.Min.Y
	}
	if startY > b.Max.Y {
		startY = b.Max.Y
	}
	x, y := startX, startY
	for {
		m.step(x, y)
		if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
			m.result = append(m.result, Point{x, y})
		}
		switch m.nextStep {
		case diretionUp:
			y--
		case directionLeft:
			x--
		case directionDown:
			y++
		case directionRight:
			x++
		}
		if x == startX && y == startY {
			break
		}
	}
}

func (m *marchingSquare) doMarch(img image.Image, f TestFunc) []Point {
	m.img = img
	m.test = f
	m.result = make([]Point, 0)
	m.walk(m.findStartPoint())
	return m.result
}

func (m *marchingSquare) doMarchWithFileName(filename string, f TestFunc) []Point {
	file, err := os.Open(filename)
	defer file.Close()
	panicIfErr(err)
	img, _, err := image.Decode(file)
	panicIfErr(err)
	return m.doMarch(img, f)
}

func Process(img image.Image, f TestFunc) []Point {
	m := marchingSquare{}
	return m.doMarch(img, f)
}

func ProcessWithFile(filename string, f TestFunc) []Point {
	m := marchingSquare{}
	return m.doMarchWithFileName(filename, f)
}

func TransparentTest(r, g, b, a uint32) bool {
	return a > 0
}
