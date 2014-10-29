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

func TestHandleSuffixIng(t *testing.T) {
  fixtures := []struct {
    Word string
    Expected string
    ExpectedBool bool
  }{
    {"bling", "", false},
    {"zing", "", false},
    {"ding", "", false},
    {"string", "", false},
    {"having", "have", true},
    {"making", "make", true},
    {"starring", "star", true},
    {"stopping", "stop", true},
    {"beginning", "begin", true},
    {"lying", "lie", true},
    {"dying", "die", true},
    {"staying", "stay", true},
  }

  for _, fixture := range fixtures {
    inputRune := getRunesFromString(fixture.Word)
    expectedRune := getRunesFromString(fixture.Expected)

    boolResult, result := handleSuffixIng(inputRune)
    if boolResult != fixture.ExpectedBool {
      t.Errorf("Did not expect boolean result. Expected '%t' received '%t'.",
        fixture.ExpectedBool, boolResult)
    }

    if !reflect.DeepEqual(result, expectedRune) {
      t.Errorf("Did not obtain expected string. Expected '%q' received '%q'.",
        expectedRune, result)
    }
  }
}

func TestHandleSuffixEd(t *testing.T) {
  fixtures := []struct {
    Word string
    Expected string
    ExpectedBool bool
  }{
    {"stopped", "stop", true},
    {"died", "die", true},
    {"lied", "lie", true},
    {"accoladed", "accolade", true},
    {"compacted", "compact", true},
    {"phoned", "phone", true},
    {"danced", "dance", true},
    {"bleed", "", false},
    {"zed", "", false},
  }

  for _, fixture := range fixtures {
    inputRune := getRunesFromString(fixture.Word)
    expectedRune := getRunesFromString(fixture.Expected)

    boolResult, result := handleSuffixEd(inputRune)
    if boolResult != fixture.ExpectedBool {
      t.Errorf("Did not expect boolean result. Expected '%t' received '%t'.",
        fixture.ExpectedBool, boolResult)
    }

    if !reflect.DeepEqual(result, expectedRune) {
      t.Errorf("Did not obtain expected string. Expected '%q' received '%q'.",
        expectedRune, result)
    }
  }
}

func TestImproperWords(t *testing.T) {
  fixtures := []struct {
    Word string
    Expected string
    ExpectedBool bool
  }{
    {"began", "begin", true},
    {"mice", "mouse", true},
  }

  for _, fixture := range fixtures {
    inputRune := getRunesFromString(fixture.Word)
    expectedRune := getRunesFromString(fixture.Expected)

    boolResult, result := handleImproperWord(inputRune)
    if boolResult != fixture.ExpectedBool {
      t.Errorf("Did not expect boolean result. Expected '%t' received '%t'.",
        fixture.ExpectedBool, boolResult)
    }

    if !reflect.DeepEqual(result, expectedRune) {
      t.Errorf("Did not obtain expected string. Expected '%q' received '%q'.",
        expectedRune, result)
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
