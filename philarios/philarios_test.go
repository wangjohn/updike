package philarios

import (
  "database/sql"
  "testing"

  "github.com/wangjohn/updike/tfidf"
)

const (
  tfidfDriverName = "postgres"
  tfidfDataSourceName = "host=localhost user=philarios dbname=philarios_tfidf_test sslmode=disable"

  storageDriverName = "postgres"
  storageDataSourceName = "host=localhost user=philarios dbname=philarios_storage_test sslmode=disable"
)

func setupWordFactory() (*WordFactory, error) {
  storageDb, err := sql.Open(storageDriverName, storageDataSourceName)
  if err != nil {
    return nil, err
  }

  storage := PostgresStorage{storageDb}
  settings := DefaultSettingsObject()
  tfidfDb, err := sql.Open(tfidfDriverName, tfidfDataSourceName)

  if err != nil {
    return nil, err
  }
  tfidf := tfidf.PersistentTFIDF{tfidfDb}
  storage.AddPublication(Publication{
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
  })

  wordFactory := WordFactory{storage, settings, tfidf}
  return &wordFactory, nil
}

func TestAlternativeWords(t *testing.T) {
  fixtures := []struct {
    Word string
    MaxWords int
  }{
    {"Name", 3},
  }

  wordFactory, err := setupWordFactory()
  if err != nil {
    t.Errorf("Error setting up word factory: %v", err)
  }

  for _, fixture := range fixtures {
    alternatives, err := wordFactory.AlternativeWords(fixture.Word, fixture.MaxWords)
    if err != nil {
      t.Errorf("Error obtaining alternative words: %v", err)
    }

    t.Errorf("Alternative words: %v", alternatives)
  }
}
