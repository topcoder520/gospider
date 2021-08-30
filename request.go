package gospider

import (
	"encoding/json"
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
	State      RequestState           //请求的状态
}

func (req *Request) GetExtras(key string) interface{} {
	v, ok := req.Extras[key]
	if ok {
		return v
	}
	return nil
}

func (req *Request) AddExtras(key string, value interface{}) {
	req.Extras[key] = value
}

type RequestState string

const (
	RequestNormal  RequestState = "normal"
	RequestSuccess RequestState = "success"
	RequestError   RequestState = "error"
)

func NewRequest() Request {
	u4 := uuid.NewV4()
	return Request{
		Id:     u4.String(),
		Method: http.MethodGet,
		Header: map[string][]string{},
		Extras: map[string]interface{}{},
		State:  RequestNormal,
		Skip:   false,
	}
}

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
