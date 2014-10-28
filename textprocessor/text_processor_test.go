package textprocessor

import (
  "testing"
  "reflect"
  "github.com/wangjohn/updike/philarios"
)

func TestProcessParagraphs(t *testing.T) {
  fixtures := []struct {
    Text string
    PublicationId int
    Expected []string
  }{
    {"This is some\nText that is broken\nup into multiple paragraphs", 1,
      []string{"This is some", "Text that is broken", "up into multiple paragraphs"}},
    {"\n\nBlah", 2,
      []string{"Blah"}},
    {"\n\nAlmost\nThere\n\nAreyou? Sure  \n", 3,
      []string{"Almost", "There", "Areyou? Sure"}},
  }

  for _, fixture := range fixtures {
    publication := philarios.Publication{
      Title: "Title",
      Author: "Author",
      Editor: "Editor",
      Date: "Date",
      SourceURL: "SourceURL",
      Encoding: "Encoding",
      Text: fixture.Text,
      Categories: []string{"Category1"},
    }

    paragraphs, error := ProcessParagraphs(publication, fixture.publicationId)
    if len(paragraphs) != len(fixtures.Expected) {
      t.Errorf("Unexpected paragraph length. Expected %d, but obtained %d",
        len(fixtures.Expected), len(paragraphs))
    }

    for i := 0; i < len(paragraphs); i++ {
      if !reflect.DeepEquals(paragraphs[i], fixtures.Expected[i]) {
        t.Errorf("Unexpected paragraph. Expected '%s', but obtained '%s'",
          fixtures.Expected[i], paragraphs[i])
      }
    }
  }
}
