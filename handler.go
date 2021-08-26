package gospider

import "context"

type Handler interface {
	//处理结果写入handleResult
	//返回 false则不处理
	Handle(resp Response, handleResult *Result, ctx context.Context) error
}
