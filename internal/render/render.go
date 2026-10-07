// Package render draws the solid magenta placeholder image with its size label.
package render

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// MaxDimension is the largest width or height Image accepts.
const MaxDimension = 8192

const (
	// textFill is the fraction of the image the label may occupy on each axis.
	textFill = 0.8
	// referenceSize is the point size used to measure the label before scaling.
	referenceSize = 100.0
	dpi           = 72
)

var (
	magenta = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	black   = color.RGBA{A: 255}

	fontOnce   sync.Once
	parsedFont *opentype.Font
	errFont    error
)

// ErrDimension reports a width or height outside 1..MaxDimension.
var ErrDimension = errors.New("dimension out of range")

// Label returns the text drawn on an image of the given size, e.g. "1024x1024".
func Label(width, height int) string {
	return fmt.Sprintf("%dx%d", width, height)
}

// Image returns a solid magenta image of the given size with its Label
// centered in black.
func Image(width, height int) (*image.RGBA, error) {
	if width < 1 || width > MaxDimension || height < 1 || height > MaxDimension {
		return nil, fmt.Errorf("%w: width and height must be between 1 and %d, got %dx%d",
			ErrDimension, MaxDimension, width, height)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(magenta), image.Point{}, draw.Src)

	if err := drawLabel(img, Label(width, height)); err != nil {
		return nil, err
	}

	return img, nil
}

func loadFont() (*opentype.Font, error) {
	fontOnce.Do(func() {
		parsedFont, errFont = opentype.Parse(gobold.TTF)
	})

	if errFont != nil {
		return nil, fmt.Errorf("parse font: %w", errFont)
	}

	return parsedFont, nil
}

func newFace(f *opentype.Font, size float64) (font.Face, error) {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, fmt.Errorf("create font face: %w", err)
	}

	return face, nil
}

func drawLabel(img *image.RGBA, text string) error {
	f, err := loadFont()
	if err != nil {
		return err
	}

	size, err := fitSize(f, text, img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return err
	}

	face, err := newFace(f, size)
	if err != nil {
		return err
	}
	defer func() { _ = face.Close() }()

	// Bounds are relative to the dot, so offsetting by their midpoint
	// centers the ink, not just the advance box.
	bounds, _ := font.BoundString(face, text)
	center := fixed.Point26_6{
		X: fixed.I(img.Bounds().Dx()) / 2,
		Y: fixed.I(img.Bounds().Dy()) / 2,
	}

	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(black),
		Face: face,
		Dot: fixed.Point26_6{
			X: center.X - (bounds.Min.X+bounds.Max.X)/2,
			Y: center.Y - (bounds.Min.Y+bounds.Max.Y)/2,
		},
	}
	d.DrawString(text)

	return nil
}

// fitSize returns the largest point size at which text fits within textFill
// of the width and height.
func fitSize(f *opentype.Font, text string, width, height int) (float64, error) {
	ref, err := newFace(f, referenceSize)
	if err != nil {
		return 0, err
	}
	defer func() { _ = ref.Close() }()

	bounds, _ := font.BoundString(ref, text)
	textW := (bounds.Max.X - bounds.Min.X).Ceil()
	textH := (bounds.Max.Y - bounds.Min.Y).Ceil()

	scale := min(textFill*float64(width)/float64(textW), textFill*float64(height)/float64(textH))

	return max(referenceSize*scale, 1), nil
}
