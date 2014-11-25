package philarios

type Settings struct {
  WordsToCapture int
}

const (
  WordsToCapture = 2
)

func DefaultSettingsObject() (Settings) {
  return Settings{
    WordsToCapture,
  }
}
