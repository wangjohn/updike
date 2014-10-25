package philarios

import (
  "github.com/lib/pq"
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

func AddBookData(book Book) (error) {
  db, err := sql.Open(DatabaseDriverName, DatabaseDataSourceName)
  if err != nil {
    return err
  }

  txn, err := db.Begin()
  if err != nil {
    return err
  }

  bookStmt, err := txn.Prepare(pq.CopyIn("books", "title", "author", "date"))
  bookResult, err = bookStmt.Exec(book.Title, book.Author, book.Date)
  if err != nil {
    return err
  }
  bookId := bookResult.LastInsertId()

  paragraphStmt, err := txn.Prepare(pq.CopyIn("paragraphs", "book", "body"))
  paragraphs := strings.Split(book.Text, "\n")
  for _, paragraph := range paragraphs {
    _, err = paragraphStmt.Exec(bookId, paragraph)
    if err != nil {
      return err
    }
  }

  categoryStmt, err := txn.Prepare(pq.CopyIn("categories", "book", "category"))
  for _, category := range categories {
    _, err = categoryStmt.Exec(bookId, category)
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
