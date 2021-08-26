package gospider

type Scheduler interface {
	Push(requests ...Request)
	Poll() Request
	//取n个 返回数据切片和取到的真是数量
	PollN(n int) ([]Request, int)
	Len() int
	//去重方法
}
