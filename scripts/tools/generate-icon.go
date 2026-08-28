//go:build ignore

// Regenerates workflow/icon.png, a placeholder pending real design —
// run with: go run scripts/tools/generate-icon.go workflow/icon.png
package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	const size = 512
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	bg := color.RGBA{0x3A, 0x3F, 0x5C, 0xFF} // slate blue
	fg := color.RGBA{0xF5, 0xF6, 0xFA, 0xFF} // near white
	radius := 96.0

	cx, cy := float64(size)/2, float64(size)/2

	inRoundedSquare := func(x, y float64) bool {
		half := float64(size)/2 - 8
		dx := math.Abs(x-cx) - (half - radius)
		dy := math.Abs(y-cy) - (half - radius)
		if dx < 0 {
			dx = 0
		}
		if dy < 0 {
			dy = 0
		}
		return math.Abs(x-cx) <= half && math.Abs(y-cy) <= half && dx*dx+dy*dy <= radius*radius
	}

	// eye outline: an ellipse ring, plus a diagonal slash — "make the
	// invisible visible, then strike it out."
	ellipseA, ellipseB := 150.0, 90.0
	ringWidth := 22.0

	onEllipseRing := func(x, y float64) bool {
		dx, dy := x-cx, y-cy
		v := (dx*dx)/(ellipseA*ellipseA) + (dy*dy)/(ellipseB*ellipseB)
		inner := (dx*dx)/((ellipseA-ringWidth)*(ellipseA-ringWidth)) + (dy*dy)/((ellipseB-ringWidth)*(ellipseB-ringWidth))
		return v <= 1.0 && inner >= 1.0
	}

	pupilRadius := 46.0
	onPupil := func(x, y float64) bool {
		dx, dy := x-cx, y-cy
		return dx*dx+dy*dy <= pupilRadius*pupilRadius
	}

	slashHalfWidth := 16.0
	onSlash := func(x, y float64) bool {
		// line from top-left to bottom-right of the eye's bounding box
		d := (x - y - (cx - cy))
		return math.Abs(d)/math.Sqrt2 <= slashHalfWidth &&
			x >= cx-ellipseA-10 && x <= cx+ellipseA+10 &&
			y >= cy-ellipseB-10 && y <= cy+ellipseB+10
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			switch {
			case !inRoundedSquare(fx, fy):
				img.Set(x, y, color.RGBA{})
			case onSlash(fx, fy), onEllipseRing(fx, fy), onPupil(fx, fy):
				img.Set(x, y, fg)
			default:
				img.Set(x, y, bg)
			}
		}
	}

	f, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
