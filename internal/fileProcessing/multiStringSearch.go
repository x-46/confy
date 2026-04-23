package fileprocessing

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

func computeTable(words []string) []map[int]int {
	T := make([]map[int]int, len(words))
	for i := range T {
		T[i] = make(map[int]int)
	}
	for i := range words {
		pos := 1
		cnd := 0
		T[i][0] = -1
		for pos < len(words[i]) {
			if words[i][pos] == words[i][cnd] {
				T[i][pos] = T[i][cnd]
			} else {
				T[i][pos] = cnd
				for cnd >= 0 && words[i][pos] != words[i][cnd] {
					cnd = T[i][cnd]
				}
			}
			pos += 1
			cnd += 1
		}
		T[i][pos] = cnd
	}
	return T
}

func KmpSearch(text string, words []string) [][]int {
	j := 0
	T := computeTable(words)
	nP := make([]int, len(words))
	occurences := make([][]int, len(words))
	k := make([]int, len(words))
	matched := make([]bool, len(words))
	for j < len(text) {
		char := text[j]
		for i := range matched {
			matched[i] = false
		}
		matchedAny := false
		for i := range words {
			if len(words[i]) == 0 || len(words[i]) > len(text) {
				continue
			}
			matched[i] = words[i][k[i]] == char
			if matched[i] {
				k[i] += 1
				matchedAny = true
				if k[i] == len(words[i]) {
					occurences[i] = append(occurences[i], j+1-k[i])
					nP[i] += 1
					k[i] = T[i][k[i]]
				}
			}
		}
		if matchedAny {
			j += 1
			for i := range matched {
				if !matched[i] {
					k[i] = T[i][k[i]]
					if k[i] < 0 {
						k[i] += 1
					}
				}
			}
		} else {
			incJ := false
			for i := range words {
				k[i] = T[i][k[i]]
				if k[i] < 0 {
					k[i] += 1
					incJ = true
				}
			}
			if incJ {
				j += 1
			}
		}
	}
	return occurences
}

func multiReplaceAll(text string, words []string, replacements []string) (string, error) {
	occurenceMatrix := KmpSearch(text, words)
	totalOccs := 0

	for i := range occurenceMatrix {
		totalOccs += len(occurenceMatrix[i])
	}
	indexWithOccurence := make([][2]int, totalOccs)
	occIdx := 0
	totalReplacementSize := 0
	for i := range occurenceMatrix {
		for j := range occurenceMatrix[i] {
			indexWithOccurence[occIdx] = [2]int{i, occurenceMatrix[i][j]}
			occIdx += 1
		}
		totalReplacementSize += len(occurenceMatrix[i]) * (len(replacements[i]) - len(words[i]))
	}
	slices.SortFunc(indexWithOccurence, func(a, b [2]int) int {
		return cmp.Compare(a[1], b[1])
	})
	fmt.Println(indexWithOccurence)
	for i := 0; i < len(indexWithOccurence)-1; i += 1 {
		if indexWithOccurence[i][1]+len(words[indexWithOccurence[i][0]]) > indexWithOccurence[i+1][1] {
			return "", fmt.Errorf("unambigious resolution between \"%s\" at %d and \"%s\" at %d", words[indexWithOccurence[i][0]], indexWithOccurence[i][1], words[indexWithOccurence[i+1][0]], indexWithOccurence[i+1][1])
		}
	}

	var b strings.Builder
	b.Grow(len(text) + totalReplacementSize)
	lastTextIndex := 0
	for i := range indexWithOccurence {
		b.WriteString(text[lastTextIndex:indexWithOccurence[i][1]])
		b.WriteString(replacements[indexWithOccurence[i][0]])
		lastTextIndex = indexWithOccurence[i][1] + len(words[indexWithOccurence[i][0]])
	}
	if lastTextIndex < len(text) {
		b.WriteString(text[lastTextIndex:])
	}
	return b.String(), nil
}
