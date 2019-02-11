package cache

import (
	"bytes"
	"testing"
)

func getResponseTest() Response {
	return Response{
		Path: []byte("/cache/"),
		Body: []byte("Response body"),
		Headers: []ResponseHeader{
			{
				Key:   []byte("Key1"),
				Value: []byte("Value1"),
			},
		},
	}
}
func TestResponse_allocHeader(t *testing.T) {
	r := getResponseTest()

	length := len(r.Headers)

	var h *ResponseHeader
	r.Headers, h = r.allocHeader(r.Headers)

	if h == nil {
		t.Errorf("Entry.allocHeader() not returns a new header pointer")
	}

	if len(r.Headers) != length+1 {
		t.Errorf("Entry.allocResponse() responses.len == '%d', want '%d'", len(r.Headers), length+1)
	}
}

func TestResponse_appendHeader(t *testing.T) {
	r := getResponseTest()

	length := len(r.Headers)

	r.Headers = r.appendHeader(r.Headers, []byte("key"), []byte("value"))

	if len(r.Headers) != length+1 {
		t.Errorf("Entry.appendResponse() responses.len == '%d', want '%d'", len(r.Headers), length+1)
	}
}

func TestResponse_SetHeader(t *testing.T) {
	r := getResponseTest()

	k := []byte("newKey")
	v := []byte("newValue")

	r.SetHeader(k, v)

	finded := false
	for _, h := range r.Headers {
		if bytes.Equal(h.Key, k) && bytes.Equal(h.Value, v) {
			finded = true
			break
		}
	}

	if !finded {
		t.Errorf("The header '%s==%s' has not been set", k, v)
	}
}

func TestResponse_Reset(t *testing.T) {
	r := getResponseTest()

	r.Reset()

	if len(r.Path) > 0 {
		t.Errorf("Response.Path has not been reset")
	}

	if len(r.Body) > 0 {
		t.Errorf("Response.Body has not been reset")
	}

	if len(r.Headers) > 0 {
		t.Errorf("Response.Headers has not been reset")
	}
}
