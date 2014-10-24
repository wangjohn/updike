package philarios

import (
  "strings"
  "unicode"
)

func SplitWords(sentence string) ([]string) {
  f := func(c rune) bool {
    return !unicode.IsLetter(c) && !unicode.IsNumber(c)
  }
  return strings.FieldsFunc(sentence)
}
