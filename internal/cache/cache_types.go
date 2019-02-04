package cache

//go:generate msgp

// ResponseHeaders ...
type ResponseHeaders map[string][]byte

// Response ...
type Response map[string]struct {
	Body    []byte
	Headers ResponseHeaders
}

//Entry ...
type Entry struct {
	Response Response
}
