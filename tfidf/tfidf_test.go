package tfidf

import (
  "database/sql"
  "testing"
  "math"
)

const (
  testDriverName = "postgres"
  testDataSourceName = "host=localhost user=philarios dbname=philarios_tfidf_test sslmode=disable"

  floatEqualThresh = 0.00001
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

  err = tfidf.EnsureSchema()
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
    {"tango", 2, 0.984848485},
    {"blend", 1, 0.515},
    {"nonex", 1, 0.5},
    {"watev", 2, 0.5},
  }

  for _, fixture := range fixtures {
    resultTF, err := tfidf.TermFrequency(fixture.Word, fixture.DocumentId)
    if err != nil {
      t.Errorf("Should not have thrown an error for term frequency: err=%v", err)
    }

    if math.Abs(resultTF - fixture.ExpectedTF) > floatEqualThresh {
      t.Errorf("Received unexpected term frequency value: result=%v, expected=%v",
        resultTF, fixture.ExpectedTF)
    }
  }
}

func TestInverseDocumentFrequency(t *testing.T) {
  tfidf, db, err := setupDatabase()
  defer clearDatabase(db)
  if err != nil {
    t.Errorf("Should not have thrown an error while setting up database: err=%v", err)
  }

  _, err = db.Exec(`
    INSERT INTO word_document_pairs
    (word, freq, doc_max_word_freq, document) VALUES
    ('hello', 15, 43, 1),
    ('tango', 32, 33, 2),
    ('hello', 1, 50, 2),
    ('blend', 3, 100, 1);

    INSERT INTO document_frequency
    (word, unique_documents) VALUES
    ('hello', 2),
    ('tango', 1),
    ('blend', 1);
  `)

  if err != nil {
    t.Errorf("Should not have thrown an error while inserting test data into database: err=%v", err)
  }

  fixtures := []struct {
    Word string
    ExpectedIDF float64
  }{
    {"hello", -0.176091259},
    {"tango", 0.0},
    {"blend", 0.0},
    {"suggarrrr", 0.3010299956},
    {"fiveras", 0.3010299956},
    {"not in dict", 0.3010299956},
  }

  for _, fixture := range fixtures {
    resultIDF, err := tfidf.InverseDocumentFrequency(fixture.Word)
    if err != nil {
      t.Errorf("Should not have thrown an error for inverse document frequency: err=%v", err)
    }

    if math.Abs(resultIDF - fixture.ExpectedIDF) > floatEqualThresh {
      t.Errorf("Received unexpected term frequency value: result=%v, expected=%v",
        resultIDF, fixture.ExpectedIDF)
    }
  }
}

func TestTFIDF(t *testing.T) {
  tfidf, db, err := setupDatabase()
  defer clearDatabase(db)
  if err != nil {
    t.Errorf("Should not have thrown an error while setting up database: err=%v", err)
  }

  _, err = db.Exec(`
    INSERT INTO word_document_pairs
    (word, freq, doc_max_word_freq, document) VALUES
    ('hello', 15, 43, 1),
    ('tango', 32, 33, 2),
    ('hello', 1, 50, 2),
    ('blend', 3, 100, 1);

    INSERT INTO document_frequency
    (word, unique_documents) VALUES
    ('hello', 2),
    ('tango', 1),
    ('blend', 1);
  `)

  if err != nil {
    t.Errorf("Should not have thrown an error while inserting test data into database: err=%v", err)
  }

  fixtures := []struct {
    Word string
    DocumentId int
    ExpectedScore float64
  }{
    {"hello", 1, -0.118759221},
    {"hello", 2, -0.089806542},
    {"tango", 2, 0.0},
    {"blend", 2, 0.0},
    {"notexistent", 1, 0.150514998},
    {"never existed before", 1, 0.150514998},
  }

  for _, fixture := range fixtures {
    resultScore, err := tfidf.TFIDFScore(fixture.Word, fixture.DocumentId)
    if err != nil {
      t.Errorf("Should not have thrown an error for inverse document frequency: err=%v", err)
    }

    if math.Abs(resultScore - fixture.ExpectedScore) > floatEqualThresh {
      t.Errorf("Received unexpected TFIDF value: result=%v, expected=%v",
        resultScore, fixture.ExpectedScore)
    }
  }
}

func TestStoreWord(t *testing.T) {
  tfidf, db, err := setupDatabase()
  defer clearDatabase(db)
  if err != nil {
    t.Errorf("Should not have thrown an error while setting up database: err=%v", err)
  }

  // Add words through the Store function
  storageFixtures := []struct {
    Word string
    Occurrences int
    DocMaxWordOccurrences int
    DocumentId int
  }{
    {"hello", 15, 43, 1},
    {"tango", 32, 33, 2},
    {"hello", 1, 50, 2},
    {"blend", 3, 100, 1},
  }

  for _, f := range storageFixtures {
    err = tfidf.Store(f.Word, f.Occurrences, f.DocMaxWordOccurrences, f.DocumentId)
    if err != nil {
      t.Errorf("Obtained an error while trying to store words: err=%v", err)
    }
  }

  // Check to make sure that the stored words have their TFIDF scores computed
  // correctly.
  scoreFixtures := []struct {
    Word string
    DocumentId int
    ExpectedTF float64
    ExpectedIDF float64
    ExpectedScore float64
  }{
    {"hello", 1, 0.674418605, -0.176091259, -0.118759221},
    {"hello", 2, 0.51, -0.176091259, -0.089806542},
    {"tango", 2, 0.984848485, 0.0, 0.0},
    {"blend", 1, 0.515, 0.0, 0.0},
    {"notexistent", 1, 0.5, 0.3010299956, 0.150514998},
    {"never existed before", 1, 0.5, 0.3010299956, 0.150514998},
  }

  for _, f := range scoreFixtures {
    // Check the Term Frequency score
    tfScore, err := tfidf.TermFrequency(f.Word, f.DocumentId)
    if err != nil {
      t.Errorf("Obtained an error while trying to get Term Frequency: err=%v", err)
    }
    if math.Abs(tfScore - f.ExpectedTF) > floatEqualThresh {
      t.Errorf("Received unexpected TF value: word=%v, result=%v, expected=%v",
        f.Word, tfScore, f.ExpectedTF)
    }

    // Check the Inverse Document Frequency score
    idfScore, err := tfidf.InverseDocumentFrequency(f.Word)
    if err != nil {
      t.Errorf("Obtained an error while trying to get Term Frequency: err=%v", err)
    }
    if math.Abs(idfScore - f.ExpectedIDF) > floatEqualThresh {
      t.Errorf("Received unexpected IDF value: word=%v, result=%v, expected=%v",
        f.Word, idfScore, f.ExpectedIDF)
    }

    // Check the TFIDF score
    score, err := tfidf.TFIDFScore(f.Word, f.DocumentId)
    if err != nil {
      t.Errorf("Obtained an error while trying to get TFIDF score: err=%v", err)
    }
    if math.Abs(score - f.ExpectedScore) > floatEqualThresh {
      t.Errorf("Received unexpected TFIDF value: word=%v, result=%v, expected=%v",
        f.Word, score, f.ExpectedScore)
    }
  }
}
