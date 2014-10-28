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
