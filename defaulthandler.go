package gospider

import (
	"context"
	"fmt"
)

type HtmlHandler struct {
}

//处理结果写入handleResult
//返回 false则不处理
func (hh *HtmlHandler) Handle(html string, handleResult *HandlerResult, ctx context.Context) error {
	fmt.Println(html)
	return nil
}
