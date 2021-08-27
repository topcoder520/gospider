package gospider

import "encoding/json"

func RequestStringify(req Request) (string, error) {
	req.Downloader = nil
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ParseRequest(str string) (*Request, error) {
	req := &Request{}
	err := json.Unmarshal([]byte(str), req)
	return req, err
}
