package internal

import (
	"image"
	"image/color"
	"math"
)

type SpotRGBA64Filter struct {
	spotX, spotY int
	spotR        float64
}

type SpotRGBAFilter struct {
	spotX, spotY int
	spotR        float64
}

func (filter *SpotRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			dX, dY := curr.X-filter.spotX, curr.Y-filter.spotY
			d := math.Sqrt(float64(dX*dX + dY*dY))

			var r, g, b, a uint32 = (*curr.Self).RGBA()
			spotAdjust(&r, &g, &b, &d, &filter.spotR, 0xffff)

			filteredImg.SetRGBA64(curr.X, curr.Y, color.RGBA64{
				uint16(r >> 16),
				uint16(g >> 16),
				uint16(b >> 16),
				uint16(a >> 16),
			})

		}
	}
}

func (filter *SpotRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			dX, dY := filter.spotX-curr.X, filter.spotY-curr.Y
			d := math.Sqrt(float64(dX*dX + dY*dY))

			var r, g, b, a uint32 = (*curr.Self).RGBA()
			spotAdjust(&r, &g, &b, &d, &filter.spotR, 0xff)

			filteredImg.SetRGBA(curr.X, curr.Y, color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			})
		}
	}
}

func spotAdjust(r, g, b *uint32, d *float64, rad *float64, cMax uint32) {
	if *d > *rad {
		*r = cMax
		*g = cMax
		*b = cMax
	} else {
		fac := 1 - (*d / *rad)
		*r = uint32(float64(*r) * fac)
		*g = uint32(float64(*g) * fac)
		*b = uint32(float64(*b) * fac)
	}
}
