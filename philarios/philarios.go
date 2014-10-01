package philarios

type TargetPhrase struct {
  Pretarget string
  Target string
  Posttarget string
}

func (t *TargetPhrase) AlternativeTargets() ([]string) {
  // TODO: implement this
}
