package render_test

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/federico-paolillo/imagen/internal/render"
)

var (
	magenta = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	black   = color.RGBA{A: 255}
)

func TestLabel(t *testing.T) {
	t.Parallel()

	if got, want := render.Label(1024, 768), "1024x768"; got != want {
		t.Fatalf("Label = %q, want %q", got, want)
	}
}

func TestImage(t *testing.T) {
	t.Parallel()

	sizes := []struct{ w, h int }{{1024, 1024}, {640, 480}, {200, 1000}, {1000, 200}, {64, 64}}

	for _, s := range sizes {
		img, err := render.Image(s.w, s.h)
		if err != nil {
			t.Fatalf("Image(%d, %d): %v", s.w, s.h, err)
		}

		if got := img.Bounds().Size(); got != (image.Point{X: s.w, Y: s.h}) {
			t.Fatalf("size = %v, want %dx%d", got, s.w, s.h)
		}

		if got := img.RGBAAt(0, 0); got != magenta {
			t.Errorf("%dx%d: corner = %v, want magenta", s.w, s.h, got)
		}

		ink, ok := inkBounds(img)
		if !ok {
			t.Fatalf("%dx%d: no black text pixels", s.w, s.h)
		}

		if ink.Min.X < 0 || ink.Min.Y < 0 || ink.Max.X > s.w || ink.Max.Y > s.h {
			t.Errorf("%dx%d: text %v overflows image", s.w, s.h, ink)
		}

		// Allow a few pixels of slack for glyph rounding.
		slack := 3 + s.w/100
		if dx := ink.Min.X + ink.Max.X - s.w; abs(dx) > 2*slack {
			t.Errorf("%dx%d: text %v not horizontally centered", s.w, s.h, ink)
		}

		if dy := ink.Min.Y + ink.Max.Y - s.h; abs(dy) > 2*slack {
			t.Errorf("%dx%d: text %v not vertically centered", s.w, s.h, ink)
		}
	}
}

func TestImageRejectsBadDimensions(t *testing.T) {
	t.Parallel()

	for _, s := range []struct{ w, h int }{{0, 10}, {10, 0}, {-1, 10}, {render.MaxDimension + 1, 10}, {10, render.MaxDimension + 1}} {
		if _, err := render.Image(s.w, s.h); !errors.Is(err, render.ErrDimension) {
			t.Errorf("Image(%d, %d) error = %v, want ErrDimension", s.w, s.h, err)
		}
	}
}

func inkBounds(img *image.RGBA) (image.Rectangle, bool) {
	var r image.Rectangle

	found := false

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.RGBAAt(x, y) != black {
				continue
			}

			p := image.Rect(x, y, x+1, y+1)
			if !found {
				r, found = p, true
			} else {
				r = r.Union(p)
			}
		}
	}

	return r, found
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
