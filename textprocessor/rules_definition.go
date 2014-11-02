package textprocessor

var TextProcessorRules = []ProcessingRule{
  FilterBy("EndsWith", "ying").Try(Word{0, -4}, 'i', 'e'),
}
