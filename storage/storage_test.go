package philarios

import (
  "database/sql"
  "testing"
)

const (
  testDriverName = "postgres"
  testDataSourceName = "host=localhost user=philarios dbname=philarios_storage_test sslmode=disable"
)

func teardownDatabase(db *sql.DB) {
  db.Exec(`
    DROP TABLE IF EXISTS paragraphs;
    DROP TABLE IF EXISTS categories;
    DROP TABLE IF EXISTS publications;
    DROP TABLE IF EXISTS frequencies;
  `)
}

func setupDatabase() (Storage, error) {
  db, err := sql.Open(testDriverName, testDataSourceName)
  if err != nil {
    return nil, err
  }

  philariosDatabase := PostgresStorage{db}
  teardownDatabase(db)

  publication := Publication{
    Title: "Great Expectations",
    Author: "Charles Dickens",
    Editor: "",
    Date: "1860-01-12",
    SourceURL: "",
    Encoding: "utf-8",
    Type: "book",
    Categories: []string{"classic", "dickens"},
    Text:
`My father's family name being Pirrip, and my Christian name Philip, my infant tongue could make of both names nothing longer or more explicit than Pip. So, I called myself Pip, and came to be called Pip.
I give Pirrip as my father's family name, on the authority of his tombstone and my sister,—Mrs. Joe Gargery, who married the blacksmith. As I never saw my father or my mother, and never saw any likeness of either of them (for their days were long before the days of photographs), my first fancies regarding what they were like were unreasonably derived from their tombstones. The shape of the letters on my father's, gave me an odd idea that he was a square, stout, dark man, with curly black hair. From the character and turn of the inscription, "Also Georgiana Wife of the Above," I drew a childish conclusion that my mother was freckled and sickly. To five little stone lozenges, each about a foot and a half long, which were arranged in a neat row beside their grave, and were sacred to the memory of five little brothers of mine,—who gave up trying to get a living, exceedingly early in that universal struggle,—I am indebted for a belief I religiously entertained that they had all been born on their backs with their hands in their trousers-pockets, and had never taken them out in this state of existence.
Ours was the marsh country, down by the river, within, as the river wound, twenty miles of the sea. My first most vivid and broad impression of the identity of things seems to me to have been gained on a memorable raw afternoon towards evening. At such a time I found out for certain that this bleak place overgrown with nettles was the churchyard; and that Philip Pirrip, late of this parish, and also Georgiana wife of the above, were dead and buried; and that Alexander, Bartholomew, Abraham, Tobias, and Roger, infant children of the aforesaid, were also dead and buried; and that the dark flat wilderness beyond the churchyard, intersected with dikes and mounds and gates, with scattered cattle feeding on it, was the marshes; and that the low leaden line beyond was the river; and that the distant savage lair from which the wind was rushing was the sea; and that the small bundle of shivers growing afraid of it all and beginning to cry, was Pip.
"Hold your noise!" cried a terrible voice, as a man started up from among the graves at the side of the church porch. "Keep still, you little devil, or I'll cut your throat!"
A fearful man, all in coarse gray, with a great iron on his leg. A man with no hat, and with broken shoes, and with an old rag tied round his head. A man who had been soaked in water, and smothered in mud, and lamed by stones, and cut by flints, and stung by nettles, and torn by briars; who limped, and shivered, and glared, and growled; and whose teeth chattered in his head as he seized me by the chin.
"Oh! Don't cut my throat, sir," I pleaded in terror. "Pray don't do it, sir."
"Tell us your name!" said the man. "Quick!"
"Pip, sir."
"Once more," said the man, staring at me. "Give it mouth!"
"Pip. Pip, sir."
`,
  }

  err = philariosDatabase.AddPublication(publication)
  if err != nil {
    return philariosDatabase, err
  }

  return philariosDatabase, nil
}

func TestCreatingAndQueryingPostgresDatabaseWithoutCategories(t *testing.T) {
  philariosDatabase, err := setupDatabase()
  if err != nil {
    t.Errorf("Error setting up database and seeding with data: %s", err.Error())
  }

  paragraphs, err := philariosDatabase.QueryForWord("Georgiana", nil)
  if err != nil {
    t.Errorf("Shouldn't have thrown an error when querying for word: %s", err.Error())
  }

  expectedParagraphs := 2
  if len(paragraphs) != expectedParagraphs {
    t.Errorf("Should have obtained %d paragraphs, instead obtained %d", expectedParagraphs, len(paragraphs))
  }

  expectedPublicationId := 1
  for _, paragraph := range paragraphs {
    if paragraph.PublicationId != expectedPublicationId {
      t.Errorf("Should have obtained paragraphs with PublicationId %d, instead obtained %d",
        expectedPublicationId, paragraph.PublicationId)
    }
  }
}
