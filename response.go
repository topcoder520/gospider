package gospider

type Response struct {
	Body          []byte
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	Header        map[string][]string
	ContentLength int64
	Request       *Request
}
