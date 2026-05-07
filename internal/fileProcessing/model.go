package fileprocessing

type ReplacementPattern struct {
	Pattern     string
	Replacement string
}

func MultiReplaceAll(fileContent string, patterns []ReplacementPattern) (string, error) {
	if len(fileContent) == 0 || len(patterns) == 0 {
		return fileContent, nil
	}
	words := make([]string, len(patterns))
	replacements := make([]string, len(patterns))
	for i := range patterns {
		words[i] = patterns[i].Pattern
		replacements[i] = patterns[i].Replacement
	}
	return multiReplaceAll(fileContent, words, replacements)
}
