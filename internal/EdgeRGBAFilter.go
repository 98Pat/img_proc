package internal

import (
	"image"
	"image/color"
)

type EdgeRGBA64Filter struct {
	amp int64
}

type EdgeRGBAFilter struct {
	amp int64
}

func (filter *EdgeRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, DIRECT, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			var iv, ih int64
			addIntensity(&iv, curr.North, 16)
			subIntensity(&iv, curr.South, 16)
			if iv < 0 {
				iv = -iv
			}

			addIntensity(&ih, curr.West, 16)
			subIntensity(&ih, curr.East, 16)
			if ih < 0 {
				ih = -ih
			}

			iv = min((iv+ih)*filter.amp, 0xffff)

			filteredImg.SetRGBA64(curr.X, curr.Y, color.RGBA64{
				uint16(iv >> 16),
				uint16(iv >> 16),
				uint16(iv >> 16),
				0xffff,
			})
		}
	}
}

func (filter *EdgeRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, DIRECT, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			var iv, ih int64
			addIntensity(&iv, curr.North, 8)
			subIntensity(&iv, curr.South, 8)
			if iv < 0 {
				iv = -iv
			}

			addIntensity(&ih, curr.West, 8)
			subIntensity(&ih, curr.East, 8)
			if ih < 0 {
				ih = -ih
			}

			iv = min((iv+ih)*filter.amp, 0xff)

			filteredImg.SetRGBA(curr.X, curr.Y, color.RGBA{
				uint8(iv),
				uint8(iv),
				uint8(iv),
				0xff,
			})
		}
	}
}

func addIntensity(intensity *int64, c *color.Color, bitSize uint32) {
	if c != nil {
		*intensity += int64(calcIntensity(c, bitSize))
	}
}

func subIntensity(intensity *int64, c *color.Color, bitSize uint32) {
	if c != nil {
		*intensity -= int64(calcIntensity(c, bitSize))
	}
}
