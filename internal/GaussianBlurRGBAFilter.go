package internal

import (
	"image"
	"image/color"
	"math"
)

type GaussianBlurRGBA64Filter struct {
	kernelSize int
	kernel     [][]float64
}

func (filter *GaussianBlurRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			values := getKernelValues(img, filter.kernelSize, curr.X, curr.Y)
		}
	}
}

func applyKernelToValues(values [][]color.Color, kernel [][]float64, kernelSize int) {
	var r, g, b, a float64 = 0, 0, 0, 0
	for x := range kernelSize {
		for y := range kernelSize {
			if clr := values[x][y]; clr != nil {
				vr, vg, vb, va := clr.RGBA()
				kVal := kernel[x][y]
				r += float64(vr) * kVal
				g += float64(vg) * kVal
				b += float64(vb) * kVal
				a += float64(va) * kVal
			} else {
				// mirroring? check kernel image proc wiki on convolution
				clr := values[]
			}
		}
	}
}

func getKernelValues[T image.Image](img T, kernelSize int, x int, y int) [][]color.Color {
	var kernelValues [][]color.Color = make([][]color.Color, kernelSize)
	for k := range kernelSize {
		kernelValues[k] = make([]color.Color, kernelSize)
	}

	bnds := img.Bounds()
	kOffset := int(math.Floor(float64(kernelSize) / 2.0))
	kOffsetX := x - kOffset
	kOffsetY := y - kOffset
	for kX := kOffsetX; kX < kX+kernelSize; kX++ {
		for kY := kOffsetY; kY < kY+kernelSize; kY++ {
			if kX < bnds.Min.X || kX >= bnds.Max.X {
				continue
			}
			if kY < bnds.Min.Y || kY >= bnds.Max.Y {
				continue
			}

			kernelValues[kX-kOffsetX][kY-kOffsetY] = img.At(kX-kOffsetX, kY-kOffsetY)
		}
	}

	return kernelValues
}

func buildKernel(kernelSize int) [][]float64 {
	const sigma float64 = 1
	var kernel [][]float64 = make([][]float64, kernelSize)
	for k := range kernelSize {
		kernel[k] = make([]float64, kernelSize)
	}

	mean := float64(kernelSize) / 2.0
	sum := 0.0

	for x := 0; x < kernelSize; x++ {
		for y := 0; y < kernelSize; y++ {
			kernel[x][y] = math.Exp(-0.5*(math.Pow((float64(x)-mean)/sigma, 2.0)+math.Pow((float64(y)-mean)/sigma, 2.0))) / (2 * math.Pi * sigma * sigma)

			sum += kernel[x][y]
		}
	}

	for x := 0; x < kernelSize; x++ {
		for y := 0; y < kernelSize; y++ {
			kernel[x][y] /= sum
		}
	}

	return kernel
}
