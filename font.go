package hzk

import (
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	glyphWidth = 16
	gb2312Cols = 94
)

//go:embed fonts/HZK12
var embeddedHZK12 []byte

//go:embed fonts/HZK16
var embeddedHZK16 []byte

var ErrUnsupportedRune = errors.New("unsupported rune")

// Size identifies a bundled HZK bitmap font size.
type Size int

const (
	HZK12 Size = 12
	HZK16 Size = 16
)

// ParseSize parses "12" or "16" into a supported HZK size.
func ParseSize(s string) (Size, error) {
	switch s {
	case "12":
		return HZK12, nil
	case "16":
		return HZK16, nil
	default:
		return 0, fmt.Errorf("unsupported HZK size %q", s)
	}
}

func (s Size) String() string {
	if s == HZK12 || s == HZK16 {
		return strconv.Itoa(int(s))
	}
	return fmt.Sprintf("Size(%d)", int(s))
}

func (s Size) metrics() (height int, bytesPerGlyph int, err error) {
	switch s {
	case HZK12:
		return 12, 24, nil
	case HZK16:
		return 16, 32, nil
	default:
		return 0, 0, fmt.Errorf("unsupported HZK size %d", s)
	}
}

// Font converts supported GB2312 runes into HZK bitmap glyphs.
type Font struct {
	size          Size
	height        int
	bytesPerGlyph int
	data          []byte
}

// New returns a font backed by the bundled HZK12 or HZK16 data.
func New(size Size) (*Font, error) {
	switch size {
	case HZK12:
		return NewFromBytes(size, embeddedHZK12)
	case HZK16:
		return NewFromBytes(size, embeddedHZK16)
	default:
		return nil, fmt.Errorf("unsupported HZK size %d", size)
	}
}

// LoadFile loads an HZK font file from disk.
func LoadFile(size Size, path string) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(size, data)
}

// NewFromBytes creates a font from raw HZK bytes.
func NewFromBytes(size Size, data []byte) (*Font, error) {
	height, bytesPerGlyph, err := size.metrics()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("empty HZK data")
	}
	if len(data)%bytesPerGlyph != 0 {
		return nil, fmt.Errorf("invalid HZK%d data length %d: not a multiple of %d", size, len(data), bytesPerGlyph)
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	return &Font{
		size:          size,
		height:        height,
		bytesPerGlyph: bytesPerGlyph,
		data:          buf,
	}, nil
}

func (f *Font) Size() Size {
	return f.size
}

func (f *Font) Width() int {
	return glyphWidth
}

func (f *Font) Height() int {
	return f.height
}

// Contains reports whether r can be read from this font.
func (f *Font) Contains(r rune) bool {
	if r == ' ' || r == '\u3000' {
		return true
	}
	_, offset, ok := f.offset(r)
	return ok && offset+f.bytesPerGlyph <= len(f.data)
}

// Glyph returns the bitmap for a single supported rune.
func (f *Font) Glyph(r rune) (Bitmap, error) {
	if r == ' ' {
		return NewBlankBitmap(glyphWidth, f.height), nil
	}

	_, offset, ok := f.offset(r)
	if !ok || offset+f.bytesPerGlyph > len(f.data) {
		return Bitmap{}, fmt.Errorf("%w %q", ErrUnsupportedRune, r)
	}

	bitmap := NewBlankBitmap(glyphWidth, f.height)
	glyphData := f.data[offset : offset+f.bytesPerGlyph]
	for y := 0; y < f.height; y++ {
		for byteX := 0; byteX < 2; byteX++ {
			b := glyphData[y*2+byteX]
			for bit := 0; bit < 8; bit++ {
				if b&(0x80>>bit) != 0 {
					bitmap.Set(byteX*8+bit, y, true)
				}
			}
		}
	}
	return bitmap, nil
}

func (f *Font) offset(r rune) (uint16, int, bool) {
	code, ok := gb2312Code(r)
	if !ok {
		return 0, 0, false
	}

	area := int(code>>8) - 0xA0
	index := int(code&0xFF) - 0xA0
	if area < 1 || index < 1 || area > gb2312Cols || index > gb2312Cols {
		return 0, 0, false
	}

	offset := (gb2312Cols*(area-1) + (index - 1)) * f.bytesPerGlyph
	return code, offset, true
}

