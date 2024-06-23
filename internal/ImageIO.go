package internal

import (
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

func ReadImage(filepath string) (image.Image, error) {
	if file, err := os.Open(filepath); err != nil {
		return nil, err
	} else {
		defer file.Close()
		img, _, err := image.Decode(file)

		switch img.(type) {
		case *image.RGBA64:
		case *image.RGBA:
		case *image.YCbCr:
			rgbaImg := image.NewRGBA(img.Bounds())
			draw.Draw(rgbaImg, img.Bounds(), img, image.Point{}, draw.Src)
			img = rgbaImg
		default:
			return nil, errors.New("unsupported image type")
		}

		return img, err
	}
}

func WriteImage[T draw.Image](filepath string, img *T) (string, error) {
	if ofs, err := os.Create(filepath); err != nil {
		return "", err
	} else {
		defer ofs.Close()
		png.Encode(ofs, *img)

		return ofs.Name(), nil
	}
}
