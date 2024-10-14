package classification

type Classifier interface {
	Classify(features []float64) string
}

type DecisionTreeClassifier struct{}

func (dtc *DecisionTreeClassifier) Classify(features []float64, numContours int, avgContourArea float64) string {
	//! This is basic, a proper decision tree classifier should be used in practice

	contrast, energy, homogeneity, correlation := features[0], features[1], features[2], features[3]

	// Uses Haralick and contour features
	if contrast > 1000.0 && energy < 0.1 && numContours > 3 {
		return "crack"
	} else if homogeneity > 0.5 && correlation > 0.2 && avgContourArea > 50.0 {
		return "spalling"
	}
	return "other"
}
