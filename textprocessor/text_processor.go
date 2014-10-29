package textprocessor

import (
  "strings"
  "fmt"
  "reflect"
  "unicode/utf8"
)

var vowels = map[rune]bool {
  'a': true,
  'e': true,
  'i': true,
  'o': true,
  'u': true,
  'y': true,
}
var ingRunes = []rune{'i', 'n', 'g'}
var edRunes = []rune{'e', 'd'}

/*
ProcessParagraphs processes a publication's text into a format which can be
used for storage, and returns the paragraphs that compose the publication.
*/
func ProcessParagraphs(text string) ([]string, error) {
  var paragraphs []string
  preprocessedText, err := preprocessPublicationText(text)
  if err != nil {
    return paragraphs, err
  }

  preprocessedParagraphs := strings.Split(preprocessedText, "\n")
  var postprocessedParagraph string
  for _, preprocessedParagraph := range preprocessedParagraphs {
    postprocessedParagraph, err = postprocessParagraph(preprocessedParagraph)
    if err != nil {
      return paragraphs, err
    }

    if postprocessedParagraph != "" {
      paragraphs = append(paragraphs, postprocessedParagraph)
    }
  }

  return paragraphs, nil
}

/*
PreprocessPublicationText takes a publication and preprocesses it so that the
text is ready to be turned into paragraphs.
*/
func preprocessPublicationText(text string) (string, error) {
  return text, nil
}

/*
PostprocessParagraph takes a paragraph and processes the text so that it is
ready to be stored in Storage.
*/
func postprocessParagraph(paragraph string) (string, error) {
  return strings.TrimSpace(paragraph), nil
}

/*
NormaizedWord takes a word and returns a word that is normalized and ready to
be stored.
*/
func NormalizedWord(word string) (string, error) {
  if !utf8.ValidString(word) {
    return "", fmt.Errorf("The word '%s' is not a valid utf-8 string.", word)
  }

  word = strings.ToLower(word)
  word = strings.TrimSpace(word)

  //runes, err := getRunesFromString(word)

  return word, nil
}

func getRunesFromString(str string) ([]rune) {
  runes := make([]rune, 0)
  for i, w := 0, 0; i < len(str); i += w {
    runeValue, width := utf8.DecodeRuneInString(str[i:])
    runes = append(runes, runeValue)
    w += width
  }

  return runes
}

func handleSuffixIng(runes []rune) (bool, []rune) {
  length := len(runes)
  if length >= 5 {
    lastThree := runes[(length-3):length]
    if reflect.DeepEqual(lastThree, ingRunes) {
      precedingRune := runes[length-4]

      // For `ing` suffixes that are preceded by a `y`, like `staying` or `lying`
      if precedingRune == 'y' {
        result := deepCopy(runes[:(length-4)])
        result = append(result, 'i', 'e')
        result = deepCopy(runes[:(length-3)])
        if isWord(result) {
          return true, result
        } else {
          return true, result
        }

      // For `ing` suffixes that are preceded by a double consonant, like `stopping`
      } else if !isVowel(precedingRune) {
        if precedingRune == runes[length-5] {
          result := deepCopy(runes[:(length-4)])
          if isWord(result) {
            return true, result
          }
        }

      // For all other `ing` suffixes
      } else {
        result := deepCopy(runes[:(length-3)])
        if isWord(result) {
          return true, result
        }
      }
    }
  }

  return false, []rune{}
}

func handleSuffixEd(runes []rune) (bool, []rune) {
  length := len(runes)
  if length >= 4 {
    lastTwo := runes[(length-2):length]
    if reflect.DeepEqual(lastTwo, edRunes) {

    }
  }

  return false, []rune{}
}

func isVowel(char rune) (bool) {
  return vowels[char]
}

func deepCopy(runes []rune) ([]rune) {
  result := make([]rune, len(runes))
  for i, r := range runes {
    result[i] = r
  }

  return result
}

func isWord(runes []rune) (bool) {
  return true
}
