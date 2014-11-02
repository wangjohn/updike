package textprocessor

var TextProcessorRules = []ProcessingRule{
  FilterBy("EndsWith", "ying").Try(Word{0, -4}, 'i', 'e'),
}

func NormalizedWord(word string) (string, error) {
  for _, rule := range TextProcessorRules {
    applies, normalizedRunes := rule.Rule(word)
    if applies {
      return string(normalizedRunes), nil
    }
  }

  return word, nil
}

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

func preprocessPublicationText(text string) (string, error) {
  return text, nil
}

func postprocessParagraph(paragraph string) (string, error) {
  return strings.TrimSpace(paragraph), nil
}
