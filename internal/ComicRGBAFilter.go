package internal

import (
	"image"
	"image/color"
	"math"
)

const (
	redFactor   float64 = 0.2126
	greenFactor float64 = 0.7152
	blueFactor  float64 = 0.0722
)

type ComicRGBA64Filter struct {
	colorStep, colorOffset uint16
	colorStepF             float64
}

type ComicRGBAFilter struct {
	colorStep, colorOffset uint8
	colorStepF             float64
}

func (filter *ComicRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			clr := uint16(math.Floor(calcIntensity(curr.Self, 16)/filter.colorStepF))*filter.colorStep + filter.colorOffset

			filteredImg.SetRGBA64(curr.X, curr.Y, color.RGBA64{clr, clr, clr, 0xffff})
		}
	}
}

func (filter *ComicRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			clr := uint8(math.Floor(calcIntensity(curr.Self, 8)/filter.colorStepF))*filter.colorStep + filter.colorOffset

			filteredImg.SetRGBA(curr.X, curr.Y, color.RGBA{clr, clr, clr, 0xff})
		}
	}
}

func calcIntensity(c *color.Color, bitSize uint32) float64 {
	var r, g, b, _ = (*c).RGBA()
	r >>= bitSize
	g >>= bitSize
	b >>= bitSize
	return float64(r)*redFactor + float64(g)*greenFactor + float64(b)*blueFactor
}
