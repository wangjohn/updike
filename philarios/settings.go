package philarios

type Settings struct {
  WordsToCapture int
}

const (
  WordsToCapture = 4
)

func DefaultSettingsObject() (Settings) {
  return Settings{
    WordsToCapture,
  }
}
