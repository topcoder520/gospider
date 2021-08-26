package gospider

import "context"

//监听接口
type Listener interface {
	OnError(req Request, e error, ctx context.Context)
	OnSuccess(req Request, ctx context.Context)
}
