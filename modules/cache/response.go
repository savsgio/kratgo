package cache

import (
	"bytes"
	"sync"
)

var responsePool = sync.Pool{
	New: func() interface{} {
		return new(Response)
	},
}

// AcquireResponse ...
func AcquireResponse() *Response {
	return responsePool.Get().(*Response)
}

// ReleaseResponse ...
func ReleaseResponse(r *Response) {
	r.Reset()
	responsePool.Put(r)
}

func (r *Response) allocHeader(data []ResponseHeader) ([]ResponseHeader, *ResponseHeader) {
	n := len(data)

	if cap(data) > n {
		data = data[:n+1]
	} else {
		data = append(data, ResponseHeader{})
	}

	return data, &data[n]
}

func (r *Response) appendHeader(data []ResponseHeader, k, v []byte) []ResponseHeader {
	data, h := r.allocHeader(data)

	h.Key = append(h.Key[:0], k...)
	h.Value = append(h.Value[:0], v...)

	return data
}

// HasHeader ...
func (r *Response) HasHeader(k, v []byte) bool {
	for i, n := 0, len(r.Headers); i < n; i++ {
		h := &r.Headers[i]
		if bytes.Equal(h.Key, k) && bytes.Equal(h.Value, v) {
			return true
		}
	}

	return false
}

// SetHeader ...
func (r *Response) SetHeader(k, v []byte) {
	r.Headers = r.appendHeader(r.Headers, k, v)
}

// Reset reset response
func (r *Response) Reset() {
	r.Path = r.Path[:0]
	r.Body = r.Body[:0]
	r.Headers = r.Headers[:0]
}
