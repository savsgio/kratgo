package proxy

import "sync"

var evalParamsPool = sync.Pool{
	New: func() interface{} {
		return &evalParams{
			p: make(map[string]interface{}),
		}
	},
}

func acquireEvalParams() *evalParams {
	return evalParamsPool.Get().(*evalParams)
}

func releaseEvalParams(ep *evalParams) {
	ep.reset()
	evalParamsPool.Put(ep)
}

func (ep *evalParams) set(k string, v interface{}) {
	ep.p[k] = v
}

func (ep *evalParams) get(k string) (interface{}, bool) {
	v, ok := ep.p[k]
	return v, ok
}

func (ep *evalParams) all() map[string]interface{} {
	return ep.p
}

func (ep *evalParams) del(k string) {
	delete(ep.p, k)
}

func (ep *evalParams) reset() {
	for k := range ep.p {
		ep.del(k)
	}
}
