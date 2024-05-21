package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/secriy/go-hzk"
	"github.com/spf13/cobra"
)

type commonOptions struct {
	sizeValue    string
	fontPath     string
	glyphSpacing int
	lineSpacing  int
	cellWidth    int
}

func main() {
	root := newRootCommand(os.Stdin, os.Stdout)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand(stdin *os.File, stdout io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           "hzk",
		Short:         "Render text with HZK12 or HZK16 bitmap fonts",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(newTextCommand(stdin, stdout))
	root.AddCommand(newImageCommand(stdin, stdout))
	return root
}

func addCommonFlags(cmd *cobra.Command, opts *commonOptions) {
	cmd.Flags().StringVar(&opts.sizeValue, "size", opts.sizeValue, "font size: 12 or 16")
	cmd.Flags().StringVar(&opts.fontPath, "font", opts.fontPath, "path to a raw HZK font file")
	cmd.Flags().IntVar(&opts.glyphSpacing, "spacing", opts.glyphSpacing, "blank columns between glyphs")
	cmd.Flags().IntVar(&opts.lineSpacing, "line-spacing", opts.lineSpacing, "blank rows between lines")
	cmd.Flags().IntVar(&opts.cellWidth, "cell-width", opts.cellWidth, "output cell width per glyph, 0 keeps font width")
}

func loadFont(opts commonOptions) (*hzk.Font, error) {
	size, err := hzk.ParseSize(opts.sizeValue)
	if err != nil {
		return nil, err
	}
	if opts.fontPath == "" {
		return hzk.New(size)
	}
	return hzk.LoadFile(size, opts.fontPath)
}

func readText(args []string, stdin *os.File) (string, error) {
	text := strings.Join(args, " ")
	if text != "" {
		return text, nil
	}

	stat, err := stdin.Stat()
	if err != nil {
		return "", err
	}
	if stat.Mode()&os.ModeCharDevice != 0 {
		return "", fmt.Errorf("missing text")
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\n"), nil
}
