package philarios

import (
  "github.com/lib/pq"
  "database/sql"
  "time"
)

const (
  DatabaseDriverName = "postgres"
  DatabaseDataSourceName = "user=philarios dbnamephilarios sslmode=verify-full"
)

type Publication struct {
  Title string
  Author string
  Editor string
  Date time.Time
  SourceURL string
  Encoding string
  Text string
  Categories []Category
}

type Category struct {
  Name string
}

/*
QueryForWord returns SQL rows of paragraphs containing the query word given as
an argument. These are returned from the database.
*/
func QueryForWord(word string) (*sql.Rows, error) {
  db, err := sql.Open(DatabaseDriverName, DatabaseDataSourceName)
  if err != nil {
    return nil, err
  }

  rows, err := performWordQuery(word, db)

  if err != nil {
    return nil, err
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return rows, nil
}

func performWordQuery(word string, db *sql.DB) (*sql.Rows, error) {
  return db.Query(`SELECT body FROM paragraphs
    WHERE to_tsvector(body) @@ to_tsquery(?)`, word)
}

func AddPublication(publication Publication) (error) {
  db, err := sql.Open(DatabaseDriverName, DatabaseDataSourceName)
  if err != nil {
    return err
  }

  txn, err := db.Begin()
  if err != nil {
    return err
  }

  publicationStmt, err := txn.Prepare(pq.CopyIn("publications", "title", "author", "date"))
  publicationResult, err = publicationStmt.Exec(publication.Title, publication.Author, publication.Date)
  if err != nil {
    return err
  }
  publicationId := publicationResult.LastInsertId()

  paragraphStmt, err := txn.Prepare(pq.CopyIn("paragraphs", "publication", "body"))
  paragraphs := strings.Split(publication.Text, "\n")
  for _, paragraph := range paragraphs {
    _, err = paragraphStmt.Exec(publicationId, paragraph)
    if err != nil {
      return err
    }
  }

  categoryStmt, err := txn.Prepare(pq.CopyIn("categories", "publication", "category"))
  for _, category := range publication.Categories {
    _, err = categoryStmt.Exec(publicationId, category.Name)
    if err != nil {
      return err
    }
  }

  _, err = stmt.Exec()
  if err != nil {
    return err
  }

  _, err = stmt.Close()
  if err != nil {
    return err
  }

  err = txn.Commit()
  if err != nil {
    return err
  }

  return nil
}
