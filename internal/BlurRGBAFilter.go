package internal

import (
	"image"
	"image/color"
)

type BlurRGBA64Filter struct{}

type BlurRGBAFilter struct{}

func (filter *BlurRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan uint8) {
	if iter, err := NewImageIterator(img, DIRECT, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()
			var r, g, b, a uint32 = (*curr.Self).RGBA()
			var i uint32 = 1

			addRGBAI(&r, &g, &b, &a, &i, curr.North)
			addRGBAI(&r, &g, &b, &a, &i, curr.West)
			addRGBAI(&r, &g, &b, &a, &i, curr.East)
			addRGBAI(&r, &g, &b, &a, &i, curr.South)

			r, g, b, a = r/i, g/i, b/i, a/i

			filteredImg.SetRGBA64(curr.X, curr.Y, color.RGBA64{
				uint16(r >> 16),
				uint16(g >> 16),
				uint16(b >> 16),
				uint16(a >> 16),
			})
		}
	}
}

func (filter *BlurRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int, prgrsCh chan uint8) {
	if iter, err := NewImageIterator(img, DIRECT, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()
			var r, g, b, a uint32 = (*curr.Self).RGBA()
			var i uint32 = 1

			addRGBAI(&r, &g, &b, &a, &i, curr.North)
			addRGBAI(&r, &g, &b, &a, &i, curr.West)
			addRGBAI(&r, &g, &b, &a, &i, curr.East)
			addRGBAI(&r, &g, &b, &a, &i, curr.South)

			r, g, b, a = r/i, g/i, b/i, a/i

			filteredImg.SetRGBA(curr.X, curr.Y, color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			})
		}
	}
}

func addRGBAI(r, g, b, a, i *uint32, c *color.Color) {
	if c != nil {
		tr, tg, tb, ta := (*c).RGBA()
		*r, *g, *b, *a = *r+tr, *g+tg, *b+tb, *a+ta
		*i++
	}
}