// Matrix converts text to a [][]bool dot matrix.
func (f *Font) Matrix(text string, options ...TextOptions) ([][]bool, error) {
	bitmap, err := f.Bitmap(text, options...)
	if err != nil {
		return nil, err
	}
	return bitmap.Matrix(), nil
}

// Bitmap converts text to a single bitmap. Newlines are supported.
func (f *Font) Bitmap(text string, options ...TextOptions) (Bitmap, error) {
	opt := mergeTextOptions(options)
	if opt.GlyphSpacing < 0 {
		return Bitmap{}, fmt.Errorf("glyph spacing must be >= 0")
	}
	if opt.LineSpacing < 0 {
		return Bitmap{}, fmt.Errorf("line spacing must be >= 0")
	}
	if opt.CellWidth < 0 {
		return Bitmap{}, fmt.Errorf("cell width must be >= 0")
	}
	if text == "" {
		return NewBlankBitmap(0, 0), nil
	}

	lines := strings.Split(text, "\n")

	lineGlyphs := make([][]Bitmap, len(lines))
	maxWidth := 0
	for lineIndex, line := range lines {
		glyphs := make([]Bitmap, 0, utf8.RuneCountInString(line))
		lineWidth := 0
		for _, r := range line {
			if !opt.DisableWiden {
				r = widenRune(r)
			}
			glyph, err := f.Glyph(r)
			if err != nil {
				return Bitmap{}, err
			}
			if len(glyphs) > 0 {
				lineWidth += opt.GlyphSpacing
			}
			lineWidth += glyphCellWidth(glyph, opt.CellWidth)
			glyphs = append(glyphs, glyph)
		}
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
		lineGlyphs[lineIndex] = glyphs
	}

	totalHeight := len(lines) * f.height
	if len(lines) > 1 {
		totalHeight += (len(lines) - 1) * opt.LineSpacing
	}
	result := NewBlankBitmap(maxWidth, totalHeight)

	y := 0
	for _, glyphs := range lineGlyphs {
		x := 0
		for glyphIndex, glyph := range glyphs {
			if glyphIndex > 0 {
				x += opt.GlyphSpacing
			}
			width := glyphCellWidth(glyph, opt.CellWidth)
			blit(result, glyph, x, y, width)
			x += width
		}
		y += f.height + opt.LineSpacing
	}

	return result, nil
}

func glyphCellWidth(glyph Bitmap, width int) int {
	if width > 0 {
		return width
	}
	return glyph.Width
}

func blit(dst Bitmap, src Bitmap, offsetX, offsetY, width int) {
	if width > src.Width {
		width = src.Width
	}
	for y := 0; y < src.Height; y++ {
		for x := 0; x < width; x++ {
			if src.At(x, y) {
				dst.Set(offsetX+x, offsetY+y, true)
			}
		}
	}
}

// Render converts text to printable dot-art using the provided options.
func (f *Font) Render(text string, options ...RenderOptions) (string, error) {
	opt := mergeRenderOptions(options)
	bitmap, err := f.Bitmap(text, TextOptions{
		GlyphSpacing: opt.GlyphSpacing,
		LineSpacing:  opt.LineSpacing,
		CellWidth:    opt.CellWidth,
		DisableWiden: opt.DisableWiden,
	})
	if err != nil {
		return "", err
	}
	return bitmap.Render(opt.Foreground, opt.Background), nil
}

// Image converts text to a pixel image.
func (f *Font) Image(text string, options ...ImageOptions) (image.Image, error) {
	opt := mergeImageOptions(options)
	bitmap, err := f.Bitmap(text, TextOptions{
		GlyphSpacing: opt.GlyphSpacing,
		LineSpacing:  opt.LineSpacing,
		CellWidth:    opt.CellWidth,
		DisableWiden: opt.DisableWiden,
	})
	if err != nil {
		return nil, err
	}
	return bitmap.Image(opt), nil
}

// EncodePNG writes text as a PNG image.
func (f *Font) EncodePNG(w io.Writer, text string, options ...ImageOptions) error {
	img, err := f.Image(text, options...)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

func widenRune(r rune) rune {
	if r == ' ' {
		return '\u3000'
	}
	if r >= 0x21 && r <= 0x7E {
		return r + 0xFEE0
	}
	return r
}
