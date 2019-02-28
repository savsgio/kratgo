package proxy

import (
	"reflect"
	"testing"
)

func TestAcquireEvalParams(t *testing.T) {
	ep := acquireEvalParams()
	if ep == nil {
		t.Errorf("acquireEvalParams() returns '%v'", nil)
	}
}

func TestReleaseEvalParams(t *testing.T) {
	ep := acquireEvalParams()
	ep.set("key", "value")

	releaseEvalParams(ep)

	if len(ep.p) > 0 {
		t.Errorf("releaseEvalParams() entry has not been reset")
	}
}

func Test_evalParams_set(t *testing.T) {
	ep := acquireEvalParams()

	k := "Kratgo"
	v := "fast"

	ep.set(k, v)

	val, ok := ep.p[k]
	if !ok {
		t.Errorf("evalParams.set() the key '%s' not found", k)
	}

	if val.(string) != v {
		t.Errorf("evalParams.set() value == '%s', want '%s'", val, v)
	}
}

func Test_evalParams_get(t *testing.T) {
	ep := acquireEvalParams()

	k := "Kratgo"
	v := "fast"

	ep.set(k, v)

	val, ok := ep.get(k)
	if !ok {
		t.Errorf("evalParams.set() the key '%s' not found", k)
	}

	if val.(string) != v {
		t.Errorf("evalParams.set() value == '%s', want '%s'", val, v)
	}
}

func Test_evalParams_all(t *testing.T) {
	ep := acquireEvalParams()

	data := map[string]interface{}{
		"Kratgo": "fast",
		"slow":   false,
	}
	for k, v := range data {
		ep.set(k, v)
	}

	all := ep.all()

	if !reflect.DeepEqual(data, all) {
		t.Errorf("evalParams.all() == '%v', want '%v'", all, data)
	}
}

func Test_evalParams_del(t *testing.T) {
	ep := acquireEvalParams()

	expected := map[string]interface{}{
		"Kratgo": "fast",
	}
	keyToDelete := "slow"

	data := map[string]interface{}{
		"Kratgo":    "fast",
		keyToDelete: false,
	}
	for k, v := range data {
		ep.set(k, v)
	}

	ep.del(keyToDelete)

	all := ep.all()

	if !reflect.DeepEqual(all, expected) {
		t.Errorf("evalParams.all() == '%v', want '%v'", all, expected)
	}
}

func Test_evalParams_reset(t *testing.T) {
	ep := acquireEvalParams()

	data := map[string]interface{}{
		"Kratgo": "fast",
		"slow":   false,
	}
	for k, v := range data {
		ep.set(k, v)
	}

	ep.reset()

	if len(ep.p) > 0 {
		t.Errorf("evalParams.reset() not reset, current value: %v", ep.p)
	}
}
