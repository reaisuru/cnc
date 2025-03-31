package renderer

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image/gif"
	"os"
	"time"
)

type GifRenderer struct {
	source *gif.GIF
	*writeOpts
}

func GifFromPath(path string, options ...Option) (*GifRenderer, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer imageFile.Close()

	// attempt to decode image
	giff, err := gif.DecodeAll(imageFile)
	if err != nil {
		return nil, err
	}

	return NewGif(giff, options...)
}

func NewGif(source *gif.GIF, options ...Option) (*GifRenderer, error) {
	opts := &writeOpts{
		width:    100,
		height:   50,
		writer:   os.Stdout,
		drawType: TypeANSI,
	}

	for _, opt := range options {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}

	return &GifRenderer{
		source:    source,
		writeOpts: opts,
	}, nil
}

func (a *GifRenderer) Write() (int, error) {
	_, _ = a.writer.Write([]byte("\x1b[2J"))
	_, _ = a.writer.Write([]byte("\x1b[H"))

	for i := 0; i < len(a.source.Image); i++ {
		buffer := new(bytes.Buffer)
		buffer.WriteString("\x1b[0m")

		frame := imaging.Resize(a.source.Image[i], a.width, a.height, imaging.Lanczos)

		for y := 0; y < a.height; y++ {
			for x := 0; x < a.width; x++ {
				r, g, b, _ := frame.At(x*frame.Bounds().Dx()/a.width, y*frame.Bounds().Dy()/a.height).RGBA()
				buffer.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm ", r>>8, g>>8, b>>8))
			}

			if !a.full && y < a.height-1 {
				buffer.WriteString("\x1b[0m\r\n")
			} else {
				buffer.WriteString("\x1b[0m")
			}
		}

		if _, err := buffer.WriteTo(a.writer); err != nil {
			return 0, err
		}

		time.Sleep(time.Duration(a.source.Delay[i]) * time.Second / 100)
		buffer.Reset()
	}

	return 0, nil
}
