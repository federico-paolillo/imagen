// Command imagen writes a solid magenta PNG labelled with its own size.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/federico-paolillo/imagen/internal/render"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the CLI and returns the process exit code. The PNG is written
// to the current working directory and its filename is printed to stdout.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("imagen", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// -h is the height flag, so help is requested with -help or --help.
	width := fs.Int("w", 0, "image width in pixels (required)")
	height := fs.Int("h", 0, "image height in pixels (required)")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(stderr, "Usage: imagen -w WIDTH -h HEIGHT\n\n")
		_, _ = fmt.Fprintf(stderr, "Writes WIDTHxHEIGHT.png to the current directory: solid magenta with the size in black.\n")
		_, _ = fmt.Fprintf(stderr, "WIDTH and HEIGHT must be between 1 and %d.\n\n", render.MaxDimension)

		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		return 2
	}

	if fs.NArg() > 0 {
		_, _ = fmt.Fprintf(stderr, "imagen: unexpected argument %q\n", fs.Arg(0))
		fs.Usage()

		return 2
	}

	if *width == 0 || *height == 0 {
		_, _ = fmt.Fprintln(stderr, "imagen: both -w and -h are required")

		fs.Usage()

		return 2
	}

	name, err := generate(*width, *height)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "imagen: %v\n", err)

		return 1
	}

	_, _ = fmt.Fprintln(stdout, name)

	return 0
}

func generate(width, height int) (name string, err error) {
	img, err := render.Image(width, height)
	if err != nil {
		return "", err
	}

	name = render.Label(width, height) + ".png"

	file, err := os.Create(name)
	if err != nil {
		return "", fmt.Errorf("create %s: %w", name, err)
	}

	defer func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = fmt.Errorf("close %s: %w", name, cerr)
		}
	}()

	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("encode %s: %w", name, err)
	}

	return name, nil
}
