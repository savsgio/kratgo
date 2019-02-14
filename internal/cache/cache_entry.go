package cache

import (
	"bytes"
	"fmt"
	"sync"
)

var entryPool = sync.Pool{
	New: func() interface{} {
		return new(Entry)
	},
}

// AcquireEntry ...
func AcquireEntry() *Entry {
	return entryPool.Get().(*Entry)
}

// ReleaseEntry ...
func ReleaseEntry(e *Entry) {
	e.Reset()
	entryPool.Put(e)
}

// Reset ...
func (e *Entry) Reset() {
	e.Responses = e.Responses[:0]
}

func (e *Entry) swap(data []Response, i, j int) []Response {
	data[i], data[j] = data[j], data[i]

	return data
}

func (e *Entry) allocResponse(data []Response) ([]Response, *Response) {
	n := len(data)

	if cap(data) > n {
		data = data[:n+1]
	} else {
		data = append(data, Response{})
	}

	return data, &data[n]
}

func (e *Entry) appendResponse(data []Response, resp Response) []Response {
	data, r := e.allocResponse(data)

	r.Path = append(r.Path[:0], resp.Path...)
	r.Body = append(r.Body[:0], resp.Body...)
	r.Headers = resp.Headers

	return data
}

// HasResponse ...
func (e Entry) HasResponse(path []byte) bool {
	for i, n := 0, len(e.Responses); i < n; i++ {
		resp := &e.Responses[i]
		if bytes.Equal(path, resp.Path) {
			return true
		}
	}

	return false
}

// GetAllResponses ...
func (e Entry) GetAllResponses() []Response {
	return e.Responses
}

// Len ...
func (e Entry) Len() int {
	return len(e.Responses)
}

// GetResponse ...
func (e Entry) GetResponse(path []byte) *Response {
	n := len(e.Responses)
	for i := 0; i < n; i++ {
		resp := &e.Responses[i]
		if bytes.Equal(path, resp.Path) {
			return resp
		}
	}

	return nil
}

// SetResponse ...
func (e *Entry) SetResponse(resp Response) {
	r := e.GetResponse(resp.Path)
	if r != nil {
		r.Body = append(r.Body[:0], resp.Body...)
		r.Headers = resp.Headers

		return
	}

	e.Responses = e.appendResponse(e.Responses, resp)
}

// DelResponse ...
func (e *Entry) DelResponse(path []byte) {
	responses := e.GetAllResponses()

	for i, n := 0, len(responses); i < n; i++ {
		resp := &responses[i]
		if bytes.Equal(path, resp.Path) {
			n--
			if i != n {
				e.swap(responses, i, n)
				i--
			}
			responses = responses[:n] // Remove last position
		}
	}

	e.Responses = responses
}

// Marshal ...
func Marshal(src Entry) ([]byte, error) {
	b, err := src.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal: %v", err)
	}

	return b, nil
}

// Unmarshal ...
func Unmarshal(dst *Entry, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if _, err := dst.UnmarshalMsg(data); err != nil {
		return fmt.Errorf("Could not unmarshal '%s': %v", data, err)
	}

	return nil
}
