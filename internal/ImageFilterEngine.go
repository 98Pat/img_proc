package internal

import (
	"errors"
	"fmt"
	"image/draw"
	"math"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type ImageFilterEngineInterface interface {
	Run(int) error
	SetFilter(string, []string) error
	GetOutput() *draw.Image
	GetFileName() (string, error)
	WriteOutputFile() error
}

type imageFilterEngine[T draw.Image] struct {
	filePath     string
	filter       *ImageFilterer[T]
	filterName   string
	imgA         *T
	imgB         *T
	outputImg    *T
	wg           sync.WaitGroup
	switchBuffer bool
}

func NewImageFilterEngine[T draw.Image](filepath string, imgA, imgB T) *imageFilterEngine[T] {
	return &imageFilterEngine[T]{filepath, nil, "", &imgA, &imgB, &imgB, sync.WaitGroup{}, false}
}

func (engine *imageFilterEngine[T]) Run(iterations int) error {
	if engine.filter == nil {
		return errors.New("filter not set")
	}

	currMaxProcs := runtime.GOMAXPROCS(0)
	totalRows := (*engine.imgA).Bounds().Max.Y
	rowsPerProc := int(math.Ceil(float64(totalRows) / float64(currMaxProcs)))

	fmt.Println("ROWS", totalRows)
	fmt.Println("RPP", rowsPerProc)
	fmt.Println("PRCS", currMaxProcs)

	for range iterations {
		if engine.switchBuffer {
			engine.switchOutputBuffer()
		}

		for i := 0; (i+1)*rowsPerProc <= ((*engine.imgA).Bounds().Max.Y - 1 + rowsPerProc); i++ {
			engine.wg.Add(1)
			go func(id int) {
				defer engine.wg.Done()

				//fmt.Println("STARTED GR", id, "FOR (i,j)", id*rowsPerProc, min((id+1)*rowsPerProc, totalRows))

				(*engine.filter).Apply(*engine.imgA, *engine.imgB, id*rowsPerProc, min((id+1)*rowsPerProc, totalRows))
			}(i)
		}
		engine.wg.Wait()

		engine.switchBuffer = true
	}

	return nil
}

func (engine *imageFilterEngine[T]) switchOutputBuffer() {
	if engine.outputImg == engine.imgB {
		tmp := engine.imgA
		engine.imgA = engine.imgB
		engine.imgB = tmp

		engine.outputImg = tmp
	} else {
		tmp := engine.imgB
		engine.imgB = engine.imgA
		engine.imgA = tmp

		engine.outputImg = tmp
	}

	engine.switchBuffer = false
}

func (engine *imageFilterEngine[T]) GetOutput() *draw.Image {
	output := draw.Image(*engine.outputImg)
	return &output
}

func (engine *imageFilterEngine[T]) GetFileName() (string, error) {
	if engine.filter == nil {
		return "", errors.New("filter not set")
	} else {
		return strings.TrimSuffix(engine.filePath, filepath.Ext(engine.filePath)) + "_" + engine.filterName + ".png", nil
	}
}

func (engine *imageFilterEngine[T]) WriteOutputFile() error {
	if fileName, err := engine.GetFileName(); err != nil {
		return err
	} else {
		return WriteImageFile(fileName, *engine.outputImg)
	}
}

func (engine *imageFilterEngine[T]) SetFilter(filterName string, args []string) error {
	if tmp, err := GetFilter[T](filterName, args); err != nil {
		return err
	} else {
		engine.filter = &tmp
		engine.filterName = filterName
		return nil
	}
}
