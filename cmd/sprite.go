package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

type SpriteMeasurements struct {
	Width, Height int
	FrameCount    int
	Horizontal    bool
}

var (
	DroidZapperWakeMeasurements = SpriteMeasurements{
		FrameCount: 6,
		Width:      58,
		Height:     41,
		Horizontal: false,
	}
)

func NewSpriteFromMeasurements(p pixel.Picture, d SpriteMeasurements, n int) (*pixel.Sprite, error) {
	var minx, miny, maxy, maxx float64
	if n >= d.FrameCount {
		return nil, fmt.Errorf("requested sprite frame at index %d but image only has %d frames", n, d.FrameCount)
	}
	if d.Horizontal {
		minx = float64(n * d.Width)
		maxx = float64((n + 1) * d.Width)
		miny = 0
		maxy = float64(d.Height)
	} else {
		minx = 0
		maxx = float64(d.Width)
		miny = float64(n * d.Height)
		maxy = float64((n + 1) * d.Height)
	}
	return pixel.NewSprite(p, pixel.R(minx, miny, maxx, maxy)), nil //pic.Bounds())
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
