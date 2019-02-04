package proxy

import "sync"

var evalParamsPool = sync.Pool{
	New: func() interface{} {
		return &evalParams{
			p: make(map[string]interface{}),
		}
	},
}

func acquireEvalCtxParams() *evalParams {
	return evalParamsPool.Get().(*evalParams)
}

func releaseEvalCtxParams(ecp *evalParams) {
	ecp.reset()
	evalParamsPool.Put(ecp)
}

func (ecp *evalParams) set(k string, v interface{}) {
	ecp.p[k] = v
}

func (ecp *evalParams) get(k string) (interface{}, bool) {
	v, ok := ecp.p[k]
	return v, ok
}

func (ecp *evalParams) del(k string) {
	delete(ecp.p, k)
}

func (ecp *evalParams) reset() {
	for k := range ecp.p {
		ecp.del(k)
	}
}
