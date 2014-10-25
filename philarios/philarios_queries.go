package philarios

import (
  "database/sql"
)

const (
  DatabaseDriverName = "postgres"
  DatabaseDataSourceName = "user=philarios dbnamephilarios sslmode=verify-full"
)

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
