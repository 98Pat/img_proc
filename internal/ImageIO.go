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

func WriteImageFile[T draw.Image](filepath string, img *T) error {
	if ofs, err := os.Create(filepath); err != nil {
		fmt.Println(err)
		return err
	} else {
		defer ofs.Close()
		png.Encode(ofs, *img)
	}

	return nil
}
