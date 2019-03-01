package invalidator

const (
	invTypeHost invType = iota
	invTypePath
	invTypeHeader
	invTypePathHeader
	invTypeInvalid
)
