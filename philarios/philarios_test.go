package philarios

import (
  "testing"
)

func TestForBeginningWordTargetVector(t *testing.T) {
  fixtures := []struct {
    paragraph string
    beginningVector string
  }{
    {"My booty don't lie", "My"},
    {"", ""},
    {"Tabitha", "Tabitha"},
    {"Nineteen. Ninety five.", "Nineteen"},
    {"Chain-smokers.", "Chain-smokers"},
    {"--bad things happen", "bad"},
  }

  for _, fixture := range fixtures {
    resultingWord, _ := lookForWord(fixture.paragraph, true)
    if resultingWord != fixture.beginningVector {
      t.Errorf("Did not find the expected vector. Expected '%s' but received '%s'",
        fixture.beginningVector, resultingWord)
    }
  }
}
