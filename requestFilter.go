package gospider

import "sync"

//FilterRequest 去重复的url
type RequestFilter interface {
	Filter(requests ...Request) []Request
}

//FilterRequest的默认鸟实现
type StoreRequestFilter struct {
	mapRequest map[string]Request
	mux        sync.Mutex
}

func (filter *StoreRequestFilter) Filter(requests ...Request) []Request {
	filter.mux.Lock()
	defer filter.mux.Unlock()
	if len(requests) == 0 {
		return requests
	}
	newRequests := make([]Request, 0, len(requests))
	for _, req := range requests {
		_, ok := filter.mapRequest[req.Url]
		if !ok {
			newRequests = append(newRequests, req)
			filter.mapRequest[req.Url] = req
		}
	}
	return newRequests
}
