package gospider

import "context"

//监听接口
type Listener interface {
	OnError(req Request, e error, ctx context.Context)
	OnSuccess(req Request, ctx context.Context)
}

//Listener的认实现
type DefaultListener struct {
	spider *Spider
}

func (listen *DefaultListener) OnError(req Request, e error, ctx context.Context) {
	if req.CycleTime < listen.spider.cycleTime {
		req.CycleTime++
		req.State = RequestNormal
		listen.spider.saveRequest(&req, RequestNormal)
		listen.spider.scheduler.Push(req)
	} else {
		listen.spider.saveRequest(&req, RequestError)
	}
}
func (listen *DefaultListener) OnSuccess(req Request, ctx context.Context) {
	listen.spider.saveRequest(&req, RequestSuccess)
}
