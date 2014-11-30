package dataingestor

import (
  "fmt"
  "os"
  "regexp"
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
  pagesIngested := 0
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
        pagesIngested++
      }
    }

    if pagesIngested == 2 {
      return nil
    }
  }

  return nil
}

type wikipediaPage struct {
  Title string `xml:"title"`
  ID int `xml:"id"`
  Revision wikipediaRevision `xml:"revision"`
}

func (w wikipediaPage) SourceURL() (string) {
  reg, _ := regexp.Compile("[[:space:]]")
  underscored := reg.ReplaceAllString(w.Title, "_")
  return "en.wikipedia.org/" + underscored
}

func (w wikipediaRevision) Date() (string) {
  reg, _ := regexp.Compile("[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}")
  return reg.FindString(w.Timestamp)
}

type wikipediaRevision struct {
  ID int `xml:"id"`
  Timestamp string `xml:"timestamp"`
  Contributor wikipediaContributor `xml:"contributor"`
  Text string `xml:"text"`
  Format string `xml:"format"`
}

type wikipediaContributor struct {
  Username string `xml:"username"`
  ID int `xml:"id"`
}

func (d DataIngestor) ingestWikipediaPage(page wikipediaPage) (error) {
  reg, err := regexp.Compile("({{.*}}|\\[\\[.*\\|.*\\]\\])")
  if err != nil {
    return err
  }

  body := reg.ReplaceAllString(page.Revision.Text, "")
  pub := philarios.Publication{
    Title: page.Title,
    Author: "wikipedia",
    Editor: "",
    Date: page.Revision.Date(),
    SourceID: page.ID,
    SourceURL: page.SourceURL(),
    Encoding: "utf-8",
    Type: "wikipedia_article",
    Categories: []string{},
    Text: body,
  }
  fmt.Println(pub)

  return nil
}

