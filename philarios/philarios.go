package philarios

import (
  "github.com/wangjohn/quickselect"
  "github.com/wangjohn/updike/tfidf"
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
  Settings Settings
  TFIDF tfidf.TFIDF
}

type Philarios interface {
  SentenceFindAlternativeWords(sentence string, queryStart, queryEnd, maxWords int) ([]string, error)
  FindAlternativeWords(beforeWords, afterWords []string, queryWord string, maxWords int) ([]string, error)
  AlternativeWords(word string, maxWords int) ([]string, error)
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

  beforeWords = p.findImportantWords(beforeWords)
  afterWords = p.findImportantWords(afterWords)

  // TODO: figure out how we want to include categories. This is a stop gap.
  categories := []string{}

  for _, beforeWord := range beforeWords {
    _, err := p.Storage.QueryForWord(beforeWord, categories)
    if err != nil {
      return alternativeWords, err
    }

  }

  // TODO: do after words.

  return alternativeWords, nil
}

func (p WordFactory) findImportantWords(words []string) ([]string) {
  // TODO: for now, returning all words. In the future, we probably want to filter
  // by words with high aggregate TFIDF scores.
  return words
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
  paragraphs, err := p.Storage.QueryForWord(word, nil)
  if err != nil {
    return nil, err
  }

  scoreCollection := make(map[string]float64)
  for _, paragraph := range paragraphs {
    probWordVectors, err := p.associatedWordProbabilities(paragraph.Body, word)
    if err != nil {
      return nil, err
    }

    for _, vec := range probWordVectors {
      scoreCollection[vec.Word] += (vec.Score / float64(len(paragraphs)))
    }
  }

  var wordVectors = make([]WordVector, len(scoreCollection))
  i := 0
  for key, value := range scoreCollection {
    wordVectors[i] = WordVector{key, value}
    i++
  }
  return wordVectors, nil
}

func (p WordFactory) associatedWordProbabilities(paragraph, word string) ([]WordVector, error) {
  associatedCounts := make(map[string]int)

  paragraphWords := SplitWords(paragraph)
  wordOccurrences := 0
  for i, pw := range paragraphWords {
    if FuzzyStringEquals(pw, word) {
      for _, surroundingWord := range p.surroundingWords(paragraphWords, i) {
        canonicalSW := CanonicalWordForm(surroundingWord)
        associatedCounts[canonicalSW]++
      }
      wordOccurrences++
    }
  }

  wordVectors := make([]WordVector, len(associatedCounts))
  i := 0
  for key, value := range associatedCounts {
    occurrenceProb := float64(value) / float64(wordOccurrences)
    wordVectors[i] = WordVector{key, occurrenceProb}
    i++
  }

  return wordVectors, nil
}

func (p WordFactory) surroundingWords(words []string, wordIndex int) ([]string) {
  var start, end int
  if wordIndex > p.Settings.WordsToCapture {
    start = wordIndex - p.Settings.WordsToCapture
  } else {
    start = 0
  }

  if p.Settings.WordsToCapture + wordIndex < len(words) {
    end = wordIndex + p.Settings.WordsToCapture
  } else {
    end = len(words)
  }

  surrounding := make([]string, end - start - 1)

  j := 0
  for i := start; i < end; i++ {
    if i != wordIndex {
      surrounding[j] = words[i]
      j++
    }
  }

  return surrounding
}

func Synonyms(word string) ([]string, error) {
  // TODO: implement
  synonyms := make([]string, 0)
  return synonyms, nil
}
