package proxy

const proxyReqHeaderKey = "X-Kratgo-Cache"
const proxyReqHeaderValue = "true"

const headerLocation = "Location"
const headerContentEncoding = "Content-Encoding"

const (
	setHeaderAction typeHeaderAction = iota
	unsetHeaderAction
)
