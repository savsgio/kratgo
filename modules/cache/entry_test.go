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

func TestAcquireEntry(t *testing.T) {
	e := AcquireEntry()
	if e == nil {
		t.Errorf("AcquireEntry() returns '%v'", nil)
	}
}

func TestReleaseEntry(t *testing.T) {
	e := AcquireEntry()
	r := AcquireResponse()

	e.SetResponse(*r)

	ReleaseEntry(e)

	if e.Len() > 0 {
		t.Errorf("ReleaseEntry() entry has not been reset")
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

	wantCapacity := cap(e.Responses) * 2 // Capacity is incremented by power of two

	var r *Response
	e.Responses, r = e.allocResponse(e.Responses)

	if r == nil {
		t.Errorf("Entry.allocResponse() not returns a new response pointer")
	}

	if cap(e.Responses) != wantCapacity {
		t.Errorf("Entry.allocResponse() responses capacity == '%d', want '%d'", cap(e.Responses), wantCapacity)
	}

	e.Responses = e.Responses[:len(e.Responses)-1]
	e.Responses, _ = e.allocResponse(e.Responses)

	if cap(e.Responses) != wantCapacity {
		t.Errorf("Entry.allocResponse() responses capacity == '%d', want '%d'", cap(e.Responses), wantCapacity)
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

func TestEntry_Len(t *testing.T) {
	e := getEntryTest()

	length := len(e.GetAllResponses())

	if e.Len() != length {
		t.Errorf("Entry.Len() == '%d', want '%d'", e.Len(), length)
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

	wantLength := e.Len() + 1

	r := AcquireResponse()
	r.Path = []byte("/kratgo/fast")
	r.Body = []byte("Body Kratgo Fast")

	e.SetResponse(*r)

	length := e.Len()

	if !e.HasResponse(r.Path) {
		t.Errorf("Entry.SetResponse() has not been set a new response")
	}

	// Update respose (same r.Path)
	r.Body = []byte("UPDATED Body Kratgo Fast")
	r.SetHeader([]byte("key"), []byte("value"))

	e.SetResponse(*r)

	if length != wantLength {
		t.Errorf("Entry.SetResponse() has not been update the existing response")
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

	marshalData, err := Marshal(e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !bytes.Equal(marshalData, expected) {
		t.Errorf("Marshal() == '%s', want '%s'", marshalData, expected)
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		charsToDelete int
	}

	type want struct {
		empty bool
		err   bool
	}

	e := getEntryTest()
	marshalData, _ := Marshal(e)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				charsToDelete: 0,
			},
			want: want{
				empty: false,
				err:   false,
			},
		},
		{
			name: "Empty",
			args: args{
				charsToDelete: len(marshalData),
			},
			want: want{
				empty: true,
				err:   false,
			},
		},
		{
			name: "Error",
			args: args{
				charsToDelete: 1,
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := AcquireEntry()

			err := Unmarshal(entry, marshalData[tt.args.charsToDelete:])
			if (err != nil) != tt.want.err {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.want.err {
				return
			}

			if tt.want.empty {
				if entry.Len() > 0 {
					t.Errorf("Unmarshal() entry has not been empty: %v'", *entry)
				}

			} else if !reflect.DeepEqual(e, *entry) {
				t.Errorf("Unmarshal() == '%v', want '%v'", *entry, e)
			}
		})
	}
}
