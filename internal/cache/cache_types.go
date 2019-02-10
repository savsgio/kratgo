package cache

//go:generate msgp

// ResponseHeader ...
type ResponseHeader struct {
	Key   []byte
	Value []byte
}

// Response ...
type Response struct {
	Path    []byte
	Body    []byte
	Headers []ResponseHeader
}

//Entry ...
type Entry struct {
	Responses []Response
}
