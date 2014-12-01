package main

import (
  "database/sql"
  "github.com/wangjohn/updike/philarios"
  "github.com/wangjohn/updike/tfidf"
  "github.com/wangjohn/updike/dataingestor"
)

const (
  tfidfDriverName = "postgres"
  tfidfDataSourceName = "host=localhost user=philarios dbname=philarios_tfidf sslmode=disable"

  storageDriverName = "postgres"
  storageDataSourceName = "host=localhost user=philarios dbname=philarios_storage sslmode=disable"
)

func main() {
  wordFactory, _ := createWordFactory()

  ingestor := dataingestor.DataIngestor{wordFactory.Storage}
  ingestor.IngestWikipedia("/home/wangjohn/wikipedia/enwiki-latest-pages-articles.xml")
}

func createWordFactory() (*philarios.WordFactory, error) {
  storageDb, err := sql.Open(storageDriverName, storageDataSourceName)
  if err != nil {
    return nil, err
  }

  storage := philarios.PostgresStorage{storageDb}
  settings := philarios.DefaultSettingsObject()
  tfidfDb, err := sql.Open(tfidfDriverName, tfidfDataSourceName)
  if err != nil {
    return nil, err
  }

  tfidf := tfidf.PersistentTFIDF{tfidfDb}
  wordFactory := philarios.WordFactory{storage, settings, tfidf}
  return &wordFactory, nil
}
