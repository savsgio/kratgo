package cache

//go:generate msgp

// ResponseHeaders ...
type ResponseHeaders struct {
	Key   []byte
	Value []byte
}

// Response ...
type Response struct {
	Path    []byte
	Body    []byte
	Headers []ResponseHeaders
}

//Entry ...
type Entry struct {
	Responses []Response
}
