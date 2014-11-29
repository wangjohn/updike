package dataingestor

import (
  "strings"
  "fmt"
  "os"
  "bufio"
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

  scanner := bufio.NewScanner(f)
  scanner.Split(bufio.ScanLines)

  pageBegin, err := regexp.Compile("<page>")
  if err != nil {
    return err
  }

  pageEnd, err := regexp.Compile("</page>")
  if err != nil {
    return err
  }

  var onPage bool
  currentPage := make([]string, 0)
  for scanner.Scan() {
    text := scanner.Text()
    if onPage {
      currentPage = append(currentPage, text)

      if pageEnd.MatchString(text) {
        onPage = false
        d.ingestWikipediaPage(strings.Join(currentPage, ""))
        break
      }
    } else if pageBegin.MatchString(text) {
      currentPage = make([]string, 1)
      currentPage[0] = text
      onPage = true
    }
  }

  return nil
}

type wikipediaStruct struct {
  Title string `xml:"title,attr"`
  ID int `xml:"id,attr"`
}

func (d DataIngestor) ingestWikipediaPage(page string) (error) {
  fmt.Println(page)
  stringReader := strings.NewReader(page)
  decoder := xml.NewDecoder(stringReader)

  fmt.Println(decoder.Decode(wikipediaStruct))

  return nil
}
