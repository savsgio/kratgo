package invalidator

import (
	"testing"
)

func getEntryTest() Entry {
	return Entry{
		Action: "delete",
		Host:   "www.kratgo.com",
		Path:   "/fast/",
		Header: Header{
			Key:   "X-Data",
			Value: "1",
		},
	}
}

func TestHeader_Reset(t *testing.T) {
	e := getEntryTest()

	e.Reset()

	if e.Action != "" {
		t.Errorf("Entry.Action has not been reset")
	}

	if e.Host != "" {
		t.Errorf("Entry.Host has not been reset")
	}

	if e.Path != "" {
		t.Errorf("Entry.Path has not been reset")
	}

	if e.Header.Key != "" {
		t.Errorf("Entry.Header.Key has not been reset")
	}

	if e.Header.Value != "" {
		t.Errorf("Entry.Header.Value has not been reset")
	}
}

func TestEntry_Reset(t *testing.T) {
	e := getEntryTest()

	e.Header.Reset()

	if e.Header.Key != "" {
		t.Errorf("Entry.Header.Key has not been reset")
	}

	if e.Header.Value != "" {
		t.Errorf("Entry.Header.Value has not been reset")
	}
}
