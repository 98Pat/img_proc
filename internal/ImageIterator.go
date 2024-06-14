package internal

import (
	"image"
	"image/color"
	"math"
)

type ImageIteratorNeighbourCount int

const (
	NONE   ImageIteratorNeighbourCount = 0
	DIRECT ImageIteratorNeighbourCount = 4

	WORK_PROGRESS_STEP_MULT int = 2
)

type ImageIterator interface {
	HasNext() bool
	Next() *imageIteratorYield
}

type imageIterator[T image.Image] struct {
	img                      T
	curX, curY, startY, endY int
	rowBufferNorth           []color.Color
	rowBufferSouth           []color.Color
	neighbourCount           ImageIteratorNeighbourCount
	current                  imageIteratorYield
	prgrsCh                  chan int
	workProgressStep         int
}

type imageIteratorYield struct {
	self, north, east, south, west color.Color

	Self, North, East, South, West *color.Color

	X, Y int
}

func NewImageIterator[T image.Image](img T, neighbourCount ImageIteratorNeighbourCount, startY, endY int, prgrsCh chan int) (*imageIterator[T], error) {
	var rowBufferNorth, rowBufferSouth []color.Color = nil, nil

	rowBufferNorth = make([]color.Color, img.Bounds().Max.X)
	rowBufferSouth = make([]color.Color, img.Bounds().Max.X)

	workProgressStep := int(math.Max(float64((endY-startY)/WORK_PROGRESS_STEP_MULT), 1))

	return &imageIterator[T]{img, img.Bounds().Min.X, startY, startY, endY, rowBufferNorth, rowBufferSouth, neighbourCount, imageIteratorYield{}, prgrsCh, workProgressStep}, nil
}

func (iter *imageIterator[T]) HasNext() bool {
	if iter.curY >= iter.endY {
		return false
	}

	return true
}

func (iter *imageIterator[T]) Next() *imageIteratorYield {
	iter.current.X = iter.curX
	iter.current.Y = iter.curY

	if iter.neighbourCount == NONE {
		iter.current.self = iter.img.At(iter.curX, iter.curY)
	} else {
		if iter.curX == 0 {
			iter.current.self = iter.img.At(iter.curX, iter.curY)
			iter.current.West = nil
			iter.current.east = iter.img.At(iter.curX+1, iter.curY)
			iter.current.East = &iter.current.east
		} else {
			iter.current.west = iter.current.self
			iter.current.West = &iter.current.west
			iter.current.self = iter.current.east

			if iter.curX+1 == iter.img.Bounds().Max.X {
				iter.current.East = nil
			} else {
				if iter.curY != iter.startY {
					iter.current.east = iter.rowBufferSouth[iter.curX+1]
				} else {
					iter.current.east = iter.img.At(iter.curX+1, iter.curY)
				}
				iter.current.East = &iter.current.east
			}
		}

		if iter.curY == iter.startY {
			if iter.curY == 0 {
				iter.current.North = nil
			} else {
				iter.current.north = iter.img.At(iter.curX, iter.curY-1)
				iter.current.North = &iter.current.north
			}

			iter.current.south = iter.img.At(iter.curX, iter.curY+1)
			iter.current.South = &iter.current.south
		} else {
			iter.current.north = iter.rowBufferNorth[iter.curX]
			iter.current.North = &iter.current.north

			if iter.curY+1 == iter.img.Bounds().Max.Y {
				iter.current.South = nil
			} else {
				iter.current.south = iter.img.At(iter.curX, iter.curY+1)
				iter.current.South = &iter.current.south
			}
		}

		iter.rowBufferNorth[iter.curX] = iter.current.self
		iter.rowBufferSouth[iter.curX] = iter.current.south
	}

	iter.current.Self = &iter.current.self

	iter.curX++

	if iter.curX >= iter.img.Bounds().Max.X {
		iter.curY++

		if (iter.curY-iter.startY)%iter.workProgressStep == 0 {
			iter.prgrsCh <- iter.workProgressStep
		} else if iter.curY == iter.endY {
			iter.prgrsCh <- (iter.endY - iter.startY) % iter.workProgressStep
		}

		iter.curX = 0
	}

	return &iter.current
}
