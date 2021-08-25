package gospider

import "sync"

type UrlScheduler struct {
	urls []string
	mux  sync.Mutex
}

func (s *UrlScheduler) Push(url ...string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.urls == nil {
		s.urls = make([]string, 0, len(url))
	}
	s.urls = append(s.urls, url...)
}
func (s *UrlScheduler) Poll() string {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.urls == nil {
		s.urls = make([]string, 0)
	}
	m := len(s.urls)
	if m == 0 {
		return ""
	}
	u := s.urls[0]
	if m == 1 {
		s.urls = s.urls[0:0]
	} else {
		s.urls = s.urls[1:]
	}
	return u
}
func (s *UrlScheduler) PollN(n int) ([]string, int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.urls == nil {
		s.urls = make([]string, 0)
	}
	m := len(s.urls)
	if m == 0 {
		return nil, m
	}
	if n >= m {
		u := s.urls[0:m]
		s.urls = s.urls[0:0]
		return u, m
	} else {
		u := s.urls[0:n]
		s.urls = s.urls[n:m]
		return u, n
	}
}
func (s *UrlScheduler) Len() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.urls)
}
