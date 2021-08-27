package gospider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/topcoder520/gospider/data"
)

type Spider struct {
	downloader       Downloader          //下载器 负责下载网页
	listHandler      []Handler           //处理器 负责处理网页
	listPipeline     []Pipeline          //管道 负责持久化数据或者下载资源的任务
	scheduler        Scheduler           //调度器 负责待爬取的url的管理
	sleepTime        time.Duration       //控制访问的速度，单个协程每执行一次沉睡sleepTime
	goroutines       int                 //开启协程数量
	header           map[string][]string //设置请求头
	initRequests     []Request           //种子url
	timeOut          time.Duration       //没有数据的情况下，程序结束运行的时间
	isTimeOut        bool                //没有数据的情况下是否自动退出，默认true
	listListener     []Listener          //程序监听器
	PreHandleRequest RequestHandle       //执行请求前的请求处理
	totalPage        int                 //爬取的总链接
	totalPageMux     sync.RWMutex
	store            data.Store //保存请求对象数据
}

//NewSpider 创建一个爬虫程序
//seedUrl 种子Url
func NewSpider(seedUrl ...string) *Spider {
	spider := &Spider{
		downloader:   NewDownloader(),
		scheduler:    &RequestScheduler{},
		listHandler:  make([]Handler, 0),
		listPipeline: make([]Pipeline, 0),
		sleepTime:    time.Second * 1, //默认1s
		goroutines:   1,               //默认开1个
		header:       make(map[string][]string),
		initRequests: make([]Request, 0),
		isTimeOut:    true,
		timeOut:      10 * time.Second,
		listListener: make([]Listener, 0),
	}
	spider.checkUrls(seedUrl)
	return spider
}

//SetTimeOut 设置程序在没有数据之后退出的时间
//当t<0时 程序一直运行不退出
func (s *Spider) SetTimeOut(t time.Duration) {
	if t < 0 {
		s.isTimeOut = false
	} else {
		s.timeOut = t
	}
}

func (s *Spider) checkUrls(urls []string) {
	if len(urls) == 0 {
		return
	}
	for i := 0; i < len(urls); i++ {
		u, err := url.Parse(urls[i])
		if err != nil {
			panic(err.Error())
		}
		req := NewRequest()
		req.Url = u.String()
		req.Method = http.MethodGet
		req.Header = s.header
		req.Skip = false
		s.initRequests = append(s.initRequests, req)
	}
}

//AddInitUrl 添加种子链接
func (s *Spider) AddSeedUrl(seedUrls ...string) {
	s.checkUrls(seedUrls)
}

//SetSleepTime 睡眠时间
func (s *Spider) SetSleepTime(t time.Duration) {
	s.sleepTime = t
}

//AddHeader 添加请求头
func (s *Spider) AddHeader(key, value string) {
	v, ok := s.header[key]
	if ok {
		v = append(v, value)
	} else {
		v = make([]string, 1)
		v[0] = value
	}
	s.header[key] = v
}

//SetGoroutines 协程数
func (s *Spider) SetGoroutines(n int) {
	if n <= 0 {
		n = 1
	}
	s.goroutines = n
}

//SetDownloader 设置下载器
func (s *Spider) SetDownloader(downloader Downloader) {
	s.downloader = downloader
}

//AddHandler 添加处理器
func (s *Spider) AddHandler(handler Handler) {
	s.listHandler = append(s.listHandler, handler)
}

//AddPipeline 添加结果处理器
func (s *Spider) AddPipeline(pipeline Pipeline) {
	s.listPipeline = append(s.listPipeline, pipeline)
}

//SetScheduler 设置调度器
func (s *Spider) SetScheduler(scheduler Scheduler) {
	s.scheduler = scheduler
}

//AddListener 添加监听器
func (s *Spider) AddListener(listener Listener) {
	s.listListener = append(s.listListener, listener)
}

//initCompent init
func (s *Spider) initCompent() {
	if _, ok := s.header["User-Agent"]; !ok {
		s.AddHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	}
	if s.downloader == nil {
		httpDownloader := NewDownloader()
		s.downloader = httpDownloader
	}

	if s.scheduler == nil {
		s.scheduler = &RequestScheduler{}
	}
	if len(s.initRequests) == 0 {
		panic("seed url can not empty")
	}
	s.scheduler.Push(s.initRequests...)
	s.totalPage = len(s.initRequests)

	if s.listHandler == nil {
		s.listHandler = make([]Handler, 0, 1)
		s.listHandler = append(s.listHandler, &ConsoleHandler{})
	} else if len(s.listHandler) == 0 {
		s.listHandler = append(s.listHandler, &ConsoleHandler{})
	}
	if s.listPipeline == nil {
		s.listPipeline = make([]Pipeline, 0, 1)
		s.listPipeline = append(s.listPipeline, &ConsolePipeline{})
	} else if len(s.listPipeline) == 0 {
		s.listPipeline = append(s.listPipeline, &ConsolePipeline{})
	}
	if s.goroutines == 0 {
		s.goroutines = 1
	}
	if s.sleepTime == time.Second*0 {
		s.sleepTime = time.Second * 2
	}
	if s.store == nil {
		s.store = data.CreateLeveldbStore("./data/db")
	}
}

