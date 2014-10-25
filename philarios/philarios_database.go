package philarios

import (
  "github.com/lib/pq"
  "database/sql"
  "strings"
  "time"
)

type PhilariosDatabase interface {
  QueryForWord(word string, categories []string) ([]Paragraph, error)
  AddPublication(publication Publication) (error)
  EnsureSchema(db *sql.DB) (error)
}

type PhilariosPostgresDatabase struct {
  DriverName string
  DataSourceName string
}

var philariosSchema = `
CREATE TABLE IF NOT EXISTS publications (
  id integer PRIMARY KEY,
  title text,
  author text,
  editor text,
  date timestamp,
  source_url text,
  encoding text,
  type text
);

CREATE TABLE IF NOT EXISTS categories (
  id integer PRIMARY KEY,
  publication integer REFERENCES publications (id),
  category text
);

CREATE TABLE IF NOT EXISTS paragraphs (
  id integer PRIMARY KEY,
  publication integer REFERENCES publications (id),
  body text
);`

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
  Categories []string
}

type Paragraph struct {
  PublicationName string
  Body string
}

/*
QueryForWord returns SQL rows of paragraphs containing the query word given as
an argument. These are returned from the database.
*/
func (p PhilariosPostgresDatabase) QueryForWord(word string) ([]Paragraph, error) {
  db, err := sql.Open(p.DriverName, p.DataSourceName)
  if err != nil {
    return nil, err
  }

  rows, err := performWordQuery(word, db)
  defer rows.Close()
  if err != nil {
    return nil, err
  }

  paragraphs := make([]Paragraph, 0)
  var publicationName, body string
  for rows.Next() {
    err = rows.Scan(&publicationName, &body)
    if err != nil {
      return nil, err
    }
    paragraphs = append(paragraphs, Paragraph{publicationName, body})
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return paragraphs, nil
}

func performWordQuery(word string, db *sql.DB) (*sql.Rows, error) {
  return db.Query(`SELECT body FROM paragraphs
    WHERE to_tsvector(body) @@ to_tsquery(?)`, word)
}

/*
AddPublication adds a new publication to the database, adding data to the
publications, categories, and paragraphs tables.
*/
func (p PhilariosPostgresDatabase) AddPublication(publication Publication) (error) {
  db, err := sql.Open(p.DriverName, p.DataSourceName)
  if err != nil {
    return err
  }
  p.EnsureSchema(db)

  txn, err := db.Begin()
  if err != nil {
    return err
  }

  var publicationId int
  err = db.QueryRow(
    `INSERT INTO publications(
        title, author, editor, date, source_url, type, encoding)
      VALUES (
        ?, ?, ?, ?, ?, ?, ?)
      RETURNING id`,
    publication.Title,
    publication.Author,
    publication.Editor,
    publication.Date.Unix(),
    publication.SourceURL,
    publication.Type,
    publication.Encoding).Scan(&publicationId)
  if err != nil {
    return err
  }

  categoryStmt, err := txn.Prepare(pq.CopyIn("categories", "publication", "category"))
  for _, category := range publication.Categories {
    _, err = categoryStmt.Exec(publicationId, category)
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

/*
EnsureSchema is called on a database and ensures that a correct schema has been
applied so that Queries on the database can occur.
*/
func (p PhilariosPostgresDatabase) EnsureSchema(db *sql.DB) (error) {
  _, err := db.Exec(philariosSchema)
  return err
}
