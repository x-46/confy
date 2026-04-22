package fileprocessing

type ReplacementPattern struct {
	Pattern     string
	Replacement string
}

type FileProcessing interface {
	ProcessFile(fileContent string, patterns []ReplacementPattern) (string, error)
}
