package hzk

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"
)

// Bitmap is a row-major boolean dot matrix.
type Bitmap struct {
	Width  int
	Height int
	Pixels []bool
}

// NewBlankBitmap creates an empty bitmap with all pixels off.
func NewBlankBitmap(width, height int) Bitmap {
	return Bitmap{
		Width:  width,
		Height: height,
		Pixels: make([]bool, width*height),
	}
}

func (b Bitmap) At(x, y int) bool {
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		return false
	}
	return b.Pixels[y*b.Width+x]
}

func (b Bitmap) Set(x, y int, on bool) {
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		return
	}
	b.Pixels[y*b.Width+x] = on
}

// Matrix returns a copy of the bitmap as [][]bool.
func (b Bitmap) Matrix() [][]bool {
	rows := make([][]bool, b.Height)
	for y := 0; y < b.Height; y++ {
		row := make([]bool, b.Width)
		copy(row, b.Pixels[y*b.Width:(y+1)*b.Width])
		rows[y] = row
	}
	return rows
}

// Render converts the bitmap into text using foreground for on pixels and
// background for off pixels.
func (b Bitmap) Render(foreground, background string) string {
	var out strings.Builder
	for y := 0; y < b.Height; y++ {
		if y > 0 {
			out.WriteByte('\n')
		}
		for x := 0; x < b.Width; x++ {
			if b.At(x, y) {
				out.WriteString(foreground)
			} else {
				out.WriteString(background)
			}
		}
	}
	return out.String()
}

// Image converts the bitmap into a pixel image.
func (b Bitmap) Image(options ...ImageOptions) image.Image {
	opt := mergeImageOptions(options)

	width := b.Width*opt.Scale + opt.Padding*2
	height := b.Height*opt.Scale + opt.Padding*2
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, opt.Background)
		}
	}

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			if !b.At(x, y) {
				continue
			}
			drawCell(img, opt.Padding+x*opt.Scale, opt.Padding+y*opt.Scale, opt.Scale, opt.Foreground)
		}
	}
	return img
}

func drawCell(img *image.RGBA, x, y, size int, c color.Color) {
	for yy := y; yy < y+size; yy++ {
		for xx := x; xx < x+size; xx++ {
			img.Set(xx, yy, c)
		}
	}
}

// EncodePNG writes the bitmap as a PNG image.
func (b Bitmap) EncodePNG(w io.Writer, options ...ImageOptions) error {
	return png.Encode(w, b.Image(options...))
}
