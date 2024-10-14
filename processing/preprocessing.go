package processing

import (
	"fmt"

	"github.com/ZanzyTHEbar/concrete-damage-detection/global"
	"gocv.io/x/gocv"
)

// PreprocessStrategy defines the preprocessing strategy interface.
type PreprocessStrategy interface {
	Process(img gocv.Mat) gocv.Mat
}

// GrayscaleStrategy converts an image to grayscale.
type GrayscaleStrategy struct{}

func (gs *GrayscaleStrategy) Process(img gocv.Mat) gocv.Mat {
	grayImg := gocv.NewMat()
	//defer grayImg.Close()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)
	return grayImg
}

// EdgeDetectionStrategy applies Canny edge detection.
type EdgeDetectionStrategy struct {
	Threshold1, Threshold2 float32
}

func (es *EdgeDetectionStrategy) Process(img gocv.Mat) gocv.Mat {
	edges := gocv.NewMat()
	//defer edges.Close()

	gocv.Canny(img, &edges, es.Threshold1, es.Threshold2)
	fmt.Printf("Canny Edge Detection: Threshold1=%.2f, Threshold2=%.2f\n", es.Threshold1, es.Threshold2)
	return edges
}

// AutoTuneEdgeDetection tries a series of thresholds to find the best edge detection parameters
func AutoTuneEdgeDetection(img gocv.Mat, imgName string) gocv.Mat {
	edges := gocv.NewMat()
	defer edges.Close()

	for t1 := 50.0; t1 <= 100.0; t1 += 10 {
		for t2 := 100.0; t2 <= 300.0; t2 += 50 {
			gocv.Canny(img, &edges, float32(t1), float32(t2))

			global.SaveResult(edges, fmt.Sprintf("%s_edges_t1_%.0f_t2_%.0f.jpg", imgName, t1, t2), "edge_tuning")
		}
	}

	return edges
}

func PreprocessImage(img gocv.Mat, strategy PreprocessStrategy) gocv.Mat {
	return strategy.Process(img)
}
