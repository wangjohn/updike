package philarios

import (
  "github.com/lib/pq"
  "database/sql"
  "strings"
  "time"
  "fmt"
)

const (
  DatabaseDriverName = "postgres"
  DatabaseDataSourceName = "user=philarios dbnamephilarios sslmode=verify-full"
)

/*
Publication is a structure which represents any type of publication (such as
books or articles) which contains text.
*/
type Publication struct {
  Title string
  Author string
  Editor string
  Date time.Time
  SourceURL string
  Encoding string
  Type string
  Text string
  Categories []Category
}

/*
Category is a stuct which represents a category type for a publication.
*/
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

/*
AddPublication adds a new publication to the database, adding data to the
publications, categories, and paragraphs tables.
*/
func AddPublication(publication Publication) (error) {
  db, err := sql.Open(DatabaseDriverName, DatabaseDataSourceName)
  if err != nil {
    return err
  }

  txn, err := db.Begin()
  if err != nil {
    return err
  }

  var publicationId int
  publicationInsertQuery := fmt.Sprintf(
    `INSERT INTO publications(
        title, author, editor, date, source_url, type, encoding)
      VALUES (
        %s, %s, %s, %d, %s, %s, %s)
      RETURNING id`,
    publication.Title,
    publication.Author,
    publication.Editor,
    publication.Date.Unix(),
    publication.SourceURL,
    publication.Type,
    publication.Encoding)

  err = db.QueryRow(publicationInsertQuery).Scan(&publicationId)
  if err != nil {
    return err
  }

  categoryStmt, err := txn.Prepare(pq.CopyIn("categories", "publication", "category"))
  for _, category := range publication.Categories {
    _, err = categoryStmt.Exec(publicationId, category.Name)
    if err != nil {
      return err
    }
  }
  err = categoryStmt.Close()
  if err != nil {
    return err
  }

  paragraphStmt, err := txn.Prepare(pq.CopyIn("paragraphs", "publication", "body"))
  paragraphs := strings.Split(publication.Text, "\n")
  for _, paragraph := range paragraphs {
    _, err = paragraphStmt.Exec(publicationId, paragraph)
    if err != nil {
      return err
    }
  }
  err = paragraphStmt.Close()
  if err != nil {
    return err
  }

  err = txn.Commit()
  if err != nil {
    return err
  }

  return nil
}
