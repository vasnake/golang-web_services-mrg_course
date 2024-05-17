package csrf_token

var _ Storage = &StDb{} // type check

// empty struct alias
type Unit struct{}

var unit = Unit{}
