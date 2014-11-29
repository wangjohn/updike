package dataingestor

import (
  "fmt"
  "os"
  "encoding/xml"

  "github.com/wangjohn/updike/philarios"
)

const (
  readerBufferSize = 4096
)

type DataIngestor struct {
  Storage philarios.Storage
}

func (d DataIngestor) IngestWikipedia(filename string) (error) {
  f, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer f.Close()

  decoder := xml.NewDecoder(f)
  for {
    t, err := decoder.Token()
    if err != nil {
      return err
    }
    if t == nil {
      break
    }

    switch se := t.(type) {
    case xml.StartElement:
      if se.Name.Local == "page" {
        var p wikipediaPage
        decoder.DecodeElement(&p, &se)
        d.ingestWikipediaPage(p)
        return nil
      }
    }
  }

  return nil
}

type wikipediaPage struct {
  Title string `xml:"title"`
  ID int `xml:"id"`
  Revision wikipediaRevision `xml:"revision"`
}

type wikipediaRevision struct {
  Text string `xml:"text"`
}

func (d DataIngestor) ingestWikipediaPage(page wikipediaPage) (error) {
  fmt.Println(page.Revision.Text)

  return nil
}
