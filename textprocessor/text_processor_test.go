package textprocessor

import (
  "testing"
  "reflect"
)

func TestProcessParagraphs(t *testing.T) {
  fixtures := []struct {
    Text string
    Expected []string
  }{
    {"This is some\nText that is broken\nup into multiple paragraphs",
      []string{"This is some", "Text that is broken", "up into multiple paragraphs"}},
    {"\n\nBlah",
      []string{"Blah"}},
    {"\n\nAlmost\nThere\n\nAreyou? Sure  \n",
      []string{"Almost", "There", "Areyou? Sure"}},
  }

  for _, fixture := range fixtures {
    paragraphs, err := ProcessParagraphs(fixture.Text)
    if err != nil {
      t.Errorf("Did not expect error for ProcessParagraphs: %s", err.Error())
    }

    if len(paragraphs) != len(fixture.Expected) {
      t.Errorf("Unexpected paragraph length. Expected %d, but obtained %d",
        len(fixture.Expected), len(paragraphs))
    }

    for i := 0; i < len(paragraphs); i++ {
      if !reflect.DeepEqual(paragraphs[i], fixture.Expected[i]) {
        t.Errorf("Unexpected paragraph. Expected '%s', but obtained '%s'",
          fixture.Expected[i], paragraphs[i])
      }
    }
  }
}

func TestNormalizedWord(t *testing.T) {
  fixtures := []struct {
    Word string
    Expected string
  }{
    {"Naming", "name"},
    {"boom-cam", "boom-cam"},
    {"safes", "safe"},
    {"headphones", "headphone"},
    {"Ryan's", "ryan"},
    {"CustomErs'", "customer"},
    {"zamboni", "zamboni"},
    {"facing", "face"},
    {"payed", "pay"},
    {"lying", "lie"},
    {"stopping", "stop"},
    {"stopped", "stop"},
    {"dancing", "dance"},
    {"danced", "dance"},
  }

  for _, fixture := range fixtures {
    result, err := NormalizedWord(fixture.Word)
    if err != nil {
      t.Errorf("Did not expect error for NormalizedWord: %s", err.Error())
    }

    if result != fixture.Expected {
      t.Errorf("Non-matching result. Expected '%s', but obtained '%s'",
        fixture.Expected, result)
    }
  }
}

/*
Testing for `ing` rules
*/

var textProcessorIngRules = []ProcessingRule{
  FilterBy("EndsWith", "ying").Try(Word{0,-4}, 'i', 'e'),
  FilterBy("EndsWith", Consonant{1}, Consonant{1}, "ing").Try(Word{0,-4}),
  FilterBy("EndsWith", "ing").Try(Word{0,-3}),
}

func handleSuffixIng(word string) (string, error) {
  for _, rule := range textProcessorIngRules {
    applies, normalizedRunes := rule.Rule(word)
    if applies {
      return string(normalizedRunes), nil
    }
  }

  return word, nil
}

func TestHandleSuffixIng(t *testing.T) {
  fixtures := []struct {
    Word string
    Expected string
  }{
    {"bling", "bling"},
    {"zing", "zing"},
    {"ding", "ding"},
    {"string", "string"},
    {"having", "have"},
    {"making", "make"},
    {"starring", "star"},
    {"stopping", "stop"},
    {"beginning", "begin"},
    {"lying", "lie"},
    {"dying", "die"},
    {"staying", "stay"},
  }

  for _, fixture := range fixtures {
    result, err := handleSuffixIng(fixture.Word)
    if err != nil {
      t.Errorf("Did not expect handleSuffixIng to throw an error: %s", err.Error())
    }

    if !reflect.DeepEqual(result, fixture.Expected) {
      t.Errorf("Did not obtain expected string. Expected '%q' received '%q'.",
        fixture.Expected, result)
    }
  }
}

/*
Testing for `ed` rules
*/

var textProcessorEdRules = []ProcessingRule{
  FilterBy("EndsWith", "ied").Try(Word{0,-1}),
  FilterBy("EndsWith", Consonant{1}, Consonant{1}, "ed").Try(Word{0,-3}),
  FilterBy("EndsWith", Consonant{1}, "ed").Try(Word{0,-1}),
  FilterBy("EndsWith", "ed").Try(Word{0,-2}),
}

func handleSuffixEd(word string) (string, error) {
  for _, rule := range textProcessorEdRules {
    applies, normalizedRunes := rule.Rule(word)
    if applies {
      return string(normalizedRunes), nil
    }
  }

  return word, nil
}

func TestHandleSuffixEd(t *testing.T) {
  fixtures := []struct {
    Word string
    Expected string
  }{
    {"stopped", "stop"},
    {"died", "die"},
    {"lied", "lie"},
    {"accoladed", "accolade"},
    {"compacted", "compact"},
    {"phoned", "phone"},
    {"danced", "dance"},
    {"bleed", "bleed"},
    {"zed", "zed"},
  }

  for _, fixture := range fixtures {
    result, err := handleSuffixIng(fixture.Word)
    if err != nil {
      t.Errorf("Did not expect handleSuffixIng to throw an error: %s", err.Error())
    }

    if !reflect.DeepEqual(result, fixture.Expected) {
      t.Errorf("Did not obtain expected string. Expected '%q' received '%q'.",
        fixture.Expected, result)
    }
  }
}

func TestGetRunesFromString(t *testing.T) {
  fixtures := []struct {
    String string
    ExpectedRune []rune
  }{
    {"string", []rune{'s', 't', 'r', 'i', 'n', 'g'}},
    {"hello", []rune{'h', 'e', 'l', 'l', 'o'}},
    {"dangalang", []rune{'d', 'a', 'n', 'g', 'a', 'l', 'a', 'n', 'g'}},
  }

  for _, fixture := range fixtures {
    result := getRunesFromString(fixture.String)
    if !reflect.DeepEqual(result, fixture.ExpectedRune) {
      t.Errorf("Obtained unexpected runes from string. Expected '%q' received '%q'.",
        fixture.ExpectedRune, result)
    }
  }
}
