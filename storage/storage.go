package philarios

import (
  "github.com/wangjohn/updike/textprocessor"
  "github.com/lib/pq"
  "database/sql"
  "fmt"
)

type Storage interface {
  QueryForWord(word string, categories []string) ([]Paragraph, error)
  AddPublication(publication Publication) (error)
}

type PostgresStorage struct {
  DriverName string
  DataSourceName string
}

var philariosSchema = `
CREATE TABLE IF NOT EXISTS publications (
  id bigserial PRIMARY KEY,
  title text,
  author text,
  editor text,
  date date,
  source_url text,
  encoding text,
  type text
);

CREATE TABLE IF NOT EXISTS categories (
  id bigserial PRIMARY KEY,
  publication integer REFERENCES publications (id),
  category text
);

CREATE TABLE IF NOT EXISTS paragraphs (
  id bigserial PRIMARY KEY,
  publication integer REFERENCES publications (id),
  body text
);

CREATE TABLE IF NOT EXISTS frequencies (
  id bigserial PRIMARY KEY,
  word text,
  count integer,
);
`

/*
Publication is a structure which represents any type of publication (such as
books or articles) which contains text.
*/
type Publication struct {
  Title string
  Author string
  Editor string
  Date string
  SourceURL string
  Encoding string
  Type string
  Text string
  Categories []string
}

type Paragraph struct {
  PublicationId int
  Body string
}

/*
QueryForWord returns SQL rows of paragraphs containing the query word given as
an argument. These are returned from the database.
*/
func (p PostgresStorage) QueryForWord(word string, categories []string) ([]Paragraph, error) {
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
  var publicationId int
  var body string
  for rows.Next() {
    err = rows.Scan(&publicationId, &body)
    if err != nil {
      return nil, err
    }
    paragraphs = append(paragraphs, Paragraph{publicationId, body})
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return paragraphs, nil
}

func performWordQuery(word string, db *sql.DB) (*sql.Rows, error) {
  wordQuery := fmt.Sprintf(`SELECT publication, body FROM paragraphs
    WHERE to_tsvector(body) @@ to_tsquery('%s')`, word)
  return db.Query(wordQuery)
}

/*
AddPublication adds a new publication to the database, adding data to the
publications, categories, and paragraphs tables.
*/
func (p PostgresStorage) AddPublication(publication Publication) (error) {
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
  publicationQuery := fmt.Sprintf(
    `INSERT INTO publications(
        title, author, editor, date, source_url, type, encoding)
      VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')
      RETURNING id`,
    publication.Title,
    publication.Author,
    publication.Editor,
    publication.Date,
    publication.SourceURL,
    publication.Type,
    publication.Encoding)
  err = db.QueryRow(publicationQuery).Scan(&publicationId)
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
  paragraphs, err := textprocessor.ProcessParagraphs(publication.Text)
  if err != nil {
    return err
  }

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
func (p PostgresStorage) EnsureSchema(db *sql.DB) (error) {
  _, err := db.Exec(philariosSchema)
  return err
}
