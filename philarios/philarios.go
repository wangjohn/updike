package philarios

import (
  "database/sql"
)

type WordVector struct {
  Word string
  Score float64
}

/*
This method returns alternative words that can be used in place of the current
target in the target phrase.
*/
func AlternativeWords(word string) ([]string, error) {
  var alternativeWords []string

  targetVectors, err := TargetVectors(word)
  if err != nil {
    return alternativeWords, err
  }

  var potentialWords map[string]float64
  for _, wordVector := range targetVectors {
    potentialWords[wordVector.Word] = wordVector.Score
  }

  synonyms, err := Synonyms(word)
  if err != nil {
    return alternativeWords, err
  }

  for _, synonym := range synonyms {
    synonymVectors, err := TargetVectors(synonym)
    if err != nil {
      return alternativeWords, err
    }

    for _, synonymVector := range synonymVectors {
      potentialWords[synonymVector.Word] = SynonymScore(synonymVector.Score)
    }
  }

  // TODO: find the best n alternative words.
}

func SynonymScore(score float64) (float64) {
  return score
}

func TargetVectors(word string) ([]WordVector, error) {
  var wordVectors []WordVector

  // TODO: figure out a way to get pre and post target vectors
}

func Synonyms(word string) ([]string, error) {
  // TODO: implement
}
