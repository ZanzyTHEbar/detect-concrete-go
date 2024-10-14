package features

import (
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
)

// ExtractContours extracts contours from the edge-detected image.
func ExtractContours(img gocv.Mat) gocv.PointsVector {

	if img.Empty() {
		fmt.Println("Input image is empty, cannot extract contours.")
		return gocv.PointsVector{}
	}

	thresholded := gocv.NewMat()
	defer thresholded.Close()

	gocv.Threshold(img, &thresholded, 50, 255, gocv.ThresholdBinary)

	if thresholded.Empty() {
		fmt.Println("Thresholding resulted in an empty image.")
		return gocv.PointsVector{}
	}

	// Save thresholded image for debugging
	//gocv.IMWrite("thresholded_image.jpg", thresholded)

	// Perform contour detection
	contours := gocv.FindContours(thresholded, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// Check if contours are found
	if contours.Size() == 0 {
		fmt.Println("Found 0 contours")
		return gocv.PointsVector{}
	}

	filteredContours := gocv.NewPointsVector()

	fmt.Printf("Found %d contours\n", contours.Size())

	// Filter contours based on area and aspect ratio

	maxContours := 500

	for i := 0; i < contours.Size() && i < maxContours; i++ {
		contour := contours.At(i)

		// Ensure contour is valid and has sufficient size to process
		if contour.Size() < 5 { // Skip small/invalid contours
			continue
		}

		area := gocv.ContourArea(contour)
		aspectRatio := CalculateAspectRatio(contour)

		// Log the area and aspect ratio for debugging
		fmt.Printf("Contour %d: Area = %.2f, Aspect Ratio = %.2f\n", i, area, aspectRatio)

		if area > 12.0 && aspectRatio > 2.0 {
			filteredContours.Append(contour)
		}
	}

	return filteredContours
}

// OverlayContours draws contours onto the original image.
func OverlayContours(original gocv.Mat, contours gocv.PointsVector) gocv.Mat {
	// Create a copy of the original image
	output := original.Clone()

	// Draw the contours on the image (Red color, thickness 2)
	for i := 0; i < contours.Size(); i++ {
		gocv.DrawContours(&output, contours, i, color.RGBA{255, 0, 0, 0}, 2)
	}

	return output
}


// CalculateAspectRatio calculates the aspect ratio of a contour.
func CalculateAspectRatio(contour gocv.PointVector) float64 {
	rect := gocv.BoundingRect(contour)
	if rect.Dy() == 0 { // Prevent division by zero
		return 0
	}
	return float64(rect.Dx()) / float64(rect.Dy())
}
