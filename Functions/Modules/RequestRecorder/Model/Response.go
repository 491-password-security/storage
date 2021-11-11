package Model

type DeviceAnalysisQueryResponse struct {
	Id                              string `json:"id"`
	DeviceModel                     string `json:"deviceModel"`
	ServiceName                     string `json:"serviceName"`
	Build                           string `json:"build"`
	Version                         string `json:"version"`
	AverageFrameCount               string `json:"averageFrameCount"`
	AverageFaceDetectionCount       string `json:"averageFaceDetectionCount"`
	AverageSkinExtractionCount      string `json:"averageSkinExtractionCount"`
	AverageQualityIndex             string `json:"averageQualityIndex"`
	FaceDetectionToTotalFrameRatio  string `json:"faceDetectionToTotalFrameRatio"`
	SkinExtractionToTotalFrameRatio string `json:"skinExtractionToTotalFrameRatio"`
	TotalMeasurementCount           int64  `json:"totalMeasurementCount"`
	LastChanged                     string `json:"lastChanged"`
	FirstSet						string `json:"firstSet"`
}
