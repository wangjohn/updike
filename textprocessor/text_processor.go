package textprocessor

import (
  "log"
  "strings"
  "unicode/utf8"
)

var vowels = map[rune]bool {
  'a': true,
  'e': true,
  'i': true,
  'o': true,
  'u': true,
  'y': true,
}

var TextProcessorRules = []ProcessingRule{
  FilterBy("EndsWith", Consonant{1}, Consonant{1}, "ing").FilterBy("LongerThan", 5).Try(Word{0,-4}),
  FilterBy("EndsWith", "ying").FilterBy("LongerThan", 5).Try(Word{0,-4}, 'i', 'e'),
  FilterBy("EndsWith", "ing").FilterBy("LongerThan", 5).Try(Word{0,-3}),

  FilterBy("EndsWith", "ied").Try(Word{0,-1}),
  FilterBy("EndsWith", Consonant{1}, Consonant{1}, "ed").Try(Word{0,-3}),
  FilterBy("EndsWith", Consonant{1}, "ed").Try(Word{0,-1}),
  FilterBy("EndsWith", "ed").Try(Word{0,-2}),
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

  preprocessedParagraphs := []string{preprocessedText}
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

type Vowel struct {
  Index int
}

type Consonant struct {
  Index int
}

func IsVowel(char rune) (bool) {
  return vowels[char]
}

type Word struct {
  StartIndex int
  EndIndex int
}

func (w Word) GetSlice(word string) ([]rune) {
  wordRunes := getRunesFromString(word)
  startIndex := w.StartIndex

  var endIndex int
  if w.EndIndex < 0 {
    endIndex = len(wordRunes) + w.EndIndex
  } else {
    endIndex = w.EndIndex
  }

  if endIndex < 0 {
    log.Fatalf(`Specified an EndIndex on a word that cannot exist. You specified
               '%d', but the word only has length %d`, w.EndIndex, len(wordRunes))
  }

  if startIndex > endIndex {
    log.Fatalf(`Specified a StartIndex '%d' on a word that is greater than the
              EndIndex '%d' (for word of length '%d').`,
              w.StartIndex, w.EndIndex, len(wordRunes))
  }

  return wordRunes[startIndex:endIndex]
}


type FilterFunction func(word string) (bool)

/*
ProcessingRule is a rule for processing a string. It contains all the information
needed to determine whether a string should be processed, and what the resulting
normalized word will be.
*/
type ProcessingRule struct {
  Rule func(word string) (bool, []rune)
}

/*
IntermediateProcessingResult is an intermediate object used by the text processor
for chaining together declarations of a processing rule.
*/
type IntermediateProcessingResult interface {
  FilterBy(filterType string, args ...interface{}) (IntermediateProcessingResult)
  Try(args ...interface{}) (ProcessingRule)
}

/*
DefaultIntermediateProcessingResult is the default implementation of the
IntermediateProcessingResult interface. It carries everything necessary to
chain together a set of declarations to create a ProcessingRule.
*/
type DefaultIntermediateProcessingResult struct {
  FilterFunctions []FilterFunction
}

/*
FilterBy is a convenience method for creating an IntermediateProcessingResult and
calling FilterBy on the resulting struct.
*/
func FilterBy(filterType string, args ...interface{}) (IntermediateProcessingResult) {
  filterFunctions := []FilterFunction{}
  processingResult := DefaultIntermediateProcessingResult{filterFunctions}
  return processingResult.FilterBy(filterType, args...)
}

/*
FilterBy defines a declaration for filtering words by a particular set of instructions.
The filterType can be one of the following:
  - "LongerThan"(integer): Filters words to be longer than a certain integer
  - "EndsWith"(Consonants/runes): Filters words to end with a certain set of consonants
    or rune characters.
*/
func (d DefaultIntermediateProcessingResult) FilterBy(filterType string, args ...interface{}) (IntermediateProcessingResult) {
  var currentFilterFunction FilterFunction

  switch filterType {
  case "LongerThan":
    currentFilterFunction = func(word string) (bool) {
      wordRunes := getRunesFromString(word)
      if len(args) != 1 {
        log.Fatal("Must use a single argument for 'LongerThan' filter type.")
      }
      length, validInt := args[0].(int)
      if !validInt {
        log.Fatal("Must use an integer as the argument for 'LongerThan' filter type.")
      }
      return len(wordRunes) > length
    }
  case "EndsWith":
    // Convert all string arguments into individual rune arguments
    adjustedArgs := make([]interface{}, 0)
    for _, arg := range args {
      stringArg, isValidString := arg.(string)
      if isValidString {
        stringRunes := getRunesFromString(stringArg)
        for _, stringRune := range stringRunes {
          adjustedArgs = append(adjustedArgs, interface{}(stringRune))
        }
      } else {
        adjustedArgs = append(adjustedArgs, arg)
      }
    }
    args = adjustedArgs

    currentFilterFunction = func(word string) (bool) {
      wordRunes := getRunesFromString(word)

      // Not enough characters in word, so it can't possibly end with the arguments provided
      if len(wordRunes) < len(args) {
        return false
      }

      consonantMapping := make(map[int][]int)
      vowelMapping := make(map[int][]int)
      for i, argument := range args {
        wordIndex := len(wordRunes) - len(args) + i

        switch arg := argument.(type) {
        case rune:
          if wordRunes[wordIndex] != arg {
            return false
          }
        case Consonant:
          consonantMapping[arg.Index] = append(consonantMapping[arg.Index], wordIndex)
        case Vowel:
          vowelMapping[arg.Index] = append(vowelMapping[arg.Index], wordIndex)
        default:
          log.Printf("Invalid argument type '%T' for 'EndsWith' filter type. Argument is: %s", argument, argument)
        }
      }

      for consonantIndex := range consonantMapping {
        var requiredLetter rune
        for j, wordIndex := range consonantMapping[consonantIndex] {
          currentLetter := wordRunes[wordIndex]
          if IsVowel(currentLetter) {
            return false
          }
          if j == 0 {
            requiredLetter = currentLetter
          } else if currentLetter != requiredLetter {
            return false
          }
        }
      }

      for vowelIndex := range vowelMapping {
        var requiredLetter rune
        for j, wordIndex := range vowelMapping[vowelIndex] {
          currentLetter := wordRunes[wordIndex]
          if !IsVowel(currentLetter) {
            return false
          }
          if j == 0 {
            requiredLetter = currentLetter
          } else if currentLetter != requiredLetter {
            return false
          }
        }
      }

      return true
    }
  default:
    log.Fatalf(`Specified an invalid filter type '%s'.`, filterType)
  }

  d.FilterFunctions = append(d.FilterFunctions, currentFilterFunction)

  return d
}

func (d DefaultIntermediateProcessingResult) Try(args ...interface{}) (ProcessingRule) {
  processFunc := func (word string) (bool, []rune) {
    for _, filterFunc := range d.FilterFunctions {
      if !filterFunc(word) {
        return false, []rune{}
      }
    }

    runeResult := make([]rune, 0)
    for _, argument := range args {
      switch argument.(type) {
      case Word:
        wordArgument, _ := argument.(Word)
        currentRunes := wordArgument.GetSlice(word)
        runeResult = append(runeResult, currentRunes...)
      case rune:
        runeArgument, _ := argument.(rune)
        runeResult = append(runeResult, runeArgument)
      default:
        log.Fatal(`Used an unknown argument type in the 'Try' method. Argument
                  '%q' with type '%T' cannot but used.`, argument, argument)
      }
    }
    return true, runeResult
  }

  return ProcessingRule{processFunc}
}

func getRunesFromString(str string) ([]rune) {
  runes := make([]rune, 0)
  for i, w := 0, 0; i < len(str); i += w {
    runeValue, width := utf8.DecodeRuneInString(str[i:])
    runes = append(runes, runeValue)
    w = width
  }

  return runes
}

