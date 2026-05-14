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

type PatternMatch struct {
	PatternIndex int
	Index        int
	Line         int
	Column       int
}

func MultipleIndex(fileContent string, patterns []string) []PatternMatch {
	if len(fileContent) == 0 || len(patterns) == 0 {
		return []PatternMatch{}
	}
	patterns = append(patterns, "\n")
	lineEndId := len(patterns) - 1
	occurrences, _ := kmpSearch(fileContent, patterns, PickBoth)

	lineEnds := []int{-1}
	matches := make([]PatternMatch, 0, len(occurrences))
	for i := range occurrences {
		if occurrences[i].occurrence == lineEndId {
			lineEnds = append(lineEnds, occurrences[i].index)
		} else {
			matches = append(matches, PatternMatch{
				occurrences[i].occurrence,
				occurrences[i].index,
				len(lineEnds),
				occurrences[i].index - lineEnds[len(lineEnds)-1],
			})
		}
	}

	return matches
}
