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

type GaussianBlurRGBAFilter struct {
	kernelSize int
	kernel     [][]float64
}

func (filter *GaussianBlurRGBA64Filter) Apply(img, filteredImg *image.RGBA64, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			values := getKernelValues(img, filter.kernelSize, curr.X, curr.Y)
			r, g, b, a := applyKernelToValues(values, filter.kernel, filter.kernelSize, 16)
			clr := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			filteredImg.SetRGBA64(curr.X, curr.Y, clr)
		}
	}
}

func (filter *GaussianBlurRGBAFilter) Apply(img, filteredImg *image.RGBA, startY, endY int, prgrsCh chan int) {
	if iter, err := NewImageIterator(img, NONE, startY, endY, prgrsCh); err == nil {
		for iter.HasNext() {
			curr := iter.Next()

			values := getKernelValues(img, filter.kernelSize, curr.X, curr.Y)
			r, g, b, a := applyKernelToValues(values, filter.kernel, filter.kernelSize, 8)
			clr := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			filteredImg.SetRGBA(curr.X, curr.Y, clr)
		}
	}
}

func applyKernelToValues(values [][]color.Color, kernel [][]float64, kernelSize int, bitSize uint32) (r, g, b, a float64) {
	for x := range kernelSize {
		for y := range kernelSize {
			if clr := values[x][y]; clr != nil {
				vr, vg, vb, va := clr.RGBA()
				vr >>= bitSize
				vg >>= bitSize
				vb >>= bitSize
				va >>= bitSize

				kFac := kernel[x][y]
				r += float64(vr) * kFac
				g += float64(vg) * kFac
				b += float64(vb) * kFac
				a += float64(va) * kFac
			}
		}
	}

	maxValue := 1 << bitSize
	r = math.Min(r, float64(maxValue))
	g = math.Min(g, float64(maxValue))
	b = math.Min(b, float64(maxValue))
	a = math.Min(a, float64(maxValue))

	return
}

func getKernelValues[T image.Image](img T, kernelSize int, x int, y int) [][]color.Color {
	kernelValues := make([][]color.Color, kernelSize)
	for k := 0; k < kernelSize; k++ {
		kernelValues[k] = make([]color.Color, kernelSize)
	}

	bnds := img.Bounds()
	kOffset := int(math.Floor(float64(kernelSize) / 2.0))
	for kY := 0; kY < kernelSize; kY++ {
		for kX := 0; kX < kernelSize; kX++ {
			kIdxX := x - kOffset + kX
			kIdxY := y - kOffset + kY

			if kIdxX < bnds.Min.X {
				kIdxX = -kIdxX
			} else if kIdxX >= bnds.Max.X {
				kIdxX = 2*bnds.Max.X - kIdxX - 1
			}

			if kIdxY < bnds.Min.Y {
				kIdxY = -kIdxY
			} else if kIdxY >= bnds.Max.Y {
				kIdxY = 2*bnds.Max.Y - kIdxY - 1
			}

			kernelValues[kX][kY] = img.At(kIdxX, kIdxY)
		}
	}

	return kernelValues
}

func buildKernel(kernelSize int, sigma float64) [][]float64 {
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
