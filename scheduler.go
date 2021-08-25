package gospider

type Scheduler interface {
	Push(url ...string)
	Poll() string
	//取n个 返回数据切片和取到的真是数量
	PollN(n int) ([]string, int)
	Len() int
	//去重方法
}
