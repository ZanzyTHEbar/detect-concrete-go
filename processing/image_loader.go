package processing

import (
	"fmt"

	"gocv.io/x/gocv"
)

type ImageLoader interface {
	LoadImage(path string) (gocv.Mat, error)
}

type ConcreteImageLoader struct{}

func (loader ConcreteImageLoader) LoadImage(path string) (gocv.Mat, error) {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		return img, fmt.Errorf("Error reading image from: %v", path)
	}
	return img, nil
}

func NewImageLoader() ImageLoader {
	return &ConcreteImageLoader{}
}
