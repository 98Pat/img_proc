package internal

import (
	"image"
	"image/color"
	"math"
)

const (
	HEAT_COLOR_STEP     = 42
	HEAT_COLOR_ARR_SIZE = 6
)

var heatColorArr = [HEAT_COLOR_ARR_SIZE]byte{
	0b000,
	0b001,
	0b011,
	0b010,
	0b110,
	0b100,
}

type HeatRGBA64Filter struct {
}

type HeatRGBAFilter struct {
}

func (filter *HeatRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			clrByte := heatColorArr[min(HEAT_COLOR_ARR_SIZE-1, int(math.Floor(calcIntensity(curr.Self, 16)/HEAT_COLOR_STEP)))]

			r := uint16(clrByte&0b100) * 0xffff
			g := uint16(clrByte&0b010) * 0xffff
			b := uint16(clrByte&0b001) * 0xffff

			filteredImg.SetRGBA64(curr.X, curr.Y, color.RGBA64{r, g, b, 0xffff})
		}
	}
}

func (filter *HeatRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			clrByte := heatColorArr[min(HEAT_COLOR_ARR_SIZE-1, int(math.Floor(calcIntensity(curr.Self, 8)/HEAT_COLOR_STEP)))]

			r := uint8(clrByte&0b100) * 0xff
			g := uint8(clrByte&0b010) * 0xff
			b := uint8(clrByte&0b001) * 0xff

			filteredImg.SetRGBA(curr.X, curr.Y, color.RGBA{r, g, b, 0xff})
		}
	}
}
