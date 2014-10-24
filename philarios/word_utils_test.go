package philarios

import (
  "testing"
)

func TestSplitWords(t *testing.T) {
  fixtures := []struct {
    sentence string
    expectedWords []string
  }{
    {"Hello my name is John", []string{"Hello", "my", "name", "is", "John"}},
    {"A...B.'c'd?''", []string{"A", "B", "c", "d"}},
    {"123--hello-my.,23", []string{"123", "hello", "my", "23"}},
    {"Nick+emily--just/my/type", []string{"Nick", "emily", "just", "my", "type"}},
  }

  for _, fixture := range fixtures {
    words := SplitWords(fixture.sentence)
    if words != fixture.expectedWords {
      t.Errorf("Did not obtain the expected words. Expected %s but obtained %s",
        fixture.expectedWords, words)
    }
  }
}
