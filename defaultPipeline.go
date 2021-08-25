package main

import (
	"context"
	"fmt"
)

type ConsolePipeline struct {
}

func (c *ConsolePipeline) Process(handleResult *HandlerResult, ctx context.Context) error {
	fmt.Println("process")
	return nil
}
