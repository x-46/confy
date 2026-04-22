package fileprocessing

type ReplacementPattern struct {
	Pattern     string
	Replacement string
}

type FileProcessing interface {
	ProcessFile(filePath string, patterns []ReplacementPattern) (string, error)
}
