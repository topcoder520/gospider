package gospider

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

//ProxyProvider 代理提供器
type ProxyProvider interface {
	GetProxy() *Proxy
	AddProxy(pxy ...Proxy)
}

//SimpleProxyProvider  ProxyProvider默认实现
type SimpleProxyProvider struct {
	proxies []Proxy
	pointer int32
	mux     sync.Mutex
}

//GetProxy 实现ProxyProvider接口
func (sp *SimpleProxyProvider) GetProxy() *Proxy {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	if len(sp.proxies) == 0 {
		return nil
	}
	p := sp.incrForLoop()
	pxy := sp.proxies[p]
	return &pxy
}

func (sp *SimpleProxyProvider) AddProxy(pxy ...Proxy) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	sp.proxies = append(sp.proxies, pxy...)
}

func (sp *SimpleProxyProvider) incrForLoop() int32 {
	sp.pointer = sp.pointer + 1
	size := int32(len(sp.proxies))
	if sp.pointer >= size {
		sp.pointer = 0
	}
	return sp.pointer
}

//Proxy 代理对象
type Proxy struct {
	Scheme   string
	Host     string
	Port     string
	Username string
	Password string
}

func CreateProxy(url url.URL) Proxy {
	proxy := Proxy{
		Scheme: url.Scheme,
		Host:   url.Host,
		Port:   url.Port(),
	}
	if url.User != nil {
		proxy.Username = url.User.Username()
		pwd, isSet := url.User.Password()
		if isSet {
			proxy.Password = pwd
		}
	}
	return proxy
}

func (pxy Proxy) String() string { //http://huangj:123456@golang.cn:8000
	userInfo := ""
	if len(strings.TrimSpace(pxy.Username)) > 0 && len(strings.TrimSpace(pxy.Password)) > 0 {
		userInfo = fmt.Sprintf("%s%s%s%s", pxy.Username, ":", pxy.Password, "@")
	}

	if len(strings.TrimSpace(pxy.Port)) > 0 && !strings.Contains(pxy.Port, ":") {
		if !strings.HasSuffix(pxy.Host, ":"+pxy.Port) {
			pxy.Port = ":" + pxy.Port
		} else {
			pxy.Port = ""

		}
	}
	return fmt.Sprintf("%s%s%s%s%s", pxy.Scheme, "://", userInfo, pxy.Host, pxy.Port)
}
