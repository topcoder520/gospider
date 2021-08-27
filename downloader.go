package gospider

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Downloader interface {
	Download(req *Request, ctx context.Context) (resp *Response, err error)
}

//Downloader的默认实现
type HttpDownloader struct {
	HttpClient http.Client
}

func NewDownloader() *HttpDownloader {
	return &HttpDownloader{
		HttpClient: http.Client{
			Timeout: 10 * time.Second, // 10s 超时
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //跳过证书
				},
			},
		},
	}
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
	resp, err := d.HttpClient.Do(req)
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
	r.Body = string(b)
	request.State = RequestSuccess
	return
}
