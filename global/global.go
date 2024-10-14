package global

import (
	"fmt"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
)

// SaveResult saves the processed image to the results directory.
func SaveResult(output gocv.Mat, imageName string, folder string) error {
	resultsDir := "/imgs/results"
	if folder != "" {
		resultsDir = filepath.Join(resultsDir, folder)
	}

	// Ensure the results directory exists
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		err := os.MkdirAll(resultsDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create results directory: %w", err)
		}
	}

	// Construct the full output path
	outputPath := filepath.Join(resultsDir, imageName)

	// Save the result
	if ok := gocv.IMWrite(outputPath, output); !ok {
		return fmt.Errorf("failed to save image: %s", outputPath)
	}

	fmt.Printf("Result saved to: %s\n", outputPath)
	return nil
}