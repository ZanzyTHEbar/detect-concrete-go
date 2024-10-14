package features

import (
	"image"
	"image/color"

	"gonum.org/v1/gonum/mat"
)

// ComputeGLCM calculates the GLCM matrix for texture analysis.
// Note: Could be done better, but a nested for loop is fine for now
// Iterate over each pixel in the image
func ComputeGLCM(img image.Image, angle, distance int) *mat.Dense {

	bounds := img.Bounds()
	glcm := mat.NewDense(256, 256, nil) // GLCM matrix for 256 gray levels

	// Iterate over each pixel in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			currentPixel := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			neighborX, neighborY := getNeighborCoordinates(x, y, angle, distance)
			if neighborX >= bounds.Min.X && neighborX < bounds.Max.X && neighborY >= bounds.Min.Y && neighborY < bounds.Max.Y {
				neighborPixel := color.GrayModel.Convert(img.At(neighborX, neighborY)).(color.Gray).Y
				glcm.Set(int(currentPixel), int(neighborPixel), glcm.At(int(currentPixel), int(neighborPixel))+1)
			}
		}
	}
	return glcm
}

// getNeighborCoordinates returns the neighboring pixel's coordinates based on angle and distance.
func getNeighborCoordinates(x, y, angle, distance int) (int, int) {
	switch angle {
	case 0: // Right
		return x + distance, y
	case 45: // Top-right diagonal
		return x + distance, y - distance
	case 90: // Top
		return x, y - distance
	case 135: // Top-left diagonal
		return x - distance, y - distance
	default:
		return x, y // Default to the same pixel (no movement)
	}
}

// ExtractHaralickFeatures extracts Haralick features from a GLCM matrix.
func ExtractHaralickFeatures(glcm *mat.Dense) []float64 {
	contrast := 0.0
	energy := 0.0
	homogeneity := 0.0
	correlation := 0.0
	meanI, meanJ := 0.0, 0.0
	numRows, numCols := glcm.Dims()

	// Calculate means for correlation
	total := 0.0
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			value := glcm.At(i, j)
			total += value
			meanI += float64(i) * value
			meanJ += float64(j) * value
		}
	}
	meanI /= total
	meanJ /= total

	// Iterate through GLCM to compute the features
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			value := glcm.At(i, j)
			diff := float64(i - j)

			contrast += diff * diff * value
			energy += value * value
			homogeneity += value / (1.0 + diff*diff)
			correlation += (float64(i) - meanI) * (float64(j) - meanJ) * value
		}
	}

	// Normalizing correlation
	if numRows > 1 && numCols > 1 {
		correlation /= float64(numRows * numCols)
	}

	return []float64{contrast, energy, homogeneity, correlation}
}
