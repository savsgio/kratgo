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

func TestAcquireResponse(t *testing.T) {
	r := AcquireResponse()
	if r == nil {
		t.Errorf("AcquireResponse() returns '%v'", nil)
	}
}

func TestReleaseResponse(t *testing.T) {
	r := AcquireResponse()
	r.Path = []byte("/kratgo")
	r.Body = []byte("Kratgo is ultra fast")
	r.SetHeader([]byte("key"), []byte("value"))

	ReleaseResponse(r)

	if len(r.Path) > 0 || len(r.Body) > 0 || len(r.Headers) > 0 {
		t.Errorf("ReleaseResponse() response has not been reset")
	}
}

func TestResponse_allocHeader(t *testing.T) {
	r := AcquireResponse()
	r.SetHeader([]byte("key1"), []byte("value1"))

	wantCapacity := cap(r.Headers) * 2 // Capacity is incremented by power of two

	var h *ResponseHeader
	r.Headers, h = r.allocHeader(r.Headers)

	if h == nil {
		t.Errorf("Response.allocHeader() not returns a new header pointer")
	}

	if cap(r.Headers) != wantCapacity {
		t.Errorf("Response.allocHeader() headers capacity == '%d', want '%d'", cap(r.Headers), wantCapacity)
	}

	r.Headers = r.Headers[:len(r.Headers)-1]
	r.Headers, h = r.allocHeader(r.Headers)

	if cap(r.Headers) != wantCapacity {
		t.Errorf("Response.allocHeader() headers capacity == '%d', want '%d'", cap(r.Headers), wantCapacity)
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

func TestResponse_HasHeader(t *testing.T) {
	r := getResponseTest()

	k := []byte("newKey")
	v := []byte("newValue")

	r.SetHeader(k, v)

	if !r.HasHeader(k, v) {
		t.Errorf("The header '%s = %s' not found", k, v)
	}

	if r.HasHeader(k, []byte("other value")) {
		t.Errorf("The header '%s = %s' found", k, v)
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
