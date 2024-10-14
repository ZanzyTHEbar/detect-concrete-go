package classification

import (
	"fmt"

	"gocv.io/x/gocv"
)

// Grade defines damage severity levels.
type Grade int

const (
	Minor Grade = iota
	Moderate
	Severe
)

type Grader struct{}

func (g *Grader) GradeDamage(contours gocv.PointsVector, totalArea, avgArea float64) Grade {

	fmt.Printf("Number of contours: %d, Total Area: %.2f, Average Area: %.2f\n", contours.Size(), totalArea, avgArea)

	if contours.Size() < 3 && totalArea < 150.0 {
		return Minor
	}

	if contours.Size() >= 3 && totalArea < 400.0 && avgArea < 150.0 {
		return Moderate
	}

	return Severe
}

func (g *Grader) TotalContourArea(contours gocv.PointsVector) (float64, float64) {
	var totalArea float64
	for i := 0; i < contours.Size(); i++ {
		totalArea += gocv.ContourArea(contours.At(i))
	}

	avgArea := totalArea / float64(contours.Size())

	return totalArea, avgArea
}

func (g Grade) String() string {
	return [...]string{"Minor", "Moderate", "Severe"}[g]
}
