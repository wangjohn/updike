package textprocessor

import (
  "strings"
  "fmt"
  "log"
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

type Consonant struct {
  Index int
}

func (c Consonant) IsConsonant(char rune) (bool) {
  return vowels[char]
}

type Word struct {
  StartIndex int
  EndIndex int
}

type filteredFunction struct {
  filterFunc func
}

type processingRule struct {
  processFunc func
}

func (w Word) GetSlice(word string) ([]rune) {
  wordRunes := getRunesFromString(word)
  startIndex := w.StartIndex

  var endIndex int
  if w.EndIndex < 0 {
    endIndex = len(wordRunes) - w.EndIndex + 1
  } else {
    endIndex = w.EndIndex
  }

  if endIndex < 0 {
    log.Fatal(`Specified an EndIndex on a word that cannot exist. You specified
               '%d', but the word only has length %d`, w.EndIndex, len(wordRunes))
  }

  if startIndex > endIndex {
    log.Fatal(`Specified a StartIndex '%d' on a word that is greater than the
              EndIndex '%d' (for word of length '%d').`,
              w.StartIndex, w.EndIndex, len(wordRunes))
  }

  return wordRunes[startIndex:endIndex]
}

func filterBy(filterFunc func) (filterFunction) {
  return filteredFunction{filterFunc}
}

func (f filterFunction) try(args ...interface{}) (func) {
  processFunc := func (word string) (bool, []rune) {
    if !f.filterFunc(word) {
      return false, []rune{}
    }

    runeResult := make([]rune, 0)
    for _, argument := range args {
      case argument.(type) {
      switch Word
        currentRunes := argument.GetSlice(word)
        runeResult = append(runeResult, currentRunes...)
      switch rune
        runeResult = append(runeResult, argument)
      default
        log.Fatal(`Used an unknown argument type in the 'try' method. Argument
                   '%q' with type '%T' cannot but used.`, argument, argument)
      }
    }
    return true, runeResult
  }

  return processingRule{processFunc}
}

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
    w = width
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
        ieResult := deepCopy(runes[:(length-4)])
        ieResult = append(ieResult, 'i', 'e')
        if isWord(ieResult) {
          return true, ieResult
        }

      // For `ing` suffixes that are preceded by a double consonant, like `stopping`
      } else if !isVowel(precedingRune) {
        if precedingRune == runes[length-5] {
          result := deepCopy(runes[:(length-4)])
          if isWord(result) {
            return true, result
          }
        }
      }

      // For all other `ing` suffixes
      result := deepCopy(runes[:(length-3)])
      if isWord(result) {
        return true, result
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
      precedingRune := runes[length-3]

      // For `ed` suffixes that are preceded by i, like `lied` or `died`
      if precedingRune == 'i' {
        result := deepCopy(runes[:(length-1)])
        if isWord(result) {
          return true, result
        }

      // For `ed` suffixes that are preceded by a constant
      } else if !isVowel(precedingRune) {
        // If preceded by a double consonant, liked `stopped`
        if precedingRune == runes[length-4] {
          result := deepCopy(runes[:(length-3)])
          if isWord(result) {
            return true, result
          }

        // If preceded by a dropped e, like `phoned` or `danced`
        } else {
          eResult := deepCopy(runes[:(length-2)])
          eResult = append(eResult, 'e')
          if isWord(eResult) {
            return true, eResult
          }
        }
      }

      // For all other `ed` suffixes
      result := deepCopy(runes[:(length-2)])
      if isWord(result) {
        return true, result
      }
    }
  }

  return false, []rune{}
}

func handleApostrophe(runes []rune) (bool, []rune) {

}

func handleImproperWord(runes []rune) (bool, []rune) {
  // TODO: implement
  return false, []rune{}
}

func isVowel(char rune) (bool) {
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
