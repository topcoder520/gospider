package gospider

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Downloader interface {
	Download(req *Request, ctx context.Context) (resp *Response, err error)
	SetClientGenerator(generator ClientGenerator)
}

//Downloader的默认实现
type HttpDownloader struct {
	generator ClientGenerator
}

func (d *HttpDownloader) SetClientGenerator(generator ClientGenerator) {
	d.generator = generator
}

func (d *HttpDownloader) Download(request *Request, ctx context.Context) (r *Response, err error) {
	request.State = RequestError
	r = &Response{}
	r.Request = request
	urlstr := strings.TrimSpace(request.Url)
	if len(urlstr) == 0 {
		err = errors.New("url can not empty")
		return
	}
	u, err := url.Parse(urlstr)
	if err != nil {
		return
	}
	req, err := http.NewRequest(request.Method, u.String(), nil)
	if err != nil {
		return
	}
	req.Header = http.Header(request.Header)
	for k, v := range request.Header {
		if strings.ToLower(k) == "host" {
			if len(v) > 0 {
				req.Host = v[0]
			}
		}
	}
	client := d.generator.Generate()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	r.StatusCode = resp.StatusCode
	r.Header = resp.Header
	r.ContentLength = resp.ContentLength
	r.Status = resp.Status
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	r.Body = b
	request.State = RequestSuccess
	return
}
