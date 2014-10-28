package textprocessor

import (
  "github.com/wangjohn/updike/philarios"
  "strings"
)

/*
ProcessParagraphs processes a publication's text into a format which can be
used for storage, and returns the paragraphs that compose the publication.
*/
func ProcessParagraphs(publication philarios.Publication, publicationId int) ([]Paragraph, error) {
  var paragraphs []Paragraph
  preprocessedText, err := preprocessPublicationText(publication)
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

    if postprocessedParagrah != "" {
      paragraphs = append(paragraphs, Paragraph{publicationId, postprocessedParagraph})
    }
  }

  return paragraphs, nil
}

/*
PreprocessPublicationText takes a publication and preprocesses it so that the
text is ready to be turned into paragraphs.
*/
func preprocessPublicationText(publication philarios.Publication) (string, error) {
  return publication.Text, nil
}

/*
PostprocessParagraph takes a paragraph and processes the text so that it is
ready to be stored in Storage.
*/
func postprocessParagraph(paragraph string) (string, error) {
  return strings.TrimSpace(postprocessedParagrah), nil
}

/*
NormaizedWord takes a word and returns a word that is normalized and ready to
be stored.
*/
func NormalizedWord(word string) (string, error) {
  word = strings.ToLower(word)
  word = strings.Trim(word)
  return word, nil
}
