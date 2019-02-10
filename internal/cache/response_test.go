package cache

import (
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
	type fields struct {
		Path    []byte
		Body    []byte
		Headers []ResponseHeader
	}
	type args struct {
		k []byte
		v []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				Path:    tt.fields.Path,
				Body:    tt.fields.Body,
				Headers: tt.fields.Headers,
			}
			r.SetHeader(tt.args.k, tt.args.v)
		})
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
