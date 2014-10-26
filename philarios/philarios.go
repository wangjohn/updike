package philarios

import (
  "regexp"
  "github.com/wangjohn/quickselect"

  "fmt"
  "strings"
)

const (
  NumSurroundingWords = 2
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

type WordFactory struct {
  Storage Storage
  TextProcessor TextProcessor
}

/*
SentenceFindAlternativeWords takes a sentence and the start and end positions
of a query, and finds alternative words for that query. The queryStart index
should be the first byte of the query word, while the queryEnd index should be
the index of the first byte after the query word.
*/
func (p WordFactory) SentenceFindAlternativeWords(sentence string, queryStart, queryEnd, maxWords int) ([]string, error) {
  beforeString := sentence[:queryStart]
  afterString := sentence[queryEnd:]
  queryWord := sentence[queryStart:queryEnd]

  beforeWords := SplitWords(beforeString)
  afterWords := SplitWords(afterString)

  return p.FindAlternativeWords(beforeWords, afterWords, queryWord, maxWords)
}

/*
FindAlternativeWords takes a sentence and replaces the queryWord with a set of
words (not exceeding maxWords), which may be a good fit given the context. The
words before the queryWord in the sentence are given by beforeWords, and the
words after are given by afterWords.
*/
func (p WordFactory) FindAlternativeWords(beforeWords, afterWords []string, queryWord string, maxWords int) ([]string, error) {
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
func (p WordFactory) AlternativeWords(word string, maxWords int) ([]string, error) {
  alternativeWords := make([]string, 0)
  alternativeWordVectors, err := p.AlternativeWordVectors(word, maxWords)
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
func (p WordFactory) AlternativeWordVectors(word string, maxWords int) ([]WordVector, error) {
  wordVectors := make(WordVectorCollection, 0)

  targetVectors, err := p.TargetVectors(word)
  if err != nil {
    return wordVectors, err
  }
  wordVectors = append(wordVectors, targetVectors...)

  synonyms, err := Synonyms(word)
  if err != nil {
    return wordVectors, err
  }

  for _, synonym := range synonyms {
    synonymVectors, err := p.TargetVectors(synonym)
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

func (p WordFactory) TargetVectors(word string) ([]WordVector, error) {
  var wordVectors []WordVector

  storage := PostgresStorage{
    "postgres", "user=philarios dbname=philarios sslmode=verify-full"}

  paragraphs, err := storage.QueryForWord(word, nil)
  if err != nil {
    return wordVectors, err
  }

  for _, paragraph := range paragraphs {
    newVectors, err := paragraphTargetVectors(paragraph.Body, word)
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
  var potentialWords []string
  var currentVector WordVector

  if allIndices != nil {
    for i := 0; i < len(allIndices); i++ {
      matchedIndexPair := allIndices[i]
      start, end := matchedIndexPair[0], matchedIndexPair[1]

      potentialWords = lookForWords(paragraph[:start], NumSurroundingWords, false)
      if len(potentialWords) > 0 {
        for _, word := range potentialWords {
          currentVector = WordVector{word, 1.0}
          wordVectors = append(wordVectors, currentVector)
        }
      }

      potentialWords = lookForWords(paragraph[end:], NumSurroundingWords, true)
      if len(potentialWords) > 0 {
        for _, word := range potentialWords {
          currentVector = WordVector{word, 1.0}
          wordVectors = append(wordVectors, currentVector)
        }
      }
    }
  }

  return wordVectors, nil
}

func lookForWords(paragraph string, wordsToCapture int, atBeginning bool) ([]string) {
  if strings.Trim(paragraph, " ") == "" {
    return []string{}
  }

  words := SplitWords(paragraph)
  var index int
  if len(words) > wordsToCapture {
    index = wordsToCapture
  } else {
    index = len(words)
  }

  if atBeginning {
    return words[:index]
  } else {
    return words[(len(words)-index):]
  }
}

func Synonyms(word string) ([]string, error) {
  // TODO: implement
  synonyms := make([]string, 0)
  return synonyms, nil
}
