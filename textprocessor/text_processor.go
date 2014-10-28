package textprocessor

import (
  "strings"
  "errors"
  "encoding/utf8"
)

const (
  vowels = []rune{'a', 'e', 'i', 'o', 'u', 'y'}
  ingRunes = []rune{'i', 'n', 'g'}
)

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
    return "", errors.New("The word '%s' is not a valid utf-8 string.", word)
  }

  word = strings.ToLower(word)
  word = strings.TrimSpace(word)

  runes := getRunesFromString(word)

  return word, nil
}

func getRunesFromString(str string) ([]rune, error) {
  runes := make([]string, 0)
  for i, w := 0, 0; i < len(str); i += w {
    runeValue, width := utf8.DecodeRuneInString(str[i:])
    runes = append(runes, runeValue)
    w += width
  }

  return runes
}

func handleSuffixIng(runes []rune) (bool, string) {
  length := len(runes)
  if length >= 3 {
    lastThree := runes[(length-3):length]
    if reflect.DeepEqual(lastThree, ingRunes) {
      // TODO: implement ing suffix removal
    }
  }

  return false, ""
}
