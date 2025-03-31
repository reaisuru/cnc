package renderer

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

type ImageRenderer struct {
	source image.Image
	*writeOpts
}

func ImageFromPath(path string, filter imaging.ResampleFilter, options ...Option) (*ImageRenderer, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer imageFile.Close()

	// attempt to decode image
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	return NewImage(img, filter, options...)
}

func NewImage(source image.Image, filter imaging.ResampleFilter, options ...Option) (*ImageRenderer, error) {
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

	img := imaging.Resize(source, opts.width, opts.height, filter)

	return &ImageRenderer{
		source:    img,
		writeOpts: opts,
	}, nil
}

func (a *ImageRenderer) Write() (int, error) {
	buffer := new(bytes.Buffer)
	buffer.WriteString("\x1b[0m")

	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			r, g, b, alpha := a.source.At(x*a.source.Bounds().Dx()/a.width, y*a.source.Bounds().Dy()/a.height).RGBA()

			if alpha == 0 {
				buffer.WriteString("\x1b[0m ")
			} else {
				rA, gA, bA := r>>8, g>>8, b>>8
				buffer.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm ", rA, gA, bA))
			}
		}

		if !a.full && y < a.height-1 {
			buffer.WriteString("\x1b[0m\r\n")
		} else {
			buffer.WriteString("\x1b[0m")
		}
	}

	return a.writer.Write(buffer.Bytes())
}
