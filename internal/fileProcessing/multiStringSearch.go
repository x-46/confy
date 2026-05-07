package fileprocessing

import (
	"strings"
)

func computeTable(words []string) ([]int, []int) {
	offsets := make([]int, len(words))
	for i := 1; i < len(words); i++ {
		offsets[i] = offsets[i-1] + len(words[i-1]) + 1
	}
	T := make([]int, offsets[len(words)-1]+len(words[len(words)-1])+1)
	for i := range words {
		if len(words[i]) == 0 {
			continue
		}
		pos := 1
		cnd := 0
		T[offsets[i]] = -1
		for pos < len(words[i]) {
			if words[i][pos] == words[i][cnd] {
				T[offsets[i]+pos] = T[offsets[i]+cnd]
			} else {
				T[offsets[i]+pos] = cnd
				for cnd >= 0 && words[i][pos] != words[i][cnd] {
					cnd = T[offsets[i]+cnd]
				}
			}
			pos += 1
			cnd += 1
		}
		T[offsets[i]+pos] = cnd
	}
	return T, offsets
}

type occurrenceIndex struct {
	index, occurrence int
}

func wrappedKmpSearch(text string, words []string) [][]int {
	if len(text) == 0 || len(words) == 0 {
		return nil
	}
	if len(words) == 1 && len(words[0]) == 0 {
		return nil
	}
	allEmpty := true
	for i := range words {
		if len(words[i]) != 0 {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		return nil
	}
	occurrences := kmpSearch(text, words)
	if len(occurrences) == 0 {
		return nil
	}
	ret := make([][]int, len(words))
	for i := range occurrences {
		ret[occurrences[i].occurrence] = append(ret[occurrences[i].occurrence], occurrences[i].index)
	}
	return ret
}

func kmpSearch(text string, words []string) []occurrenceIndex {
	j := 0
	T, offsets := computeTable(words)
	type LoopVariables struct {
		nP      int
		k       int
		matched bool
	}
	vars := make([]LoopVariables, len(words))
	totalWordLen := 0
	for i := range words {
		totalWordLen += len(words[i])
	}
	occurrences := make([]occurrenceIndex, 0, 10*len(words))
	for j < len(text) {
		char := text[j]
		for i := range vars {
			vars[i].matched = false
		}
		matchedAny := false
		for i := range words {
			if len(words[i]) == 0 || len(words[i]) > len(text) {
				continue
			}
			vars[i].matched = words[i][vars[i].k] == char
			if vars[i].matched {
				vars[i].k += 1
				matchedAny = true
				if vars[i].k == len(words[i]) {
					occurrences = append(occurrences, occurrenceIndex{j + 1 - vars[i].k, i})
					vars[i].nP += 1
					vars[i].k = T[offsets[i]+vars[i].k]
				}
			}
		}
		if matchedAny {
			j += 1
			for i := range vars {
				if len(words[i]) == 0 {
					continue
				}
				if !vars[i].matched {
					vars[i].k = T[offsets[i]+vars[i].k]
					if vars[i].k < 0 {
						vars[i].k += 1
					}
				}
			}
		} else {
			jInc := 0
			for i := range words {
				if len(words[i]) == 0 {
					continue
				}
				vars[i].k = T[offsets[i]+vars[i].k]
				if vars[i].k < 0 {
					vars[i].k += 1
					jInc = 1
				}
			}
			j += jInc
		}
	}
	return occurrences
}

func multiReplaceAll(text string, words []string, replacements []string) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	occurrences := kmpSearch(text, words)
	if occurrences == nil {
		return "", nil
	}
	totalReplacementSize := 0
	for i := range occurrences {
		totalReplacementSize += len(replacements[occurrences[i].occurrence]) - len(words[occurrences[i].occurrence])
	}

	var b strings.Builder
	b.Grow(len(text) + totalReplacementSize)
	lastTextIndex := 0
	for i := range occurrences {
		b.WriteString(text[lastTextIndex:occurrences[i].index])
		b.WriteString(replacements[occurrences[i].occurrence])
		lastTextIndex = occurrences[i].index + len(words[occurrences[i].occurrence])
	}
	if lastTextIndex < len(text) {
		b.WriteString(text[lastTextIndex:])
	}
	return b.String(), nil
}
