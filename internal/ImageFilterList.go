package internal

import (
	"errors"
	"image"
	"image/draw"
	"strconv"
)

type FilterConstructor func(args []string) (interface{}, error)

var rgba64FilterConstructors = map[string]FilterConstructor{
	"invert": func(args []string) (interface{}, error) {
		return &InvertRGBA64Filter{}, nil
	},
	"blur": func(args []string) (interface{}, error) {
		return &BlurRGBA64Filter{}, nil
	},
	"comic": func(args []string) (interface{}, error) {
		var colorStep uint16
		if len(args) >= 1 {
			ui32, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return nil, errors.New("filter requires color-step count (uint32) as first non-flag argument")
			}
			colorStep = uint16(ui32)
		} else {
			colorStep = 0xffff / 3
		}
		return &ComicRGBA64Filter{colorStep, colorStep / 2, float64(colorStep)}, nil
	},
	"spot": func(args []string) (interface{}, error) {
		if len(args) < 3 {
			return nil, errors.New("filter needs x, y (int) coordinates and a radius (float) as non-flag arguments")
		}
		spotX, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		spotY, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		spotR, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return nil, err
		}
		return &SpotRGBA64Filter{spotX, spotY, spotR}, nil
	},
	"edge": func(args []string) (interface{}, error) {
		var amp int64

		if len(args) >= 1 {
			if a, err := strconv.ParseInt(args[0], 10, 64); err != nil {
				return nil, errors.New("first non-flag arguments needs to be an amplification modifier for this filter (int)")
			} else {
				amp = a
			}
		} else {
			amp = 1
		}
		return &EdgeRGBA64Filter{amp}, nil
	},
}

var rgbaFilterConstructors = map[string]FilterConstructor{
	"invert": func(args []string) (interface{}, error) {
		return &InvertRGBAFilter{}, nil
	},
	"blur": func(args []string) (interface{}, error) {
		return &BlurRGBAFilter{}, nil
	},
	"comic": func(args []string) (interface{}, error) {
		var colorStep uint8
		if len(args) >= 1 {
			ui32, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return nil, errors.New("filter requires color-step count (uint32) as first non-flag argument")
			}
			colorStep = 0xff / uint8(ui32)
		} else {
			colorStep = 0xff / 3
		}
		return &ComicRGBAFilter{colorStep, colorStep / 2, float64(colorStep)}, nil
	},
	"spot": func(args []string) (interface{}, error) {
		if len(args) < 3 {
			return nil, errors.New("filter needs x, y (int) coordinates and a radius (float) as non-flag arguments")
		}
		spotX, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		spotY, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		spotR, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return nil, err
		}
		return &SpotRGBAFilter{spotX, spotY, spotR}, nil
	},
	"edge": func(args []string) (interface{}, error) {
		var amp int64

		if len(args) >= 1 {
			if a, err := strconv.ParseInt(args[0], 10, 64); err != nil {
				return nil, errors.New("first non-flag arguments needs to be an amplification modifier for this filter (int)")
			} else {
				amp = a
			}
		} else {
			amp = 1
		}
		return &EdgeRGBAFilter{amp}, nil
	},
}

func GetFilter[T draw.Image](filterName string, args []string) (ImageFilterer[T], error) {
	var img T
	var constructor FilterConstructor
	var found bool

	switch any(img).(type) {
	case *image.RGBA64:
		constructor, found = rgba64FilterConstructors[filterName]
	case *image.RGBA:
		constructor, found = rgbaFilterConstructors[filterName]
	default:
		return nil, errors.New("unsupported image type")
	}

	if !found {
		return nil, errors.New("unknown filter type")
	}

	filter, err := constructor(args)
	if err != nil {
		return nil, err
	}
	if tf, ok := filter.(ImageFilterer[T]); ok {
		return tf, nil
	}

	return nil, errors.New("filter type mismatch")
}
