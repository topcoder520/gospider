package gospider

import (
	"context"
	"fmt"
)

type ConsolePipeline struct {
}

func (c *ConsolePipeline) Process(handleResult *Result, ctx context.Context) error {
	for k, v := range handleResult.TargetItems {
		fmt.Println(k, v)
	}
	return nil
}
