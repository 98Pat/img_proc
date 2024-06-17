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
			r, g, b, a := applyKernelToValues(values, filter.kernel, filter.kernelSize)
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
			r, g, b, a := applyKernelToValues(values, filter.kernel, filter.kernelSize)
			clr := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			filteredImg.SetRGBA(curr.X, curr.Y, clr)
		}
	}
}

func applyKernelToValues(values [][]color.Color, kernel [][]float64, kernelSize int) (r, g, b, a float64) {
	for x := range kernelSize {
		for y := range kernelSize {
			if clr := values[x][y]; clr != nil {
				vr, vg, vb, va := clr.RGBA()
				kVal := kernel[x][y]
				r += float64(vr) * kVal
				g += float64(vg) * kVal
				b += float64(vb) * kVal
				a += float64(va) * kVal
			}
		}
	}

	return
}

func getKernelValues[T image.Image](img T, kernelSize int, x int, y int) [][]color.Color {
	var kernelValues [][]color.Color = make([][]color.Color, kernelSize)
	for k := range kernelSize {
		kernelValues[k] = make([]color.Color, kernelSize)
	}

	//fmt.Printf("CALL X %d CALL Y %d\n", x, y)
	bnds := img.Bounds()
	kOffset := int(math.Floor(float64(kernelSize) / 2.0))
	kOffsetX := x - kOffset
	kOffsetY := y - kOffset
	for kY := -kOffset; kY < kOffset; kY++ {
		for kX := -kOffset; kX < kOffset; kX++ {
			kIdxX := kX + kOffsetX
			kIdxY := kY + kOffsetY

			if kIdxX < bnds.Min.X || kIdxX >= bnds.Max.X {
				kIdxX = (2*bnds.Max.X - (kIdxX)) % bnds.Max.X
			}
			if kIdxY < bnds.Min.Y || kIdxY >= bnds.Max.Y {
				kIdxY = (2*bnds.Max.Y - (kIdxY)) % bnds.Max.Y
			}

			//fmt.Printf("X %d Y %d\n", kIdxX, kIdxY)
			kernelValues[kX+kOffset][kY+kOffset] = img.At(kIdxX, kIdxY)
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
