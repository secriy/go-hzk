package main

import (
	"fmt"
	"io"
	"os"

	"github.com/secriy/go-hzk"
	"github.com/spf13/cobra"
)

type textOptions struct {
	commonOptions
	foreground string
	background string
}

func newTextCommand(stdin *os.File, stdout io.Writer) *cobra.Command {
	opts := textOptions{
		commonOptions: commonOptions{
			sizeValue:    "16",
			glyphSpacing: 1,
		},
		foreground: "#",
		background: " ",
	}

	cmd := &cobra.Command{
		Use:   "text [text]",
		Short: "Render text to terminal output",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			font, err := loadFont(opts.commonOptions)
			if err != nil {
				return err
			}
			text, err := readText(args, stdin)
			if err != nil {
				return err
			}
			rendered, err := font.Render(text, hzk.RenderOptions{
				Foreground:   opts.foreground,
				Background:   opts.background,
				GlyphSpacing: opts.glyphSpacing,
				LineSpacing:  opts.lineSpacing,
				CellWidth:    opts.cellWidth,
				DisableWiden: opts.disableWiden,
			})
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(stdout, rendered)
			return err
		},
	}

	addCommonFlags(cmd, &opts.commonOptions)
	cmd.Flags().StringVar(&opts.foreground, "fg", opts.foreground, "foreground character")
	cmd.Flags().StringVar(&opts.background, "bg", opts.background, "background character")
	return cmd
}
