package gospider

import (
	"context"
	"fmt"
)

type Handler interface {
	//处理结果写入handleResult
	//返回 false则不处理
	Handle(resp Response, handleResult *Result, ctx context.Context) error
}

//Handler的默认实现
type ConsoleHandler struct {
}

//处理结果写入handleResult
//返回 false则不处理
func (hh *ConsoleHandler) Handle(resp Response, handleResult *Result, ctx context.Context) error {
	fmt.Println("Status: ", resp.Status)
	fmt.Println(resp.Body)
	handleResult.AddItem("name", "小明")
	return nil
	//return ErrorSkip //skip,don`t handle
	//return errors.New("has a error")
}
