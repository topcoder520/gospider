package gospider

import (
	"context"
	"fmt"
)

type Pipeline interface {
	Process(handleResult *Result, ctx context.Context) error
}

//Pipeline的默认实现
type ConsolePipeline struct {
}

func (c *ConsolePipeline) Process(handleResult *Result, ctx context.Context) error {
	for k, v := range handleResult.TargetItems {
		fmt.Println(k, v)
	}
	return nil
}
