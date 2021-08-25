package gospider

import "context"

type Pipeline interface {
	Process(handleResult *HandlerResult, ctx context.Context) error
}
