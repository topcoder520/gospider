package gospider

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"time"
)

//ClientGenerator 客户端生成器
type ClientGenerator interface {
	Generate() *http.Client
	SetProxyProvider(pxyProvider ProxyProvider)
}

//SimpleClientGenerator ClientGenerator的默认实现
type SimpleClientGenerator struct {
	pxyProvider ProxyProvider
}

func (sg *SimpleClientGenerator) Generate() *http.Client {
	if sg.pxyProvider == nil {
		return &http.Client{
			Timeout: 10 * time.Second, // 10s 超时
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //跳过证书
				},
			},
		}
	}
	proxy := sg.pxyProvider.GetProxy()
	if proxy == nil {
		return &http.Client{
			Timeout: 10 * time.Second, // 10s 超时
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //跳过证书
				},
			},
		}
	}
	proxyUrl, err := url.Parse(proxy.String())
	if err != nil {
		log.Println("SimpleClientGenerator.Generate() err: ", err)
		return nil
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //跳过证书验证
			},
			Proxy: http.ProxyURL(proxyUrl), //设置代理
		},
		Timeout: 10 * time.Second, // 10s 超时,
	}
}

func (sg *SimpleClientGenerator) SetProxyProvider(pxyProvider ProxyProvider) {
	sg.pxyProvider = pxyProvider
}
