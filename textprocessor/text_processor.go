package textprocessor

import (
  "log"
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
    endIndex = len(wordRunes) - w.EndIndex + 1
  } else {
    endIndex = w.EndIndex
  }

  if endIndex < 0 {
    log.Fatal(`Specified an EndIndex on a word that cannot exist. You specified
               '%d', but the word only has length %d`, w.EndIndex, len(wordRunes))
  }

  if startIndex > endIndex {
    log.Fatal(`Specified a StartIndex '%d' on a word that is greater than the
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
  return processingResult.FilterBy(filterType, args)
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
      if validInt {
        log.Fatal("Must use an integer as the argument for 'LongerThan' filter type.")
      }
      return len(wordRunes) > length
    }
  case "EndsWith":
    currentFilterFunction = func(word string) (bool) {
      wordRunes := getRunesFromString(word)

      // Not enough characters in word, so it can't possibly end with the arguments provided
      if len(wordRunes) < len(args) {
        return false
      }

      consonantMapping := make(map[int][]int)
      vowelMapping := make(map[int][]int)
      for i, argument := range args {
        wordIndex := len(wordRunes) - i - 1

        switch arg := argument.(type) {
        case Consonant:
          consonantMapping[arg.Index] = append(consonantMapping[arg.Index], wordIndex)
        case Vowel:
          vowelMapping[arg.Index] = append(vowelMapping[arg.Index], wordIndex)
        case rune:
          if wordRunes[wordIndex] != arg {
            return false
          }
        default:
          log.Fatal("Invalid argument type '%T' for 'EndsWith' filter type.", argument)
        }
      }

      for consonantIndex := range consonantMapping {
        var requiredLetter rune
        for j, wordIndex := range consonantMapping[consonantIndex] {
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
    log.Fatal(`Specified an invalid filter type '%s'.`, filterType)
  }

  var filterFunctions []FilterFunction
  filterFunctions = append(filterFunctions, d.FilterFunctions...)
  filterFunctions = append(filterFunctions, currentFilterFunction)

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
