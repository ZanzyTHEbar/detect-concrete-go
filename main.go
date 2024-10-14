package main

import (
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"sync"

	"gocv.io/x/gocv"

	"github.com/ZanzyTHEbar/concrete-damage-detection/classification"
	"github.com/ZanzyTHEbar/concrete-damage-detection/features"
	"github.com/ZanzyTHEbar/concrete-damage-detection/global"
	"github.com/ZanzyTHEbar/concrete-damage-detection/processing"
)

func main() {

	imageDir := "/imgs"

	var waitGroup sync.WaitGroup

	err := filepath.WalkDir(imageDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return err
		}

		// Only process files that are images (adjust extension filter as needed)
		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".jpg" || filepath.Ext(entry.Name()) == ".png" || filepath.Ext(entry.Name()) == ".jpeg") {
			waitGroup.Add(1)
			go func(imagePath string) {
				defer waitGroup.Done()
				processImage(imagePath)
			}(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error while walking through directory: %v\n", err)
		return
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()
	fmt.Println("All images processed.")
}

func processImage(path string) {

	fmt.Println("Processing image: ", path)

	imgLoader := processing.NewImageLoader()
	img, err := imgLoader.LoadImage(path)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	grayImg := gocv.NewMat()
	defer grayImg.Close()

	blurredImg := gocv.NewMat()
	defer blurredImg.Close()

	// Convert image to grayscale for preprocessing
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	// Apply Gaussian blur to reduce texture noise
	gocv.GaussianBlur(grayImg, &blurredImg, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	// For Debugging: to find the optimum edge detection parameters
	_ = processing.AutoTuneEdgeDetection(blurredImg, filepath.Base(path))

	edges := processing.PreprocessImage(grayImg, &processing.EdgeDetectionStrategy{Threshold1: 60, Threshold2: 300})
	//gocv.AdaptiveThreshold(blurredImg, &edges, 100, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 11, 2)

	// Attempt to detect damaged regions
	contours := features.ExtractContours(edges)
	if contours.IsNil() || contours.Size() == 0 {
		fmt.Println("No contours detected")
		return
	}

	resultImg := features.OverlayContours(img, contours)

	imageName := fmt.Sprintf("%s_damages.jpg", filepath.Base(path))

	err = global.SaveResult(resultImg, imageName, "features")
	if err != nil {
		fmt.Println("Error saving result image:", err)
	}

	// Convert OpenCV Mat to Go image.Image for feature extraction
	goImage, _ := edges.ToImage()

	// Compute GLCM at 0-degree angle and 1-pixel distance
	glcm := features.ComputeGLCM(goImage, 0, 1)
	textureFeatures := features.ExtractHaralickFeatures(glcm)

	// For each contour, calculate the area and grade the severity
	grader := &classification.Grader{}
	totalArea, avgArea := grader.TotalContourArea(contours)

	damageGrade := grader.GradeDamage(contours, totalArea, avgArea)
	fmt.Printf("Image: %s, Total Area = %.2f, Average Area = %.2f, Damage Grade = %s\n", path, totalArea, avgArea, damageGrade.String())

	//for i := 0; i <= contours.Size(); i++ {
	//	contour := contours.At(i)
	//
	//	area := gocv.ContourArea(contour)
	//
	//	aspectRatio := features.CalculateAspectRatio(contour)
	//
	//	grade := grader.GradeDamage(area)
	//  fmt.Printf("Contour %d: Area = %.2f, Damage Grade = %s\n, Aspect Ratio: %.2f", i, area, grade.String(), aspectRatio)
	//
	//}

	// Classify damage using decision tree classifier
	classifier := &classification.DecisionTreeClassifier{}
	damageType := classifier.Classify(textureFeatures, contours.Size(), avgArea)

	fmt.Printf("Image %s: Detected damage type = %s\n", path, damageType)
}
