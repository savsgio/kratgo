package proxy

const headerLocation = "Location"
const headerContentEncoding = "Content-Encoding"

const (
	setHeaderAction typeHeaderAction = iota
	unsetHeaderAction
)
