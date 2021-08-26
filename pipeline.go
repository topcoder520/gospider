package gospider

import "context"

type Pipeline interface {
	Process(handleResult *Result, ctx context.Context) error
}
