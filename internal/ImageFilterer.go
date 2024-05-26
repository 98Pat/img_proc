package internal

import (
	"image/draw"
)

type ImageFilterer[T draw.Image] interface {
	Apply(T, T, int, int)
}
