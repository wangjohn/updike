package tfidf

import (
  "database/sql"
  "testing"
)

const (
  testDriverName = "postgres"
  testDataSourceName = "host=localhost user=philarios dbname=philarios_tfidf_test sslmode=disable"
)

func setupDatabase() (*PersistentTFIDF, *sql.DB, error) {
  db, err := sql.Open(testDriverName, testDataSourceName)
  if err != nil {
    return nil, nil, err
  }

  tfidf := PersistentTFIDF{db}
  err = clearDatabase(db)
  if err != nil {
    return nil, nil, err
  }

  return &tfidf, db, nil
}

func clearDatabase(db *sql.DB) (error) {
  _, err := db.Exec(`
    DROP TABLE IF EXISTS word_document_pairs;
    DROP TABLE IF EXISTS document_frequency;
  `)
  return err
}

func TestTermFrequency(t *testing.T) {
  tfidf, db, err := setupDatabase()
  defer clearDatabase(db)
  if err != nil {
    t.Errorf("Should not have thrown an error while setting up database: err=%v", err)
  }
  err = tfidf.EnsureSchema()
  if err != nil {
    t.Errorf("Should not have thrown an error while ensuring schema: err=%v", err)
  }

  _, err = db.Exec(`INSERT INTO word_document_pairs
    (word, freq, doc_max_word_freq, document) VALUES
    ('hello', 15, 43, 1),
    ('tango', 32, 33, 2),
    ('hello', 1, 50, 2),
    ('blend', 3, 100, 1);`)

  if err != nil {
    t.Errorf("Should not have thrown an error while inserting test data into database: err=%v", err)
  }

  fixtures := []struct {
    Word string
    DocumentId int
    ExpectedTF float64
  }{
    {"hello", 1, 0.674418605},
    {"hello", 2, 0.51},
    {"hello", 3, 0.5},
    {"tango", 2, 0.984848485},
    {"blend", 1, 0.515},
    {"blahd", 0, 0.5},
    {"nonex", 1, 0.5},
    {"watev", 2, 0.5},
  }

  for _, fixture := range fixtures {
    resultTF, err := tfidf.TermFrequency(fixture.Word, fixture.DocumentId)
    if err != nil {
      t.Errorf("Should not have thrown an error for term frequency: err=%v", err)
    }

    if resultTF != fixture.ExpectedTF {
      t.Error("Received unexpected term frequency value: result=%v, expected=%v",
        resultTF, fixture.ExpectedTF)
    }
  }
}
