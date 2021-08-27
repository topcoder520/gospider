package gospider

import "sync"

type Scheduler interface {
	Push(requests ...Request)
	Poll() Request
	//取n个 返回数据切片和取到的真是数量
	PollN(n int) ([]Request, int)
	Len() int
}

//Scheduler的默认实现
type RequestScheduler struct {
	requests []Request
	mux      sync.Mutex
}

func (s *RequestScheduler) Push(reqs ...Request) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.requests == nil {
		s.requests = make([]Request, 0, len(reqs))
	}
	s.requests = append(s.requests, reqs...)
}
func (s *RequestScheduler) Poll() Request {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.requests == nil {
		s.requests = make([]Request, 0)
	}
	m := len(s.requests)
	if m == 0 {
		return Request{Skip: true}
	}
	u := s.requests[0]
	if m == 1 {
		s.requests = s.requests[0:0]
	} else {
		s.requests = s.requests[1:]
	}
	return u
}
func (s *RequestScheduler) PollN(n int) ([]Request, int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.requests == nil {
		s.requests = make([]Request, 0)
	}
	m := len(s.requests)
	if m == 0 {
		return nil, m
	}
	if n >= m {
		u := s.requests[0:m]
		s.requests = s.requests[0:0]
		return u, m
	} else {
		u := s.requests[0:n]
		s.requests = s.requests[n:m]
		return u, n
	}
}
func (s *RequestScheduler) Len() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.requests)
}
