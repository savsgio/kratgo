package cache

import (
	"bytes"
	"reflect"
	"testing"
)

func getEntryTest() Entry {
	r1 := Response{
		Path: []byte("/cache/"),
		Body: []byte("Response body"),
		Headers: []ResponseHeader{
			{
				Key:   []byte("Key1"),
				Value: []byte("Value1"),
			},
		},
	}
	r2 := Response{
		Path: []byte("/cache/2/"),
		Body: []byte("Response body 2"),
		Headers: []ResponseHeader{
			{
				Key:   []byte("Key1"),
				Value: []byte("Value1"),
			},
		},
	}

	return Entry{
		Responses: []Response{
			r1,
			r2,
		},
	}
}

func TestEntry_Reset(t *testing.T) {
	e := getEntryTest()
	e.Reset()

	if len(e.Responses) > 0 {
		t.Errorf("Entry.Reset() has not been reset")
	}
}

func TestEntry_swap(t *testing.T) {
	e := getEntryTest()
	r1 := e.Responses[0]
	r2 := e.Responses[1]

	e.swap(e.Responses, 0, 1)

	if reflect.DeepEqual(e.Responses[0], r1) {
		t.Errorf("Entry.swap() == '%v', want '%v'", e.Responses[0], r2)
	}

	if reflect.DeepEqual(e.Responses[1], r2) {
		t.Errorf("Entry.swap() == '%v', want '%v'", e.Responses[1], r1)
	}
}

func TestEntry_allocResponse(t *testing.T) {
	e := getEntryTest()

	length := len(e.Responses)

	var r *Response
	e.Responses, r = e.allocResponse(e.Responses)

	if r == nil {
		t.Errorf("Entry.allocResponse() not returns a new response pointer")
	}

	if len(e.Responses) != length+1 {
		t.Errorf("Entry.allocResponse() responses.len == '%d', want '%d'", len(e.Responses), length+1)
	}
}

func TestEntry_appendResponse(t *testing.T) {
	e := getEntryTest()

	length := len(e.Responses)

	r := AcquireResponse()

	e.Responses = e.appendResponse(e.Responses, *r)

	if len(e.Responses) != length+1 {
		t.Errorf("Entry.appendResponse() responses.len == '%d', want '%d'", len(e.Responses), length+1)
	}
}

func TestEntry_HasResponse(t *testing.T) {
	e := getEntryTest()
	r1 := e.Responses[0]

	fakePath := []byte("/fake")

	if !e.HasResponse(r1.Path) {
		t.Errorf("Entry.HasResponse() path '%s' == '%v', want '%v'", r1.Path, false, true)
	}

	if e.HasResponse(fakePath) {
		t.Errorf("Entry.HasResponse() path '%s' == '%v', want '%v'", fakePath, true, false)
	}
}

func TestEntry_GetAllResponses(t *testing.T) {
	e := getEntryTest()

	all := e.GetAllResponses()

	if !reflect.DeepEqual(all, e.Responses) {
		t.Errorf("Entry.GetAllResponses() == '%v', want '%v'", all, e.Responses)
	}
}

func TestEntry_GetResponse(t *testing.T) {
	e := getEntryTest()
	r1 := e.Responses[0]

	fakePath := []byte("/fake")

	if r := e.GetResponse(r1.Path); !reflect.DeepEqual(*r, r1) {
		t.Errorf("Entry.GetResponse() path '%s' == '%v', want '%v'", r1.Path, *r, r1)
	}

	if r := e.GetResponse(fakePath); r != nil {
		t.Errorf("Entry.GetResponse() path '%s' == '%v', want '%v'", fakePath, *r, nil)
	}
}

func TestEntry_SetResponse(t *testing.T) {
	e := getEntryTest()

	r := AcquireResponse()
	r.Path = []byte("/kratgo/fast")
	r.Body = []byte("Body Kratgo Fast")

	e.SetResponse(*r)

	if !e.HasResponse(r.Path) {
		t.Errorf("Entry.SetResponse() has not been set a new response")
	}
}

func TestEntry_DelResponse(t *testing.T) {
	e := getEntryTest()
	r1 := e.Responses[0]

	e.DelResponse(r1.Path)

	if e.HasResponse(r1.Path) {
		t.Errorf("Entry.DelResponse() has not been delete the response")
	}
}

func TestMarshal(t *testing.T) {
	e := getEntryTest()

	expected, err := e.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	marshalData, err := Marshal(&e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !bytes.Equal(marshalData, expected) {
		t.Errorf("Marshal() == '%s', want '%s'", marshalData, expected)
	}
}

func TestUnmarshal(t *testing.T) {
	e := getEntryTest()
	entry := AcquireEntry()

	marshalData, err := Marshal(&e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = Unmarshal(entry, marshalData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(e, *entry) {
		t.Errorf("Marshal() == '%v', want '%v'", e, *entry)
	}
}
