package philarios

import (
  "regexp"
  "strings"
  "unicode"
)

func SplitWords(sentence string) ([]string, error) {
  f := func(c rune) bool {
    return !unicode.isLetter(c) && !unicode.IsNumber(c)
  }
  return strings.FieldsFunc(sentence)
}
