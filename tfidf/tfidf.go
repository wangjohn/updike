package tfidf

import (
  "database/sql"
  "math"
)

type TFIDF interface {
  Store(word string) (error)
  TermFrequency(word string) (float64, error)
  InverseDocumentFrequency(word string) (float64, error)
  TFIDFScore(word string) (float64, error)
}

type PersistentTFIDF struct {
  SQLDatabase *sql.DB
}

var persistentSqlSchema = `
CREATE TABLE IF NOT EXISTS tfidf (
  id bigserial PRIMARY KEY,
  word text,
  occurrences integer,
  max_occurrences integer,
  unique_documents integer
);

CREATE TABLE IF NOT EXISTS tfidf_pairs (
  id bigserial PRIMARY KEY,
  word text,
  occurrences integer,
  document integer
)
`

var totalDocuments = -1

func (p PersistentTFIDF) TermFrequency(word string) (float64, error) {
  var tf float64
  err := p.SQLDatabase.QueryRow(
    `SELECT occurrences FROM tfidf
     WHERE word=?`, word).Scan(&tf)
  if err != nil {
    return 0.0, err
  }

  return tf, nil
}

func (p PersistentTFIDF) InverseDocumentFrequency(word string) (float64, error) {
  var idf float64
  err := p.SQLDatabase.QueryRow(
    `SELECT unique_documents FROM tfidf
    WHERE word=?`, word).Scan(&idf)
  if err != nil {
    return 0.0, err
  }

  if totalDocuments == -1 {
    err = p.SQLDatabase.QueryRow(
      `SELECT COUNT(DISTINCT document) FROM tfidf_paris`).Scan(&totalDocuments)
    if err != nil {
      return 0.0, err
    }
  }

  return math.Log10(totalDocuments / (1 + idf)), nil
}

func (p PersistentTFIDF) TFIDFScore(word string) (float64, error) {
  tf, err := p.TermFrequency(word)
  if err != nil {
    return 0.0, err
  }

  idf, err := p.InverseDocumentFrequency(word)
  if err != nil {
    return 0.0, err
  }

  return tf * idf, nil
}

func (p PersistentTFIDF) Store(word string, occurrences, documentId int) (error) {
  var isNewDocument, isNewWord bool
  var pairId int
  err := p.SQLDatabase.QueryRow(
   `SELECT id FROM tfidf_pairs
    WHERE word='?'
    AND document=?`, word, documentId).Scan(&pairId)

  if err != nil {
    return err
  }

  if err == sql.ErrNoRows {
    isNewDocument = true
    insertPairQuery := fmt.Sprintf(
     `INSERT INTO tfidf_pairs(
        word, occurrences, document)
      VALUES ('%s', %d, %d)
      RETURNING id`,
      word,
      occurrences,
      documentId)
    err = p.SQLDatabase.QueryRow(insertPairQuery).Scan(&pairId)
    if err != nil {
      return err
    }
  } else {
    isNewDocument = false
    _, err := p.SQLDatabase.Exec(
     `UPDATE tfidf_pairs
      SET occurrences=?
      WHERE word='?'
      AND document=?`,
      occurrences,
      word,
      documentId)
    if err != nil {
      return err
    }
  }

  var tfidfId, curOccurrences, curMaxOccurrences, curUniqueDocuments int
  err = p.SQLDatabase.QueryRow(
   `SELECT id, occurrences, max_occurrences, unique_documents FROM tfidf
    WHERE word='?'`, word)
    .Scan(&tfidfIdf, &curOccurrences, &curMaxOccurrences, &curUniqueDocuments)

  if err != nil {
    return err
  }
  isNewWord = (err == sql.ErrNoRows)

  var adjOccurrences, adjMaxOccurrences, adjUniqueDocuments int
  if isNewDocument {
    adjOccurrences = curOccurrences + occurrences
    adjMaxOccurrences = max(occurrences, curMaxOccurrences)
    adjUniqueDocuments = curUniqueDocuments + 1
  } else {
    err = p.SQLDatabase.QueryRow(
     `SELECT SUM(occurrences), MAX(occurrences), COUNT(*)
      FROM tfidf_pairs
      WHERE word='?'`,
      word).Scan(&adjOccurrences, adjMaxOccurrences, adjUniqueDocuments)

    if err != nil {
      return err
    }
  }

  if isNewWord {
    insertTfidfQuery := fmt.Sprintf(
     `INSERT INTO tfidf(
        word, occurrences, max_occurrences, unique_documents)
      VALUES ('%s', %d, %d, %d)`,
      word,
      adjOccurrences,
      adjMaxOccurrences,
      adjUniqueDocuments)
    err = p.SQLDatabase.QueryRow(insertTfidfQuery)
  } else {
    err = p.SQLDatabase.Exec(
     `UPDATE tfidf
      SET occurrences=?, max_occurrences=?, unique_documents=?
      WHERE word='?'`,
      adjOccurrences,
      adjMaxOccurrences,
      adjUniqueDocuments,
      word)
  }

  if err != nil {
    return err
  }

  return nil
}

