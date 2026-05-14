package fileprocessing

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

type indexTestResult struct {
	Err        bool
	Occurences []occurrenceIndex
}

type indexTest struct {
	Haystack string
	Needles  []string
	Mode     AmbiguousResolutionMode
	Result   indexTestResult
}

var indexTests = []indexTest{
	{"", []string{"a"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"fo", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"foo", []string{"f"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"oofofoofooo", []string{"f"}, ReturnError, indexTestResult{false, []occurrenceIndex{{2, 0}, {4, 0}, {7, 0}}}},
	{"oofofoofooo", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{{4, 0}, {7, 0}}}},
	{"barfoobarfoo", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{{3, 0}, {9, 0}}}},
	{"foo", []string{"o"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}, {2, 0}}}},
	{"abcABCabc", []string{"A"}, ReturnError, indexTestResult{false, []occurrenceIndex{{3, 0}}}},
	{"abcABCabc", []string{"a"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}, {6, 0}}}},
	{"", []string{""}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"a"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"fo", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"foo", []string{"foo"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"foo", []string{""}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"foo", []string{"o"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}, {2, 0}}}},
	{"jrzm6jjhorimglljrea4w3rlgosts0w2gia17hno2td4qd1jz", []string{"jz"}, ReturnError, indexTestResult{false, []occurrenceIndex{{47, 0}}}},
	{"ekkuk5oft4eq0ocpacknhwouic1uua46unx12l37nioq9wbpnocqks6", []string{"ks6"}, ReturnError, indexTestResult{false, []occurrenceIndex{{52, 0}}}},
	{"999f2xmimunbuyew5vrkla9cpwhmxan8o98ec", []string{"98ec"}, ReturnError, indexTestResult{false, []occurrenceIndex{{33, 0}}}},
	{"9lpt9r98i04k8bz6c6dsrthb96bhi", []string{"96bhi"}, ReturnError, indexTestResult{false, []occurrenceIndex{{24, 0}}}},
	{"55u558eqfaod2r2gu42xxsu631xf0zobs5840vl", []string{"5840vl"}, ReturnError, indexTestResult{false, []occurrenceIndex{{33, 0}}}},
	{"", []string{"a"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"x", []string{"a"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"x", []string{"x"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"abc", []string{"a"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"abc", []string{"b"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"abc", []string{"c"}, ReturnError, indexTestResult{false, []occurrenceIndex{{2, 0}}}},
	{"abc", []string{"x"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"ab"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"bc", []string{"ab"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"ab", []string{"ab"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"xab", []string{"ab"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"xab"[:2], []string{"ab"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"abc"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xbc", []string{"abc"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"abc", []string{"abc"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"xabc", []string{"abc"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"xabc"[:3], []string{"abc"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xabxc", []string{"abc"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"abcd"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xbcd", []string{"abcd"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"abcd", []string{"abcd"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"xabcd", []string{"abcd"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"xyabcd"[:5], []string{"abcd"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xbcqq", []string{"abcqq"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"abcqq", []string{"abcqq"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"xabcqq", []string{"abcqq"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"xyabcqq"[:6], []string{"abcqq"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xabxcqq", []string{"abcqq"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xabcqxq", []string{"abcqq"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"01234567"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"32145678", []string{"01234567"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"01234567", []string{"01234567"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"x01234567", []string{"01234567"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"x0123456x01234567", []string{"01234567"}, ReturnError, indexTestResult{false, []occurrenceIndex{{9, 0}}}},
	{"xx01234567"[:9], []string{"01234567"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"3214567844", []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"0123456789", []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"x0123456789", []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"x012345678x0123456789", []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{{11, 0}}}},
	{"xyz0123456789"[:12], []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"x01234567x89", []string{"0123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"0123456789012345"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"3214567889012345", []string{"0123456789012345"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"0123456789012345", []string{"0123456789012345"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"x0123456789012345", []string{"0123456789012345"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"x012345678901234x0123456789012345", []string{"0123456789012345"}, ReturnError, indexTestResult{false, []occurrenceIndex{{17, 0}}}},
	{"", []string{"01234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"32145678890123456789", []string{"01234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"01234567890123456789", []string{"01234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"x01234567890123456789", []string{"01234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"x0123456789012345678x01234567890123456789", []string{"01234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{{21, 0}}}},
	{"xyz01234567890123456789"[:22], []string{"01234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"0123456789012345678901234567890"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"321456788901234567890123456789012345678911", []string{"0123456789012345678901234567890"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"0123456789012345678901234567890", []string{"0123456789012345678901234567890"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"x0123456789012345678901234567890", []string{"0123456789012345678901234567890"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"x012345678901234567890123456789x0123456789012345678901234567890", []string{"0123456789012345678901234567890"}, ReturnError, indexTestResult{false, []occurrenceIndex{{32, 0}}}},
	{"xyz0123456789012345678901234567890"[:33], []string{"0123456789012345678901234567890"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"", []string{"01234567890123456789012345678901"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"32145678890123456789012345678901234567890211", []string{"01234567890123456789012345678901"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"01234567890123456789012345678901", []string{"01234567890123456789012345678901"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"x01234567890123456789012345678901", []string{"01234567890123456789012345678901"}, ReturnError, indexTestResult{false, []occurrenceIndex{{1, 0}}}},
	{"x0123456789012345678901234567890x01234567890123456789012345678901", []string{"01234567890123456789012345678901"}, ReturnError, indexTestResult{false, []occurrenceIndex{{33, 0}}}},
	{"xyz01234567890123456789012345678901"[:34], []string{"01234567890123456789012345678901"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xxxxxx012345678901234567890123456789012345678901234567890123456789012", []string{"012345678901234567890123456789012345678901234567890123456789012"}, ReturnError, indexTestResult{false, []occurrenceIndex{{6, 0}}}},
	{"", []string{"0123456789012345678901234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"0123456789012345678901234567890123456789"}, ReturnError, indexTestResult{true, []occurrenceIndex{}}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"0123456789012345678901234567890123456789"}, PickFirst, indexTestResult{false, []occurrenceIndex{{2, 0}}}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"0123456789012345678901234567890123456789"}, PickSecond, indexTestResult{false, []occurrenceIndex{{22, 0}}}},
	{"xx012345678901234567890123456789012345678901234567890123456789012"[:41], []string{"0123456789012345678901234567890123456789"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"0123456789012345678901234567890123456xxx"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxx", []string{"0123456789012345678901234567890123456xxx"}, ReturnError, indexTestResult{false, []occurrenceIndex{{65, 0}}}},
	{"oxoxoxoxoxoxoxoxoxoxoxoy", []string{"oy"}, ReturnError, indexTestResult{false, []occurrenceIndex{{22, 0}}}},
	{"oxoxoxoxoxoxoxoxoxoxoxox", []string{"oy"}, ReturnError, indexTestResult{false, []occurrenceIndex{}}},
	{"oxoxoxoxoxoxoxoxoxoxox☺", []string{"☺"}, ReturnError, indexTestResult{false, []occurrenceIndex{{22, 0}}}},
	{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxx\xed\x9f\xc0", []string{"\xed\x9f\xc0"}, ReturnError, indexTestResult{false, []occurrenceIndex{{105, 0}}}},
	{"aabbaa", []string{"aa"}, ReturnError, indexTestResult{false, []occurrenceIndex{{0, 0}, {4, 0}}}},
	{"aabbaa", []string{"aa", "a"}, PickFirst, indexTestResult{false, []occurrenceIndex{{0, 1}, {1, 1}, {4, 1}, {5, 1}}}},
	{"aabbaa", []string{"aa", "a"}, PickSecond, indexTestResult{false, []occurrenceIndex{{1, 1}, {5, 1}}}},
	{"aabbaa", []string{"aab", "abba"}, PickFirst, indexTestResult{false, []occurrenceIndex{{0, 0}}}},
	{"aabbaa", []string{"aab", "abba"}, PickSecond, indexTestResult{false, []occurrenceIndex{{1, 1}}}},
	{"aabbaa", []string{"aab", "aabb"}, PickSecond, indexTestResult{false, []occurrenceIndex{{0, 1}}}},
	{"oxoxoxo", []string{"oxo", "ox"}, PickFirst, indexTestResult{false, []occurrenceIndex{{0, 1}, {2, 1}, {4, 1}}}},
	{"oxoxoxo", []string{"oxo", "xox"}, PickFirst, indexTestResult{false, []occurrenceIndex{{0, 0}, {3, 1}}}},
}

func TestKmpSearch(t *testing.T) {
	for _, test := range indexTests {
		result, err := wrappedKmpSearch(test.Haystack, test.Needles, test.Mode)
		if !slices.Equal(result, test.Result.Occurences) || ((err != nil) != test.Result.Err) {
			t.Errorf("KmpSearch(%q,%q) = %t, %v, %v; want %v", test.Haystack, test.Needles, err != nil, result, test.Mode, test.Result)
		}
	}
}

type replacementReturn struct {
	s   string
	err bool
}
type replacementTest struct {
	haystack     string
	needles      []string
	replacements []string
	result       replacementReturn
}

var multiReplaceAllTest = []replacementTest{
	{"", []string{""}, []string{""}, replacementReturn{"", false}},
	{"a", []string{"a"}, []string{"b"}, replacementReturn{"b", false}},
	{"aaa", []string{"a"}, []string{"b"}, replacementReturn{"bbb", false}},
	{"aaa", []string{"a"}, []string{"bbb"}, replacementReturn{"bbbbbbbbb", false}},
	{"aaa", []string{"aa"}, []string{"b"}, replacementReturn{"", true}},
	{"aabbaa", []string{"aa", "bb"}, []string{"bb", "aa"}, replacementReturn{"bbaabb", false}},
	{"ǂo", []string{"ǂ"}, []string{"a"}, replacementReturn{"ao", false}},
	{"ǂo", []string{"ǂ"}, []string{"Ɵ"}, replacementReturn{"Ɵo", false}},
	{"ǂȭ", []string{"ȴ"}, []string{"ȹ"}, replacementReturn{"ǂȭ", false}},
	{"ǂȭ", []string{"ȭ"}, []string{"ȴ"}, replacementReturn{"ǂȴ", false}},
	{"ǂɶȭɶ", []string{"ɶ"}, []string{"ȴ"}, replacementReturn{"ǂȴȭȴ", false}},
	{"ǂɶȭɶ", []string{"ɶ", "ȭ", "ǂ"}, []string{"ȴ", "ȴ", "ȴ"}, replacementReturn{"ȴȴȴȴ", false}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"0123456789012345678901234567890123456789"}, []string{"a"}, replacementReturn{"", true}},
	{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxx\xed\x9f\xc0", []string{"\xed\x9f\xc0"}, []string{"a"}, replacementReturn{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxxa", false}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"01234567890123456789"}, []string{"a"}, replacementReturn{"", true}},
}

func TestMultiReplaceAll(t *testing.T) {
	for _, test := range multiReplaceAllTest {
		actual, err := multiReplaceAll(test.haystack, test.needles, test.replacements)
		errorResult := err != nil
		if test.result.s != actual || errorResult != test.result.err {
			t.Errorf("multiReplaceAll(%q,%q) = %v; want %v", test.haystack, test.needles, replacementReturn{actual, errorResult}, test.result)
		}
	}
}

type pathtype int

const (
	pathtypeUrl = iota
	pathtypFilepath
)

type urlOrFilepath struct {
	s string
	t pathtype
}

type multiReplaceAllBenchmark struct {
	path         urlOrFilepath
	needles      []string
	replacements []string
}

var multiReplaceAllBenchmarkAllTests = []multiReplaceAllBenchmark{
	{urlOrFilepath{"https://raw.githubusercontent.com/wess/iotr/master/lotr.txt", pathtypeUrl}, []string{"Hobbit", "hobbit", "Gandalf", "Thorin", "sk-1234ijkl5678mnop1234ijkl5678mnop1234ijkl", "people", "abcdefabcdefabcdefabcdefabcdefabcdef12"}, []string{"tibboH", "tibboh", "Magic man", "ǂɶȭɶ", "a", "dudes", "AAAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAA"}},
}

func BenchmarkMultiReplaceAll(b *testing.B) {
	texts := make([]string, len(multiReplaceAllBenchmarkAllTests))
	for i, test := range multiReplaceAllBenchmarkAllTests {
		switch test.path.t {
		case pathtypeUrl:
			resp, err := http.Get(test.path.s)
			if err != nil {
				b.Skip(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				b.Skip(err)
			}
			texts[i] = string(body)
		case pathtypFilepath:
			cwd, err := os.Getwd()
			if err != nil {
				b.Skip(err)
			}
			path := filepath.Join(cwd, "../../test/fileProcessing/", test.path.s)
			data, err := os.ReadFile(path)
			if err != nil {
				b.Skip(err)
			}
			texts[i] = string(data)
		}
	}
	b.ResetTimer()
	for b.Loop() {
		for i, test := range multiReplaceAllBenchmarkAllTests {
			multiReplaceAll(texts[i], test.needles, test.replacements)
		}
	}
}
