package gospider

import "net/http"

type Request struct {
	Url        string                 //请求资源地址
	Method     string                 //请求方法,
	Header     map[string][]string    //请求头
	Downloader Downloader             //下载器
	Extras     map[string]interface{} //额外信息
	Skip       bool                   //跳过请求不处理
}

func NewRequest() Request {
	return Request{
		Method: http.MethodGet,
		Header: map[string][]string{},
		Extras: map[string]interface{}{},
	}
}
