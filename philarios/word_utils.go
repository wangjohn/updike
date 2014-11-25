package philarios

import (
  "strings"
  "unicode"
)

/*
SplitWords takes a string and returns the separate words that make up the
string. In essence, it splits the string according to demarcating characters
like punctuation and spaces.
*/
func SplitWords(sentence string) ([]string) {
  f := func(c rune) bool {
    return unicode.IsPunct(c) || unicode.IsSpace(c) || unicode.IsSymbol(c)
  }
  return strings.FieldsFunc(sentence, f)
}

func FuzzyStringEquals(word1, word2 string) (bool) {
  w1 := strings.TrimSpace(strings.ToLower(word1))
  w2 := strings.TrimSpace(strings.ToLower(word1))
  return w1 == w2
}
