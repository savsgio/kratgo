package proxy

const clientReqHeaderKey = "X-Kratgo-Cache"
const clientReqHeaderValue = "true"

const headerLocation = "Location"
const headerContentEncoding = "Content-Encoding"

const (
	setHeaderAction typeHeaderAction = iota
	unsetHeaderAction
)
