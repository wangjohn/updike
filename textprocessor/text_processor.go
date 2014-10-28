package textprocessor

import (
  "strings"
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
  word = strings.ToLower(word)
  word = strings.TrimSpace(word)
  return word, nil
}
