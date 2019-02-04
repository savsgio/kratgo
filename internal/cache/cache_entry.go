package cache

import (
	"fmt"
	"sync"
)

var entryPool = sync.Pool{
	New: func() interface{} {
		return &Entry{
			Response: make(Response),
		}
	},
}

// AcquireEntry ...
func AcquireEntry() *Entry {
	return entryPool.Get().(*Entry)
}

// ReleaseEntry ...
func ReleaseEntry(entry *Entry) {
	entry.Reset()
	entryPool.Put(entry)
}

// Reset ...
func (entry *Entry) Reset() {
	for k := range entry.Response {
		delete(entry.Response, k)
	}
}

// Marshal ...
func Marshal(src *Entry) ([]byte, error) {
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
