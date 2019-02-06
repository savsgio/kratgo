package cache

import "sync"

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

func (r *Response) allocHeader(data []ResponseHeaders) ([]ResponseHeaders, *ResponseHeaders) {
	n := len(data)

	if cap(data) > n {
		data = data[:n+1]
	} else {
		data = append(data, ResponseHeaders{})
	}

	return data, &data[n]
}

func (r *Response) appendHeader(data []ResponseHeaders, k, v []byte) []ResponseHeaders {
	data, h := r.allocHeader(data)

	h.Key = append(h.Key[:0], k...)
	h.Value = append(h.Value[:0], v...)

	return data
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
