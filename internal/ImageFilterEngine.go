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

const (
	CLEAR_LINE = "\033[u\033[K"
)

type ImageFilterEngineInterface interface {
	Run(int) error
	SetFilter(string, []string) error
	GetOutput() *draw.Image
	GetOutputFilePath() (string, error)
	SetOutputFilePath(string)
	WriteOutputFile() (string, error)
}

type imageFilterEngine[T draw.Image] struct {
	filePath       string
	filter         *ImageFilterer[T]
	filterName     string
	outputFilePath string
	imgA           *T
	imgB           *T
	outputImg      *T
	wg             sync.WaitGroup
	switchBuffer   bool
}

func NewImageFilterEngine[T draw.Image](filePath, outputFilePath string, imgA, imgB T) *imageFilterEngine[T] {
	return &imageFilterEngine[T]{filePath, nil, "", outputFilePath, &imgA, &imgB, &imgB, sync.WaitGroup{}, false}
}

func (engine *imageFilterEngine[T]) Run(iterations int) error {
	if engine.filter == nil {
		return errors.New("filter not set")
	}

	currMaxProcs := runtime.GOMAXPROCS(0) * 2
	totalRows := (*engine.imgA).Bounds().Max.Y
	rowsPerProc := int(math.Ceil(float64(totalRows) / float64(currMaxProcs)))

	fmt.Println("ROWS", totalRows)
	fmt.Println("PRCS", currMaxProcs)
	fmt.Println("RPP", rowsPerProc)

	prgrsCh := make(chan int, currMaxProcs)

	for it := range iterations {
		if engine.switchBuffer {
			engine.switchOutputBuffer()
		}

		for i := 0; (i+1)*rowsPerProc <= ((*engine.imgA).Bounds().Max.Y - 1 + rowsPerProc); i++ {
			engine.wg.Add(1)
			go func(id int) {
				defer engine.wg.Done()

				//fmt.Println("STARTED GR", id, "FOR (i,j)", id*rowsPerProc, min((id+1)*rowsPerProc, totalRows))

				(*engine.filter).Apply(*engine.imgA, *engine.imgB, id*rowsPerProc, min((id+1)*rowsPerProc, totalRows), prgrsCh)
			}(i)
		}

		engine.wg.Add(1)
		go func() {
			defer engine.wg.Done()

			processedRows := 0
			for workProgressUpdate := range prgrsCh {
				processedRows += workProgressUpdate
				fmt.Print("\r")
				prgrs := (processedRows * 100) / totalRows
				fmt.Printf("PRGRS: %3d%%, IT: %d / %d", prgrs, it+1, iterations)

				if processedRows == totalRows {
					fmt.Print("\r")

					if it+1 == iterations {
						fmt.Println()
					}

					return
				}
			}
		}()
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

func (engine *imageFilterEngine[T]) GetOutputFilePath() (string, error) {
	if engine.outputFilePath != "" {
		return engine.outputFilePath, nil
	}

	if engine.filter == nil {
		return "", errors.New("filter not set")
	} else {
		return strings.TrimSuffix(engine.filePath, filepath.Ext(engine.filePath)) + "_" + engine.filterName + ".png", nil
	}
}

func (engine *imageFilterEngine[T]) SetOutputFilePath(outputFilePath string) {
	engine.outputFilePath = outputFilePath
}

func (engine *imageFilterEngine[T]) WriteOutputFile() (string, error) {
	if fileName, err := engine.GetOutputFilePath(); err != nil {
		return "", err
	} else {
		return WriteImage(fileName, engine.outputImg)
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
