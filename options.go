package hzk

import "image/color"

// TextOptions controls how text is converted to a bitmap matrix.
type TextOptions struct {
	// GlyphSpacing is the number of blank columns inserted between glyphs.
	GlyphSpacing int
	// LineSpacing is the number of blank rows inserted between text lines.
	LineSpacing int
	// CellWidth is the output width of each glyph cell.
	// A value of 0 keeps the font's native glyph width. A positive value crops
	// wider glyphs from the right or pads narrower glyphs on the right.
	CellWidth int
	// DisableWiden disables the default conversion from printable ASCII
	// characters to their full-width forms before glyph lookup.
	DisableWiden bool
}

// RenderOptions controls how text is rendered to a string.
type RenderOptions struct {
	// Foreground is written for enabled pixels. Empty value defaults to "#".
	Foreground string
	// Background is written for disabled pixels. Empty value defaults to " ".
	Background string
	// GlyphSpacing is the number of blank columns inserted between glyphs.
	GlyphSpacing int
	// LineSpacing is the number of blank rows inserted between text lines.
	LineSpacing int
	// CellWidth is the output width of each glyph cell.
	// A value of 0 keeps the font's native glyph width. A positive value crops
	// wider glyphs from the right or pads narrower glyphs on the right.
	CellWidth int
	// DisableWiden disables the default conversion from printable ASCII
	// characters to their full-width forms before glyph lookup.
	DisableWiden bool
}

// ImageOptions controls how text or bitmaps are encoded as images.
type ImageOptions struct {
	// Foreground is used for enabled pixels. Nil defaults to color.Black.
	Foreground color.Color
	// Background is used for disabled pixels and padding. Nil defaults to color.White.
	Background color.Color
	// Scale is the output pixel size of each bitmap dot. Values smaller than 1 default to 1.
	Scale int
	// Padding is the number of pixels added around the rendered bitmap. Negative values default to 0.
	Padding int
	// GlyphSpacing is the number of blank columns inserted between glyphs.
	GlyphSpacing int
	// LineSpacing is the number of blank rows inserted between text lines.
	LineSpacing int
	// CellWidth is the output width of each glyph cell before image scaling.
	// A value of 0 keeps the font's native glyph width.
	CellWidth int
	// DisableWiden disables the default conversion from printable ASCII
	// characters to their full-width forms before glyph lookup.
	DisableWiden bool
}

func mergeTextOptions(options []TextOptions) TextOptions {
	if len(options) == 0 {
		return TextOptions{}
	}
	return options[0]
}

func mergeRenderOptions(options []RenderOptions) RenderOptions {
	opt := RenderOptions{
		Foreground: "#",
		Background: " ",
	}
	if len(options) > 0 {
		opt = options[0]
		if opt.Foreground == "" {
			opt.Foreground = "#"
		}
		if opt.Background == "" {
			opt.Background = " "
		}
	}
	return opt
}

func mergeImageOptions(options []ImageOptions) ImageOptions {
	opt := ImageOptions{
		Foreground: color.Black,
		Background: color.White,
		Scale:      1,
	}
	if len(options) > 0 {
		opt = options[0]
		if opt.Foreground == nil {
			opt.Foreground = color.Black
		}
		if opt.Background == nil {
			opt.Background = color.White
		}
		if opt.Scale < 1 {
			opt.Scale = 1
		}
		if opt.Padding < 0 {
			opt.Padding = 0
		}
	}
	return opt
}
