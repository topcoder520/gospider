package gospider

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type Request struct {
	Id         string
	Url        string                 //请求资源地址
	Method     string                 //请求方法,
	Header     map[string][]string    //请求头
	Downloader Downloader             //下载器
	Extras     map[string]interface{} //额外信息
	Skip       bool                   //跳过请求不处理
}

func NewRequest() Request {
	u4 := uuid.NewV4()
	return Request{
		Id:     u4.String(),
		Method: http.MethodGet,
		Header: map[string][]string{},
		Extras: map[string]interface{}{},
	}
}
