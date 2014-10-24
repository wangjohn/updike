package philarios

import (
  "reflect"
  "testing"
)

func TestLookForWords(t *testing.T) {
  fixtures := []struct {
    wordsToCapture int
    atBeginning bool
    paragraph string
    expectedWords []string
  }{
    // Test looking for words at the beginning of a paragraph
    {1, true, "My booty don't lie", []string{"My"}},
    {1, true, "", []string{}},
    {1, true, "Tabitha", []string{"Tabitha"}},
    {1, true, "Nineteen. Ninety five.", []string{"Nineteen"}},
    {1, true, "Chain-smokers.", []string{"Chain"}},
    {1, true, "--bad things happen", []string{"bad"}},

    // Test looking for words at the end of a paragraph
    {1, false, "Jimbo fisher", []string{"fisher"}},
    {1, false, "", []string{}},
    {1, false, "Onething", []string{"Onething"}},
    {1, false, "Twent-time..five/four/three/3/two...", []string{"two"}},
    {1, false, "staph-------", []string{"staph"}},

    // Test looking for multiple words at the beginning
    {3, true, "My booty don't", []string{"My", "booty", "don't"}},
    {5, true, "Just stop", []string{"Just", "stop"}},
    {0, true, "blah blah", []string{}},
    {3, true, "into...dusk?", []string{"into", "dusk"}},

    // Test looking for multiple words at the end
    {3, false, "Just stop your mouth", []string{"stop", "your", "mouth"}},
    {4, false, "Battle of st. Crispin's day", []string{"of", "st", "Crispin's", "day"}},
    {10, false, "don't move", []string{"don't", "move"}},
    {5, false, "hahah/you/url", []string{"hahah", "you", "url"}},
  }

  for _, fixture := range fixtures {
    resultingWords := lookForWords(fixture.paragraph, fixture.wordsToCapture, fixture.atBeginning)
    if !reflect.DeepEqual(resultingWords, fixture.expectedWords) {
      t.Errorf("Did not find the expected words. Expected '%s' but received '%s'",
        fixture.expectedWords, resultingWords)
    }
  }
}
