package philarios

import (
  "database/sql"
  "github.com/wangjohn/quickselect"
)

type WordVector struct {
  Word string
  Score float64
}

type WordVectorCollection []WordVector struct

func (t WordVectorCollection) Len() int {
  return len(t)
}

func (t WordVectorCollection) Less(i, j int) bool {
  return synonymScore(t[i].Score) < synonymScore(t[j].Score)
}

func (t WordVectorCollection) Swap(i, j int) {
  t[i], t[j] = t[j], t[i]
}

/*
AlternativeWords returns alternative words that can be used in place of the
current word. The words that are returned will have close meanings to the word
used as an argument, but which are usually used in place of that word.

The maxWords parameter specifies the maximum number of words to return.
*/
func AlternativeWords(word string, maxWords int) ([]string, error) {
  var alternativeWords []string
  wordVectors := make(WordVectorCollection, 0)

  targetVectors, err := TargetVectors(word)
  if err != nil {
    return alternativeWords, err
  }
  wordVectors = append(wordVectors, targetVectors...)

  synonyms, err := Synonyms(word)
  if err != nil {
    return alternativeWords, err
  }

  for _, synonym := range synonyms {
    synonymVectors, err := TargetVectors(synonym)
    wordVectors = append(wordVectors, synonymVectors)
  }

  var wordsToSelect int
  if len(wordVectors) < maxWords {
    wordsToSelect = len(wordVectors)
  } else {
    wordsToSelect = maxWords
  }
  quickselect.QuickSelect(wordVectors, wordsToSelect)

  alternativeWords = make([]string, wordsToSelect)
  for i := 0; i < wordsToSelect; i++ {
    alternativeWords[i] = wordVectors[i].Word
  }

  return alternativeWords
}

/*
The synonymScore method specifies the mapping between a word's score and its
rank relative to other words. This function allows finer control over the words
that are returned.
*/
func synonymScore(score float64) (float64) {
  return score
}

func TargetVectors(word string) ([]WordVector, error) {
  var wordVectors []WordVector

  // TODO: figure out a way to get pre and post target vectors
}

func Synonyms(word string) ([]string, error) {
  // TODO: implement
}
