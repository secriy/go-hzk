package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/secriy/go-hzk"
	"github.com/spf13/cobra"
)

type imageOptions struct {
	commonOptions
	outPath string
	scale   int
	padding int
	fgColor string
	bgColor string
}

func newImageCommand(stdin *os.File, stdout io.Writer) *cobra.Command {
	opts := imageOptions{
		commonOptions: commonOptions{
			sizeValue:    "16",
			glyphSpacing: 1,
		},
		scale:   8,
		padding: 8,
		fgColor: "#000000",
		bgColor: "#ffffff",
	}

	cmd := &cobra.Command{
		Use:   "image [text]",
		Short: "Render text to a PNG image",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.outPath == "" {
				return fmt.Errorf("required flag \"out\" not set")
			}
			font, err := loadFont(opts.commonOptions)
			if err != nil {
				return err
			}
			text, err := readText(args, stdin)
			if err != nil {
				return err
			}
			return writePNG(font, text, opts, stdout)
		},
	}

	addCommonFlags(cmd, &opts.commonOptions)
	cmd.Flags().StringVarP(&opts.outPath, "out", "o", opts.outPath, "PNG output path, or - for stdout")
	cmd.Flags().IntVar(&opts.scale, "scale", opts.scale, "PNG pixels per dot")
	cmd.Flags().IntVar(&opts.padding, "padding", opts.padding, "PNG padding in pixels")
	cmd.Flags().StringVar(&opts.fgColor, "fg-color", opts.fgColor, "PNG foreground color")
	cmd.Flags().StringVar(&opts.bgColor, "bg-color", opts.bgColor, "PNG background color")
	return cmd
}

func writePNG(font *hzk.Font, text string, opts imageOptions, stdout io.Writer) error {
	fg, err := parseHexColor(opts.fgColor)
	if err != nil {
		return fmt.Errorf("invalid foreground color: %w", err)
	}
	bg, err := parseHexColor(opts.bgColor)
	if err != nil {
		return fmt.Errorf("invalid background color: %w", err)
	}

	if opts.outPath == "-" {
		return font.EncodePNG(stdout, text, hzk.ImageOptions{
			Foreground:   fg,
			Background:   bg,
			Scale:        opts.scale,
			Padding:      opts.padding,
			GlyphSpacing: opts.glyphSpacing,
			LineSpacing:  opts.lineSpacing,
			CellWidth:    opts.cellWidth,
		})
	}

	file, err := os.Create(opts.outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return font.EncodePNG(file, text, hzk.ImageOptions{
		Foreground:   fg,
		Background:   bg,
		Scale:        opts.scale,
		Padding:      opts.padding,
		GlyphSpacing: opts.glyphSpacing,
		LineSpacing:  opts.lineSpacing,
		CellWidth:    opts.cellWidth,
	})
}

func parseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 && len(s) != 8 {
		return color.RGBA{}, fmt.Errorf("expected #RRGGBB or #RRGGBBAA")
	}

	value, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}
	if len(s) == 6 {
		return color.RGBA{
			R: uint8(value >> 16),
			G: uint8(value >> 8),
			B: uint8(value),
			A: 0xff,
		}, nil
	}
	return color.RGBA{
		R: uint8(value >> 24),
		G: uint8(value >> 16),
		B: uint8(value >> 8),
		A: uint8(value),
	}, nil
}
