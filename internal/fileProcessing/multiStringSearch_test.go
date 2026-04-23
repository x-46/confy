package fileprocessing

import (
	"slices"
	"testing"
)

type IndexTest struct {
	haystack string
	needle   string
	out      []int
}

var indexTests = []IndexTest{
	{"", "a", []int{-1}},
	{"", "foo", []int{-1}},
	{"fo", "foo", []int{-1}},
	{"foo", "foo", []int{0}},
	{"foo", "f", []int{0}},
	{"oofofoofooo", "f", []int{2, 4, 7}},
	{"oofofoofooo", "foo", []int{4, 7}},
	{"barfoobarfoo", "foo", []int{3, 9}},
	{"foo", "o", []int{1, 2}},
	{"abcABCabc", "A", []int{3}},
	{"abcABCabc", "a", []int{0, 6}},
	{"", "", []int{-1}},
	{"", "a", []int{-1}},
	{"", "foo", []int{-1}},
	{"fo", "foo", []int{-1}},
	{"foo", "foo", []int{0}},
	{"foo", "", []int{-1}},
	{"foo", "o", []int{1, 2}},
	{"jrzm6jjhorimglljrea4w3rlgosts0w2gia17hno2td4qd1jz", "jz", []int{47}},
	{"ekkuk5oft4eq0ocpacknhwouic1uua46unx12l37nioq9wbpnocqks6", "ks6", []int{52}},
	{"999f2xmimunbuyew5vrkla9cpwhmxan8o98ec", "98ec", []int{33}},
	{"9lpt9r98i04k8bz6c6dsrthb96bhi", "96bhi", []int{24}},
	{"55u558eqfaod2r2gu42xxsu631xf0zobs5840vl", "5840vl", []int{33}},
	{"", "a", []int{-1}},
	{"x", "a", []int{-1}},
	{"x", "x", []int{0}},
	{"abc", "a", []int{0}},
	{"abc", "b", []int{1}},
	{"abc", "c", []int{2}},
	{"abc", "x", []int{-1}},
	{"", "ab", []int{-1}},
	{"bc", "ab", []int{-1}},
	{"ab", "ab", []int{0}},
	{"xab", "ab", []int{1}},
	{"xab"[:2], "ab", []int{-1}},
	{"", "abc", []int{-1}},
	{"xbc", "abc", []int{-1}},
	{"abc", "abc", []int{0}},
	{"xabc", "abc", []int{1}},
	{"xabc"[:3], "abc", []int{-1}},
	{"xabxc", "abc", []int{-1}},
	{"", "abcd", []int{-1}},
	{"xbcd", "abcd", []int{-1}},
	{"abcd", "abcd", []int{0}},
	{"xabcd", "abcd", []int{1}},
	{"xyabcd"[:5], "abcd", []int{-1}},
	{"xbcqq", "abcqq", []int{-1}},
	{"abcqq", "abcqq", []int{0}},
	{"xabcqq", "abcqq", []int{1}},
	{"xyabcqq"[:6], "abcqq", []int{-1}},
	{"xabxcqq", "abcqq", []int{-1}},
	{"xabcqxq", "abcqq", []int{-1}},
	{"", "01234567", []int{-1}},
	{"32145678", "01234567", []int{-1}},
	{"01234567", "01234567", []int{0}},
	{"x01234567", "01234567", []int{1}},
	{"x0123456x01234567", "01234567", []int{9}},
	{"xx01234567"[:9], "01234567", []int{-1}},
	{"", "0123456789", []int{-1}},
	{"3214567844", "0123456789", []int{-1}},
	{"0123456789", "0123456789", []int{0}},
	{"x0123456789", "0123456789", []int{1}},
	{"x012345678x0123456789", "0123456789", []int{11}},
	{"xyz0123456789"[:12], "0123456789", []int{-1}},
	{"x01234567x89", "0123456789", []int{-1}},
	{"", "0123456789012345", []int{-1}},
	{"3214567889012345", "0123456789012345", []int{-1}},
	{"0123456789012345", "0123456789012345", []int{0}},
	{"x0123456789012345", "0123456789012345", []int{1}},
	{"x012345678901234x0123456789012345", "0123456789012345", []int{17}},
	{"", "01234567890123456789", []int{-1}},
	{"32145678890123456789", "01234567890123456789", []int{-1}},
	{"01234567890123456789", "01234567890123456789", []int{0}},
	{"x01234567890123456789", "01234567890123456789", []int{1}},
	{"x0123456789012345678x01234567890123456789", "01234567890123456789", []int{21}},
	{"xyz01234567890123456789"[:22], "01234567890123456789", []int{-1}},
	{"", "0123456789012345678901234567890", []int{-1}},
	{"321456788901234567890123456789012345678911", "0123456789012345678901234567890", []int{-1}},
	{"0123456789012345678901234567890", "0123456789012345678901234567890", []int{0}},
	{"x0123456789012345678901234567890", "0123456789012345678901234567890", []int{1}},
	{"x012345678901234567890123456789x0123456789012345678901234567890", "0123456789012345678901234567890", []int{32}},
	{"xyz0123456789012345678901234567890"[:33], "0123456789012345678901234567890", []int{-1}},
	{"", "01234567890123456789012345678901", []int{-1}},
	{"32145678890123456789012345678901234567890211", "01234567890123456789012345678901", []int{-1}},
	{"01234567890123456789012345678901", "01234567890123456789012345678901", []int{0}},
	{"x01234567890123456789012345678901", "01234567890123456789012345678901", []int{1}},
	{"x0123456789012345678901234567890x01234567890123456789012345678901", "01234567890123456789012345678901", []int{33}},
	{"xyz01234567890123456789012345678901"[:34], "01234567890123456789012345678901", []int{-1}},
	{"xxxxxx012345678901234567890123456789012345678901234567890123456789012", "012345678901234567890123456789012345678901234567890123456789012", []int{6}},
	{"", "0123456789012345678901234567890123456789", []int{-1}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", "0123456789012345678901234567890123456789", []int{2, 12, 22}},
	{"xx012345678901234567890123456789012345678901234567890123456789012"[:41], "0123456789012345678901234567890123456789", []int{-1}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", "0123456789012345678901234567890123456xxx", []int{-1}},
	{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxx", "0123456789012345678901234567890123456xxx", []int{65}},
	{"oxoxoxoxoxoxoxoxoxoxoxoy", "oy", []int{22}},
	{"oxoxoxoxoxoxoxoxoxoxoxox", "oy", []int{-1}},
	{"oxoxoxoxoxoxoxoxoxoxox☺", "☺", []int{22}},
	{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxx\xed\x9f\xc0", "\xed\x9f\xc0", []int{105}},
	{"aabbaa", "aa", []int{0, 4}},
}

func TestKmpSearch(t *testing.T) {
	for _, test := range indexTests {
		actual := KmpSearch(test.haystack, []string{test.needle})
		if len(actual[0]) == 0 {
			actual[0] = append(actual[0], -1)
		}
		if !slices.Equal(actual[0], test.out) {
			t.Errorf("KmpSearch(%q,%q) = %v; want %v", test.haystack, test.needle, actual, test.out)
		}
	}
}

type ReplacementReturn struct {
	s   string
	err bool
}
type ReplacementTest struct {
	haystack     string
	needles      []string
	replacements []string
	result       ReplacementReturn
}

var multiReplaceAllTest = []ReplacementTest{
	{"", []string{""}, []string{""}, ReplacementReturn{"", false}},
	{"a", []string{"a"}, []string{"b"}, ReplacementReturn{"b", false}},
	{"aaa", []string{"a"}, []string{"b"}, ReplacementReturn{"bbb", false}},
	{"aaa", []string{"a"}, []string{"bbb"}, ReplacementReturn{"bbbbbbbbb", false}},
	{"aaa", []string{"aa"}, []string{"b"}, ReplacementReturn{"", true}},
	{"aabbaa", []string{"aa", "bb"}, []string{"bb", "aa"}, ReplacementReturn{"bbaabb", false}},
	{"ǂo", []string{"ǂ"}, []string{"a"}, ReplacementReturn{"ao", false}},
	{"ǂo", []string{"ǂ"}, []string{"Ɵ"}, ReplacementReturn{"Ɵo", false}},
	{"ǂȭ", []string{"ȴ"}, []string{"ȹ"}, ReplacementReturn{"ǂȭ", false}},
	{"ǂȭ", []string{"ȭ"}, []string{"ȴ"}, ReplacementReturn{"ǂȴ", false}},
	{"ǂɶȭɶ", []string{"ɶ"}, []string{"ȴ"}, ReplacementReturn{"ǂȴȭȴ", false}},
	{"ǂɶȭɶ", []string{"ɶ", "ȭ", "ǂ"}, []string{"ȴ", "ȴ", "ȴ"}, ReplacementReturn{"ȴȴȴȴ", false}},
	{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxx\xed\x9f\xc0", []string{"\xed\x9f\xc0"}, []string{"a"}, ReplacementReturn{"xx0123456789012345678901234567890123456789012345678901234567890120123456789012345678901234567890123456xxxa", false}},
	{"xx012345678901234567890123456789012345678901234567890123456789012", []string{"01234567890123456789"}, []string{"a"}, ReplacementReturn{"", true}},
}

func TestMultiReplaceAll(t *testing.T) {
	for _, test := range multiReplaceAllTest {
		actual, err := multiReplaceAll(test.haystack, test.needles, test.replacements)
		errorResult := err != nil
		if test.result.s != actual || errorResult != test.result.err {
			t.Errorf("multiReplaceAll(%q,%q) = %v; want %v", test.haystack, test.needles, ReplacementReturn{actual, errorResult}, test.result)
		}
	}
}
