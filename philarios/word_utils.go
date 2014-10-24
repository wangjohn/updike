package philarios

import (
  "strings"
  "unicode"
)

func SplitWords(sentence string) ([]string) {
  f := func(c rune) bool {
    return unicode.IsPunct(c) || unicode.IsSpace(c) || unicode.IsSymbol(c)
  }
  return strings.FieldsFunc(sentence, f)
}
