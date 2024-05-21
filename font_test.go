package hzk

import (
	"bytes"
	"errors"
	"image/png"
	"strings"
	"testing"
)

func TestBundledFonts(t *testing.T) {
	tests := []struct {
		size   Size
		height int
	}{
		{HZK12, 12},
		{HZK16, 16},
	}

	for _, tt := range tests {
		font, err := New(tt.size)
		if err != nil {
			t.Fatalf("New(%s): %v", tt.size, err)
		}
		if font.Width() != 16 || font.Height() != tt.height {
			t.Fatalf("New(%s) metrics = %dx%d", tt.size, font.Width(), font.Height())
		}
		if !font.Contains('中') {
			t.Fatalf("New(%s) should contain 中", tt.size)
		}
	}
}

func TestGlyph(t *testing.T) {
	font, err := New(HZK16)
	if err != nil {
		t.Fatal(err)
	}

	glyph, err := font.Glyph('中')
	if err != nil {
		t.Fatal(err)
	}
	if glyph.Width != 16 || glyph.Height != 16 {
		t.Fatalf("glyph size = %dx%d", glyph.Width, glyph.Height)
	}

	on := 0
	for _, pixel := range glyph.Pixels {
		if pixel {
			on++
		}
	}
	if on == 0 {
		t.Fatal("glyph has no pixels")
	}
}

func TestMatrixAndRender(t *testing.T) {
	font, err := New(HZK12)
	if err != nil {
		t.Fatal(err)
	}

	matrix, err := font.Matrix("中文", TextOptions{GlyphSpacing: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(matrix) != 12 {
		t.Fatalf("matrix height = %d", len(matrix))
	}
	if len(matrix[0]) != 33 {
		t.Fatalf("matrix width = %d", len(matrix[0]))
	}

	rendered, err := font.Render("中", RenderOptions{Foreground: "X", Background: "."})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered, "X") || !strings.Contains(rendered, ".") {
		t.Fatalf("rendered output missing foreground/background: %q", rendered)
	}
}

func TestNewline(t *testing.T) {
	font, err := New(HZK16)
	if err != nil {
		t.Fatal(err)
	}

	bitmap, err := font.Bitmap("中\n文", TextOptions{LineSpacing: 2})
	if err != nil {
		t.Fatal(err)
	}
	if bitmap.Width != 16 || bitmap.Height != 34 {
		t.Fatalf("bitmap size = %dx%d", bitmap.Width, bitmap.Height)
	}
}

func TestCellWidth(t *testing.T) {
	font, err := New(HZK12)
	if err != nil {
		t.Fatal(err)
	}

	matrix, err := font.Matrix("中文", TextOptions{
		GlyphSpacing: 1,
		CellWidth:    12,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(matrix) != 12 {
		t.Fatalf("matrix height = %d", len(matrix))
	}
	if len(matrix[0]) != 25 {
		t.Fatalf("matrix width = %d", len(matrix[0]))
	}
}

func TestUnsupportedRune(t *testing.T) {
	font, err := New(HZK16)
	if err != nil {
		t.Fatal(err)
	}

	_, err = font.Glyph('A')
	if !errors.Is(err, ErrUnsupportedRune) {
		t.Fatalf("expected ErrUnsupportedRune, got %v", err)
	}
}

func TestEncodePNG(t *testing.T) {
	font, err := New(HZK12)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = font.EncodePNG(&buf, "中", ImageOptions{
		Scale:   2,
		Padding: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	img, err := png.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 38 || bounds.Dy() != 30 {
		t.Fatalf("image size = %dx%d", bounds.Dx(), bounds.Dy())
	}
}