type RequestHandle func(req *Request)

//Run 运行
func (s *Spider) Run() {
	s.initCompent()
	task := make(chan Request, 200)
	interruptChan := make(chan int, s.goroutines+1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//分配任务
	go func(c context.Context) {
		for {
			select {
			case <-c.Done():
				close(task)
				return
			case <-interruptChan:
				close(task)
				return
			default:
				l := s.scheduler.Len()
				if l == 0 {
					time.Sleep(time.Second * 1)
					continue
				}
				u := s.scheduler.Poll()
				task <- u
			}
		}
	}(ctx)
	//取任务
	wg := sync.WaitGroup{}
	for i := 0; i < s.goroutines; i++ {
		wg.Add(1)
		go func(index int, c context.Context) {
			defer wg.Done()
			for {
				select {
				case <-c.Done():
					log.Printf("协程 %d 程序结束\n", index)
					return
				case <-interruptChan:
					for req := range task {
						reqStr, err := RequestStringify(req)
						if err != nil {
							s.store.Add(fmt.Sprintf("req-%s-%s", "NO", req.Id), reqStr)
						}
					}
					for s.scheduler.Len() > 0 {
						req := s.scheduler.Poll()
						reqStr, err := RequestStringify(req)
						if err != nil {
							s.store.Add(fmt.Sprintf("req-%s-%s", "NO", req.Id), reqStr)
						}
					}
					return
				case req := <-task:
					req.Downloader = s.downloader
					req.Header = s.header
					if s.PreHandleRequest != nil {
						s.PreHandleRequest(&req)
					}
					if req.Skip {
						return
					}
					err := s.handRequest(&req, c)
					if s.listListener != nil && len(s.listListener) > 0 {
						for _, listener := range s.listListener {
							s.totalPageMux.RLock()
							req.Extras["totalPage"] = s.totalPage
							s.totalPageMux.RUnlock()
							if err != nil {
								listener.OnError(req, err, ctx)
							} else {
								listener.OnSuccess(req, ctx)
							}
						}
					}
					if err != nil {
						log.Println("handle request err: ", err)
					}
					time.Sleep(s.sleepTime)
				case <-time.After(s.timeOut):
					if s.isTimeOut {
						log.Printf("协程 %d 程序结束\n", index)
						return
					}
				}
			}
		}(i, ctx)
	}
	//信号
	go func(c context.Context) {
		for {
			select {
			case <-c.Done():
				return
			case <-signalChan:
				fmt.Println("信号打断")
				for i := 0; i < cap(interruptChan); i++ {
					interruptChan <- 1
				}
				return
			}
		}
	}(ctx)
	wg.Wait()
	fmt.Println("爬虫结束运行!")
}

func (s *Spider) handRequest(req *Request, ctx context.Context) (err error) {
	resp, err := s.downloader.Download(req, ctx)
	if err != nil {
		return
	}
	lenHandle := len(s.listHandler)
	lenPipeline := len(s.listPipeline)
	for i := 0; i < lenHandle; i++ {
		result := &Result{make([]Request, 0), make(map[string]interface{})}
		err = s.listHandler[i].Handle(*resp, result, ctx)
		if err != nil {
			if err == ErrorSkip {
				continue
			}
			return
		}
		if result.TargetRequests != nil && len(result.TargetRequests) > 0 {
			s.scheduler.Push(result.TargetRequests...)
			s.totalPageMux.Lock()
			s.totalPage = s.totalPage + len(result.TargetRequests)
			s.totalPageMux.Unlock()
		}
		for j := 0; j < lenPipeline; j++ {
			err = s.listPipeline[j].Process(result, ctx)
			if err != nil {
				if err == ErrorSkip {
					continue
				}
				return
			}
		}
	}
	return
}

//SetStoreDB 存储器 存储请求数据
func (s *Spider) SetStoreDB(store data.Store) {
	s.store = store
}
