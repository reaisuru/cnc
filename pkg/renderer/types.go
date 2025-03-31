package renderer

import "io"

type writeOpts struct {
	width  int
	height int
	writer io.Writer

	// full will full screen
	full bool

	// drawType is the draw type, braille or renderer
	drawType int
}

type Renderer interface {
	Write() (int, error)
}

type Option func(*writeOpts) error

const (
	TypeBraille = iota // braille will not work as of rn
	TypeANSI
)

func Writer(w io.Writer) Option {
	return func(opts *writeOpts) error {
		opts.writer = w
		return nil
	}
}

func Width(w int) Option {
	return func(opts *writeOpts) error {
		opts.width = w
		return nil
	}
}

func Height(h int) Option {
	return func(opts *writeOpts) error {
		opts.height = h
		return nil
	}
}

func Type(t int) Option {
	return func(opts *writeOpts) error {
		opts.drawType = t
		return nil
	}
}

func Full() Option {
	return func(opts *writeOpts) error {
		opts.full = true
		return nil
	}
}
