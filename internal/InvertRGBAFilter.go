package internal

import (
	"image"
	"image/color"
)

type InvertRGBA64Filter struct{}

type InvertRGBAFilter struct{}

func (filter *InvertRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			var r, g, b, a uint32 = (*curr.Self).RGBA()
			filteredImg.SetRGBA64(curr.X, curr.Y, color.RGBA64{uint16(0xffff - r), uint16(0xffff - g), uint16(0xffff - b), uint16(a)})
		}
	}
}

func (filter *InvertRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			var r, g, b, a uint32 = (*curr.Self).RGBA()
			filteredImg.SetRGBA(curr.X, curr.Y, color.RGBA{uint8(0xff - r), uint8(0xff - g), uint8(0xff - b), uint8(a)})
		}
	}
}
