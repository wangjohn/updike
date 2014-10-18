package philarios

import (
  "regexp"
  "github.com/wangjohn/quickselect"
)

type WordVector struct {
  Word string
  Score float64
}

type WordVectorCollection []WordVector

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
    wordVectors = append(wordVectors, synonymVectors...)
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

  return alternativeWords, nil
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

  rows, err := QueryForWord(word)
  if err != nil {
    return wordVectors, err
  }
  defer rows.Close()

  for rows.Next() {
    var body string
    if err = rows.Scan(&body); err != nil {
      return wordVectors, err
    }

    newVectors := paragraphTargetVectors(body, word)
    wordVectors = append(wordVectors, newVectors...)
  }

  return wordVectors
}

func paragraphTargetVectors(paragraph, word string) ([]WordVector) {
  wordVectors := make([]WordVector, 0)

  regexString := fmt.Printf(`[\s[[:punct:]]]*%s[\s[[:punct:]]]*`, word)
  compiledRegex := regexp.Compile(regexString)

  allIndices := compiledRegex.FindAllStringIndex(paragraph, -1)
  var potentialVector string
  if allIndices != nil {
    for i := 0; i < len(allIndices); i++ {
      matchedIndexPair := allIndices[i]
      start, end := matchedIndexPair[0], matchedIndexPair[1]

      potentialVector = lookForWord(paragraph[:start], false)
      if potentialVector != "" {
        wordVectors = append(wordVectors, potentialVector)
      }

      potentialVector := lookForword(paragraph[end:], true)
      if potentialVector != "" {
        wordVectors = append(wordVectors, potentialVector)
      }
    }
  }

  return wordVectors
}

func lookForWord(paragraph string, atBeginning bool) string {
  if strings.Trim(paragraph, " ") == "" {
    return ""
  }

  var compiledRegex regexp.Regexp
  if atBeginning {
    compiledRegex = regexp.Compile(`^(\w+)[\s[[:punct:]]]?`)
  } else {
    compiledRegex = regexp.Compile(`[\s[[:punct:]]]?(\w+)$`)
  }
  return compiledRegex.FindString(paragraph)
}

func Synonyms(word string) ([]string, error) {
  // TODO: implement
}
