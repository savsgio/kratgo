package invalidator

import (
	"testing"
)

func getEntryTest() Entry {
	return Entry{
		Host: "www.kratgo.com",
		Path: "/fast/",
		Header: EntryHeader{
			Key:   "X-Data",
			Value: "1",
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
	e.Host = "www.kratgo.es"
	e.Path = "/fast"
	e.Header.Key = "X-Data"
	e.Header.Value = "1"

	ReleaseEntry(e)

	if e.Host != "" {
		t.Errorf("ReleaseEntry() entry has not been reset")
	}
	if e.Path != "" {
		t.Errorf("ReleaseEntry() entry has not been reset")
	}
	if e.Header.Key != "" {
		t.Errorf("ReleaseEntry() entry has not been reset")
	}
	if e.Header.Value != "" {
		t.Errorf("ReleaseEntry() entry has not been reset")
	}
}

func TestHeader_Reset(t *testing.T) {
	e := getEntryTest()

	e.Reset()

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
