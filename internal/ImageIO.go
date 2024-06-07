package internal

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func ReadImage(filepath string) (image.Image, error) {
	if file, err := os.Open(filepath); err != nil {
		return nil, err
	} else {
		defer file.Close()
		image, _, err := image.Decode(file)
		return image, err
	}
}

func WriteImage[T draw.Image](filepath string, img *T) (string, error) {
	if ofs, err := os.Create(filepath); err != nil {
		fmt.Println(err)
		return "", err
	} else {
		defer ofs.Close()
		png.Encode(ofs, *img)

		return ofs.Name(), nil
	}
}
