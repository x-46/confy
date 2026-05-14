package fileprocessing

import (
	"fmt"
	"strings"
)

type AmbiguousResolutionMode int

const (
	ReturnError AmbiguousResolutionMode = iota
	PickFirst
	PickSecond
	PickBoth
)

func (r AmbiguousResolutionMode) String() string {
	switch r {
	case ReturnError:
		return "ReturnError"
	case PickFirst:
		return "PickFirst"
	case PickSecond:
		return "PickSecond"
	case PickBoth:
		return "PickBoth"
	}
	return "unknown value for enum AmbiguousResolutionMode"
}

type occurrenceIndex struct {
	index, occurrence int
}

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

func wrappedKmpSearch(text string, words []string, mode AmbiguousResolutionMode) ([]occurrenceIndex, error) {
	if len(text) == 0 || len(words) == 0 {
		return []occurrenceIndex{}, nil
	}
	if len(words) == 1 && len(words[0]) == 0 {
		return []occurrenceIndex{}, nil
	}
	allEmpty := true
	for i := range words {
		if len(words[i]) != 0 {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		return []occurrenceIndex{}, nil
	}
	return kmpSearch(text, words, mode)
}

func kmpSearch(text string, words []string, mode AmbiguousResolutionMode) ([]occurrenceIndex, error) {
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
	j := 0
	occurrences := make([]occurrenceIndex, 0, 10*len(words))
	for j < len(text) {
		char := text[j]
		for i := range words {
			if len(words[i]) == 0 || len(words[i]) > len(text) {
				continue
			}
			vars[i].matched = words[i][vars[i].k] == char
			if vars[i].matched {
				vars[i].k += 1
			}
			if vars[i].k == len(words[i]) {
				occLen := len(occurrences)
				if occLen == 0 || mode == PickBoth || occurrences[occLen-1].index+len(words[occurrences[occLen-1].occurrence]) <= j+1-vars[i].k {
					occurrences = append(occurrences, occurrenceIndex{j + 1 - vars[i].k, i})
					vars[i].nP += 1
					vars[i].k = T[offsets[i]+vars[i].k]
					continue
				}
				switch mode {
				case ReturnError:
					lastOcc := occurrences[occLen-1]
					return nil, fmt.Errorf("Ambiguous resolution between \"%s\" at %d and \"%s\" at %d", words[lastOcc.occurrence], lastOcc.index, words[i], j+1-vars[i].k)
				case PickFirst:
					vars[i].k = T[offsets[i]+vars[i].k]
				case PickSecond:
					vars[occurrences[occLen-1].occurrence].nP -= 1
					occurrences[occLen-1] = occurrenceIndex{j + 1 - vars[i].k, i}
					vars[i].nP += 1
					vars[i].k = T[offsets[i]+vars[i].k]
				}
			}
		}
		for i := range words {
			if len(words[i]) == 0 || vars[i].matched {
				continue
			}
			vars[i].k = T[offsets[i]+vars[i].k]
			if vars[i].k < 0 {
				vars[i].k += 1
			}
			if words[i][vars[i].k] == char {
				vars[i].k += 1
			}
			vars[i].matched = false
		}
		j += 1
	}
	return occurrences, nil
}

func multiReplaceAll(text string, words []string, replacements []string) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	occurrences, err := kmpSearch(text, words, ReturnError)
	if err != nil {
		return "", err
	}
	if occurrences == nil {
		return "", nil
	}
	totalReplacementSize := 0
	for i := range occurrences {
		totalReplacementSize += len(replacements[occurrences[i].occurrence]) - len(words[occurrences[i].occurrence])
	}
	totalReplacementSize = max(0, totalReplacementSize)

	var b strings.Builder
	if len(text)+totalReplacementSize < 0 {
		return "", nil
	}
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
