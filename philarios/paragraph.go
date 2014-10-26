package philarios

import (
  "strings"
)

type Paragraph struct {
  PublicationId int
  Body string
}

/*
ProcessParagraphs processes a publication's text into a format which can be
used for storage, and returns the paragraphs that compose the publication.
*/
func (p PostgresStorage) ProcessParagraphs(publication Publication, publicationId int) ([]Paragraph, error) {
  var paragraphs []Paragraph
  preprocessedText, err := PreprocessPublicationText(publication)
  if err != nil {
    return paragraphs, err
  }

  preprocessedParagraphs := strings.Split(preprocessedText, "\n")
  var postprocessedParagraph string
  for _, preprocessedParagraph := range preprocessedParagraphs {
    postprocessedParagraph, err = PostprocessParagraph(preprocessedParagraph)
    if err != nil {
      return paragraphs, err
    }

    paragraphs = append(paragraphs, Paragraph{publicationId, postprocessedParagraph})
  }

  return paragraphs, nil
}

/*
PreprocessPublicationText takes a publication and preprocesses it so that the
text is ready to be turned into paragraphs.
*/
func PreprocessPublicationText(publication Publication) (string, error) {
  return publication.Text, nil
}

/*
PostprocessParagraph takes a paragraph and processes the text so that it is
ready to be stored in Storage.
*/
func PostprocessParagraph(paragraph string) (string, error) {
  return paragraph, nil
}
