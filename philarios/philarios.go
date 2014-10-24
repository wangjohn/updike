package philarios

import (
  "regexp"
  "github.com/wangjohn/quickselect"

  "fmt"
  "strings"
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
SentenceFindAlternativeWords takes a sentence and the start and end positions
of a query, and finds alternative words for that query. The queryStart index
should be the first byte of the query word, while the queryEnd index should be
the index of the first byte after the query word.
*/
func SentenceFindAlternativeWords(sentence string, queryStart, queryEnd, maxWords int) ([]string, error) {
  beforeString := sentence[:queryStart]
  afterString := sentence[queryEnd:]
  queryWord := sentence[queryStart:queryEnd]

  beforeWords := SplitWords(beforeString)
  afterWords := SplitWords(afterString)

  return FindAlternativeWords(beforeWords, afterWords, queryWord, maxWords)
}

/*
FindAlternativeWords takes a sentence and replaces the queryWord with a set of
words (not exceeding maxWords), which may be a good fit given the context. The
words before the queryWord in the sentence are given by beforeWords, and the
words after are given by afterWords.
*/
func FindAlternativeWords(beforeWords, afterWords []string, queryWord string, maxWords int) ([]string, error) {
  alternativeWords := make([]string, 0)

  // TODO: implement me
  return alternativeWords, nil
}

/*
AlternativeWords returns alternative words that can be used in place of the
current word. The words that are returned will have close meanings to the word
used as an argument, but which are usually used in place of that word.

The maxWords parameter specifies the maximum number of words to return.
*/
func AlternativeWords(word string, maxWords int) ([]string, error) {
  alternativeWords := make([]string, 0)
  alternativeWordVectors, err := AlternativeWordVectors(word, maxWords)
  if err != nil {
    return alternativeWords, err
  }

  for _, wordVector := range alternativeWordVectors {
    alternativeWords = append(alternativeWords, wordVector.Word)
  }

  return alternativeWords, nil
}

/*
AlternativeWordVectors returns alternative word vectors for a particular word.
*/
func AlternativeWordVectors(word string, maxWords int) ([]WordVector, error) {
  wordVectors := make(WordVectorCollection, 0)

  targetVectors, err := TargetVectors(word)
  if err != nil {
    return wordVectors, err
  }
  wordVectors = append(wordVectors, targetVectors...)

  synonyms, err := Synonyms(word)
  if err != nil {
    return wordVectors, err
  }

  for _, synonym := range synonyms {
    synonymVectors, err := TargetVectors(synonym)
    if err != nil {
      return wordVectors, err
    }
    wordVectors = append(wordVectors, synonymVectors...)
  }

  var wordsToSelect int
  if len(wordVectors) < maxWords {
    wordsToSelect = len(wordVectors)
  } else {
    wordsToSelect = maxWords
  }
  quickselect.QuickSelect(wordVectors, wordsToSelect)

  return wordVectors[wordsToSelect:], nil
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

    newVectors, err := paragraphTargetVectors(body, word)
    if err != nil {
      return wordVectors, err
    }
    wordVectors = append(wordVectors, newVectors...)
  }

  return wordVectors, nil
}

func paragraphTargetVectors(paragraph, word string) ([]WordVector, error) {
  wordVectors := make([]WordVector, 0)

  regexString := fmt.Sprintf(`[\s[[:punct:]]]*%s[\s[[:punct:]]]*`, word)
  compiledRegex, err := regexp.Compile(regexString)
  if err != nil {
    return wordVectors, err
  }

  allIndices := compiledRegex.FindAllStringIndex(paragraph, -1)
  var potentialVector string
  var currentVector WordVector

  if allIndices != nil {
    for i := 0; i < len(allIndices); i++ {
      matchedIndexPair := allIndices[i]
      start, end := matchedIndexPair[0], matchedIndexPair[1]

      potentialVector, err = lookForWord(paragraph[:start], false)
      if err != nil {
        return wordVectors, err
      }
      if potentialVector != "" {
        currentVector = WordVector{potentialVector, 1.0}
        wordVectors = append(wordVectors, currentVector)
      }

      potentialVector, err := lookForWord(paragraph[end:], true)
      if err != nil {
        return wordVectors, err
      }
      if potentialVector != "" {
        currentVector = WordVector{potentialVector, 1.0}
        wordVectors = append(wordVectors, currentVector)
      }
    }
  }

  return wordVectors, nil
}

func lookForWord(paragraph string, atBeginning bool) (string, error) {
  if strings.Trim(paragraph, " ") == "" {
    return "", nil
  }

  words := SplitWords(paragraph)
  if atBeginning {
    return words[0], nil
  } else {
    return words[len(words)-1], nil
  }
}

func Synonyms(word string) ([]string, error) {
  // TODO: implement
  return nil, nil
}
