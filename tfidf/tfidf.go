package tfidf

import (
  "database/sql"
  "math"
)

type TFIDF interface {
  Store(word string) (error)
  TermFrequency(word string, documentId int) (float64, error)
  InverseDocumentFrequency(word string) (float64, error)
  TFIDFScore(word string, documentId int) (float64, error)
  NormalizeWord(word string) (string, error)
}

type PersistentTFIDF struct {
  SQLDatabase *sql.DB
}

var persistentSqlSchema = `
CREATE TABLE IF NOT EXISTS word_document_pairs (
  id bigserial PRIMARY KEY,
  word text,
  freq integer,
  doc_max_word_freq integer,
  document bigserial,
);

CREATE TABLE IF NOT EXISTS document_frequency (
  id bigserial PRIMARY KEY,
  word text,
  unique_documents integer,
);
`

func (p PersistentTFIDF) TermFrequency(word string, documentId int) (float64, error) {
  word, err := p.NormalizeWord(word)
  if err != nil {
    return 0.0, err
  }

  var freq, docMaxWordFreq float64
  err = p.SQLDatabase.QueryRow(
    `SELECT freq, doc_max_word_freq FROM word_document_pairs
     WHERE word=?
     AND document=?`, word, documentId).Scan(&freq, &docMaxWordFreq)
  if err != nil {
    return 0.0, err
  }

  return tfFunc(freq, docMaxWordFreq), nil
}

func tfFunc(frequency, docMaxWordFrequency int) float64 {
  return 0.5 + (0.5 * frequency) / docMaxWordFrequency
}

var totalDocs = -1

func (p PersistentTFIDF) InverseDocumentFrequency(word string) (float64, error) {
  word, err := p.NormalizeWord(word)
  if err != nil {
    return 0.0, err
  }

  var uniqDocs float64
  err = p.SQLDatabase.QueryRow(
    `SELECT unique_documents FROM document_frequency
    WHERE word=?`, word).Scan(&uniqDocs)
  if err != nil {
    return 0.0, err
  }

  if totalDocs == -1 {
    err = p.SQLDatabase.QueryRow(
      `SELECT COUNT(DISTINCT document) FROM word_document_pairs`).Scan(&totalDocs)
    if err != nil {
      return 0.0, err
    }
  }

  return idfFunc(uniqDocs, totalDocs), nil
}

func idfFunc(uniqDocs int, docs int) (float64) {
  return math.Log10(docs / (1 + uniqDocs))
}

func (p PersistentTFIDF) TFIDFScore(word string, documentId int) (float64, error) {
  word, err := p.NormalizeWord(word)
  if err != nil {
    return 0.0, err
  }

  tf, err := p.TermFrequency(word, documentId)
  if err != nil {
    return 0.0, err
  }

  idf, err := p.InverseDocumentFrequency(word)
  if err != nil {
    return 0.0, err
  }

  return tf * idf, nil
}

func (p PersistentTFIDF) Store(word string, occurrences, docMaxWordOccurrences, documentId int) (error) {
  word, err := p.NormalizeWord(word)
  if err != nil {
    return err
  }

  var isNewDocument bool, id int
  wordQueryErr := p.SQLDatabase.QueryRow(
   `SELECT id FROM word_document_pairs
    WHERE word='?'
    AND document=?`, word, documentId).Scan(&id)

  if wordQueryErr == sql.ErrNoRows {
    isNewDocument = true
    insertPairQuery := fmt.Sprintf(
     `INSERT INTO word_document_pairs(
        word, freq, doc_max_word_freq, document)
      VALUES ('%s', %d, %d, %d)
      RETURNING id`,
      word,
      occurrences,
      docMaxWordOccurrences,
      documentId)
    wordInsErr = p.SQLDatabase.QueryRow(insertPairQuery).Scan(&id)
    if wordInsErr != nil {
      return wordInsErr
    }
  } else if wordQueryErr == nil {
    isNewDocument = false
    _, wordUpdErr := p.SQLDatabase.Exec(
     `UPDATE word_document_pairs
      SET freq=?
      AND doc_max_word_freq=?
      WHERE word='?'
      AND document=?`,
      occurrences,
      docMaxWordOccurrences,
      word,
      documentId)
    if wordUpdErr != nil {
      return wordUpdErr
    }
  } else {
    return wordQueryErr
  }

  // Update the number of unique documents
  var docFreqId, uniqDocs int
  docFreqQueryErr := p.SQLDatabase.QueryRow(
    `SELECT id, unique_documents FROM document_frequency
     WHERE word='?'`, word).Scan(&docFreqId, &uniqDocs)

  if docFreqQueryErr == sql.ErrNoRows {
    docFreqInsErr := p.SQLDatabase.QueryRow(
     `INSERT INTO document_frequency(
        word, unique_documents)
      VALUES ('%s', %d)
      RETURNING id, unique_documents`, word, 0).Scan(&docFreqId, &uniqDocs)

    if docFreqInsErr != nil {
      return docFreqInsErr
    }
  } else if docFreqQueryErr != nil {
    return docFreqQueryErr
  }

  var updatedUniqDocs int
  if isNewDocument {
    updatedUniqDocs = uniqDocs + 1
  } else {
    updatedUniqDocs = uniqDocs
  }

  _, err = p.SQLDatabase.Exec(
    `UPDATE document_frequency
     SET unique_documents=?
     WHERE id=?`, docFreqId, updatedUniqDocs)

  return err
}

func (p PersistentTFIDF) NormalizeWord(word string) (string, error) {
  return word, nil
}
