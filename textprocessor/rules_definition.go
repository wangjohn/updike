package textprocessor

var TextProcessorRules = []ProcessingRule{
  FilterBy("EndsWith", "ying").Try(Word{0,-4}, 'i', 'e'),
  FilterBy("EndsWith", Consonant{1}, Consonant{1}, "ing").Try(Word{0,-4}),
  FilterBy("EndsWith", "ing").Try(Word{0,-3}),

  FilterBy("EndsWith", "ied").Try(Word{0,-1}),
  FilterBy("EndsWith", Consonant{1}, Consonant{1}, "ed").Try(Word{0,-3}),
  FilterBy("EndsWith", Consonant{1}, "ed").Try(Word{0,-1}),
  FilterBy("EndsWith", "ed").Try(Word{0,-2}),
}
