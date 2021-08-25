package gospider

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpDownloader struct {
	HttpClient http.Client
	Header     map[string][]string //设置请求头
}

func NewDownloader() Downloader {
	return &HttpDownloader{
		HttpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (d *HttpDownloader) SetHeader(header map[string][]string) {
	d.Header = header
}

func (d *HttpDownloader) Download(urlstr string, ctx context.Context) (html string, err error) {
	urlstr = strings.TrimSpace(urlstr)
	if len(urlstr) == 0 {
		err = errors.New("url can not empty")
		return
	}
	u, err := url.Parse(urlstr)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return
	}
	req.Header = http.Header(d.Header)
	resp, err := d.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html = string(b)
	return
}
