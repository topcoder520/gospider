package main

import "context"

type Handler interface {
	//处理结果写入handleResult
	//返回 false则不处理
	Handle(html string, handleResult *HandlerResult, ctx context.Context) error
}
