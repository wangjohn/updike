package tfidf

import (
  "github.com/reiver/go-porterstemmer"

  _ "github.com/lib/pq"
  "database/sql"
  "math"
  "fmt"
)

type TFIDF interface {
  Store(word string, occurrences, docMaxWordOccurrences, documentId int) (error)
  TermFrequency(word string, documentId int) (float64, error)
  InverseDocumentFrequency(word string) (float64, error)
  Score(word string, documentId int) (float64, error)
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
  document bigserial
);

CREATE TABLE IF NOT EXISTS document_frequency (
  id bigserial PRIMARY KEY,
  word text,
  unique_documents integer
);
`

func (p PersistentTFIDF) EnsureSchema() (error) {
  _, err := p.SQLDatabase.Exec(persistentSqlSchema)
  return err
}

func (p PersistentTFIDF) TermFrequency(word string, documentId int) (float64, error) {
  word, err := p.NormalizeWord(word)
  if err != nil {
    return 0.0, err
  }

  var freq, docMaxWordFreq int
  err = p.SQLDatabase.QueryRow(
    `SELECT freq, doc_max_word_freq FROM word_document_pairs
     WHERE word=$1
     AND document=$2`, word, documentId).Scan(&freq, &docMaxWordFreq)

  if err == sql.ErrNoRows {
    // We don't have that word, document pair, so just set freq to zero.
    freq = 0

    findMaxFreqErr := p.SQLDatabase.QueryRow(
      `SELECT doc_max_word_freq FROM word_document_pairs
       WHERE document=$1
       LIMIT 1`, documentId).Scan(&docMaxWordFreq)

    if findMaxFreqErr == sql.ErrNoRows {
      return 0.0, fmt.Errorf("Document with id=%v does not exist", documentId)
    } else if findMaxFreqErr != nil {
      return 0.0, findMaxFreqErr
    }
  } else if err != nil {
    return 0.0, err
  }

  return tfFunc(freq, docMaxWordFreq), nil
}

func tfFunc(frequency, docMaxWordFrequency int) float64 {
  return 0.5 + (0.5 * float64(frequency)) / float64(docMaxWordFrequency)
}

var totalDocs = -1

func (p PersistentTFIDF) InverseDocumentFrequency(word string) (float64, error) {
  word, err := p.NormalizeWord(word)
  if err != nil {
    return 0.0, err
  }

  var uniqDocs int
  err = p.SQLDatabase.QueryRow(
    `SELECT unique_documents FROM document_frequency
    WHERE word=$1`, word).Scan(&uniqDocs)

  if err == sql.ErrNoRows {
    uniqDocs = 0
  } else if err != nil {
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

func idfFunc(docs, totDocs int) (float64) {
  return math.Log10(float64(totDocs) / (1.0 + float64(docs)))
}

func (p PersistentTFIDF) Score(word string, documentId int) (float64, error) {
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

  var isNewDocument bool
  var id int
  wordQueryErr := p.SQLDatabase.QueryRow(
   `SELECT id FROM word_document_pairs
    WHERE word=$1
    AND document=$2`, word, documentId).Scan(&id)

  if wordQueryErr == sql.ErrNoRows {
    isNewDocument = true
    wordInsErr := p.SQLDatabase.QueryRow(
     `INSERT INTO word_document_pairs(
        word, freq, doc_max_word_freq, document)
      VALUES ($1, $2, $3, $4)
      RETURNING id`,
      word,
      occurrences,
      docMaxWordOccurrences,
      documentId).Scan(&id)
    if wordInsErr != nil {
      return wordInsErr
    }
  } else if wordQueryErr == nil {
    isNewDocument = false
    _, wordUpdErr := p.SQLDatabase.Exec(
     `UPDATE word_document_pairs
      SET freq=$1
      AND doc_max_word_freq=$2
      WHERE word=$3
      AND document=$4`,
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
  var docFreqId int
  docFreqQueryErr := p.SQLDatabase.QueryRow(
    `SELECT id FROM document_frequency
     WHERE word=$1`, word).Scan(&docFreqId)

  if docFreqQueryErr == sql.ErrNoRows {
    docFreqInsErr := p.SQLDatabase.QueryRow(
     `INSERT INTO document_frequency(
        word, unique_documents)
      VALUES ($1, $2)
      RETURNING id`, word, 0).Scan(&docFreqId)

    if docFreqInsErr != nil {
      return docFreqInsErr
    }
  } else if docFreqQueryErr != nil {
    return docFreqQueryErr
  }

  if isNewDocument {
    _, err = p.SQLDatabase.Exec(
      `UPDATE document_frequency
       SET unique_documents = unique_documents + 1
       WHERE id=$1`, docFreqId)
    return err
  } else {
    return nil
  }
}

func (p PersistentTFIDF) NormalizeWord(word string) (string, error) {
  return porterstemmer.StemString(word), nil
}
