package main

import (
	"flag"
	"fmt"
	"image"
	"time"

	"bib.de/img_proc/internal"
)

var (
	helpFlag      = flag.Bool("h", false, "display flag help")
	imageFlag     = flag.String("i", "", "path to the image")
	filterFlag    = flag.String("f", "", "type of filter to be applied\nlist of filters:\n\tblur\n\tinvert")
	iterationFlag = flag.Int("I", 1, "iteration count of filter")
)

func main() {
	fmt.Println("Image Processing Collection by Patrick Protte")
	fmt.Println()

	flag.Parse()
	args := flag.Args()

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	if *imageFlag == "" {
		fmt.Println("please enter a .png file path via -i flag.\ncheck help -h for more information")
		return
	}

	if *filterFlag == "" {
		fmt.Println("please enter filter via -f flag.\ncheck help -h for more information")
		return
	}

	var programStart = time.Now()
	var start = programStart

	fmt.Println("reading image", *imageFlag)

	img, err := internal.ReadImage(*imageFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("image read")
	fmt.Printf("reading process took %d ms\n\n", time.Now().Sub(start).Milliseconds())

	start = time.Now()
	fmt.Println("starting filter process")

	var filterEngine internal.ImageFilterEngineInterface

	switch _img := img.(type) {
	case *image.RGBA64:
		tmpImg := image.NewRGBA64(img.Bounds())
		fe := internal.NewImageFilterEngine(*imageFlag, _img, tmpImg)
		filterEngine = fe
	case *image.RGBA:
		tmpImg := image.NewRGBA(img.Bounds())
		fe := internal.NewImageFilterEngine(*imageFlag, _img, tmpImg)
		filterEngine = fe
	default:
		fmt.Println("unsupported image type")
		return
	}

	if err := filterEngine.SetFilter(*filterFlag, args); err != nil {
		fmt.Println(err)
		return
	}

	if err := filterEngine.Run(*iterationFlag); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("filter finished")
	fmt.Printf("filter process took %d ms\n\n", time.Now().Sub(start).Milliseconds())

	start = time.Now()
	if fileName, err := filterEngine.GetFileName(); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("writing file", fileName)
	}

	// add filename as return of writeoutputfile, add flag for output filepath, default to inputfilepath+filtername
	filterEngine.WriteOutputFile()

	fmt.Println("done")
	fmt.Printf("writing process took %d ms\n\n", time.Now().Sub(start).Milliseconds())
	fmt.Printf("entire process took %d ms", time.Now().Sub(programStart).Milliseconds())
}
