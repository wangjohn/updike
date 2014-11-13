package tfidf

import (
  "database/sql"
  "testing"
)

const (
  testDriverName = "postgres"
  testDataSourceName = "host=localhost user=philarios dbname=philarios_tfidf_test sslmode=disable"
)

func setupDatabase() (PersistentTFIDF, error) {
  db, err := sql.Open(testDriverName, testDataSourceName)
  if err != nil {
    return nil, err
  }
  defer db.Close()

  tfidf := PersistentTFIDF{db}
  err = clearDatabase(db)
  if err != nil {
    return nil, err
  }

  return tfidf, nil
}

func clearDatabase(db *sql.DB) (error) {
  _, err = db.Exec(`
    DROP TABLE IF EXISTS word_document_pairs;
    DROP TABLE IF EXISTS document_frequency;
  `)
  return err
}
