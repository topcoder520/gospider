package gospider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Spider struct {
	listProxy       []Proxy             //代理
	proxyProvider   ProxyProvider       //代理提供器
	clientGenerator ClientGenerator     //客户端生成器
	downloader      Downloader          //下载器 负责下载网页
	listHandler     []Handler           //处理器 负责处理网页
	listPipeline    []Pipeline          //管道 负责持久化数据或者下载资源的任务
	scheduler       Scheduler           //调度器 负责待爬取的url的管理
	sleepTime       time.Duration       //控制访问的速度，单个协程每执行一次沉睡sleepTime
	goroutines      int                 //开启协程数量
	header          map[string][]string //设置请求头
	initRequests    []Request           //种子url
	timeOut         time.Duration       //没有数据的情况下，程序结束运行的时间
	isTimeOut       bool                //没有数据的情况下是否自动退出，默认true
	listListener    []Listener          //程序监听器
	isSaveHtml      bool                //是否把下载的html页面保存下来,默认不保存
	saveHtmlPath    string              //html页面数据保存地址
	requestFilter   RequestFilter       //过滤重复请求
	isClearStoreDB  bool                //是否清空存储的数据
	suffixGenerate  func() string       //名字的后缀生成函数
	byteHandler     ByteHandler         //字节处理
	cycleTime       int                 //请求失败之后重复请求的次数

	RequestsStore    []Store       //保存请求对象数据
	PreHandleRequest RequestHandle //执行请求前的请求处理
}

//NewSpider 创建一个爬虫程序
//seedUrl 种子Url
func NewSpider(seedUrl ...string) *Spider {
	spider := &Spider{
		scheduler:     &RequestScheduler{},
		listProxy:     make([]Proxy, 0),
		listHandler:   make([]Handler, 0),
		listPipeline:  make([]Pipeline, 0),
		sleepTime:     time.Second * 1,  //默认1s
		goroutines:    runtime.NumCPU(), //默认主机cpu核数
		header:        make(map[string][]string),
		initRequests:  make([]Request, 0),
		isTimeOut:     true,
		timeOut:       10 * time.Second,
		listListener:  make([]Listener, 0),
		RequestsStore: make([]Store, 0, 1),
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

func (s *Spider) getDoman(u string) (string, error) {
	urlObj, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return urlObj.Hostname(), nil
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

//AddProxy 添加代理
func (s *Spider) AddProxy(pxy ...Proxy) {
	s.listProxy = append(s.listProxy, pxy...)
}

//SetProxyProvider 代理提供者
func (s *Spider) SetProxyProvider(proxyProvider ProxyProvider) {
	s.proxyProvider = proxyProvider
}

//SetClientGenerator 客户端生成器
func (s *Spider) SetClientGenerator(clientGenerator ClientGenerator) {
	s.clientGenerator = clientGenerator
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

//SetRequestFilter 设置请求过滤器
func (s *Spider) SetRequestFilter(filter RequestFilter) {
	s.requestFilter = filter
}

//SetByteHandler 设置字节处理器 对下载的字节进行处理
func (s *Spider) SetByteHandler(handler ByteHandler) {
	s.byteHandler = handler
}

//SetCycleTime 设置请求失败后重复请求次数
func (s *Spider) SetCycleTime(time int) {
	s.cycleTime = time
}

//initCompent init
func (s *Spider) initCompent() {
	if _, ok := s.header["User-Agent"]; !ok {
		s.AddHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	}
	if s.proxyProvider == nil {
		pxyProvider := &SimpleProxyProvider{
			proxies: make([]Proxy, 0),
		}
		s.proxyProvider = pxyProvider
	}
	if len(s.listProxy) > 0 {
		s.proxyProvider.AddProxy(s.listProxy...)
	}
	if s.clientGenerator == nil {
		sg := &SimpleClientGenerator{}
		s.clientGenerator = sg
	}
	s.clientGenerator.SetProxyProvider(s.proxyProvider)
	if s.downloader == nil {
		dlr := &HttpDownloader{}
		s.downloader = dlr
	}
	s.downloader.SetClientGenerator(s.clientGenerator)

	if s.scheduler == nil {
		s.scheduler = &RequestScheduler{}
	}
	if len(s.initRequests) == 0 {
		panic("seed url can not empty")
	}

	s.listListener = append(s.listListener, &DefaultListener{spider: s})

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
	if len(s.RequestsStore) == 0 {
		dataPath := "./data/db/requestdb/" //默认地址
		reqStore := &DataStore{}
		reqStore.dataDB = CreateDataDB(dataPath)
		s.RequestsStore = append(s.RequestsStore, reqStore)
	}
	if s.isClearStoreDB {
		for _, store := range s.RequestsStore {
			store.Clear("")
		}
	}
	if s.requestFilter == nil {
		requestFilter := &StoreRequestFilter{
			mapRequest: make(map[string]Request),
		}
		for _, store := range s.RequestsStore {
			listRequest, err := store.List("")
			if err != nil {
				log.Println("store.List err: ", err)
				continue
			}
			for _, reqStr := range listRequest {
				req, err := ParseRequest(reqStr)
				if err != nil {
					log.Println("ParseRequest err", err)
					continue
				}
				requestFilter.mapRequest[req.Url] = *req
			}
		}
		s.requestFilter = requestFilter
	}
	//
	s.scheduler.Push(s.requestFilter.Filter(s.initRequests...)...)
	for _, store := range s.RequestsStore {
		listRequest, err := store.List(string(RequestNormal))
		if err != nil {
			log.Println("store.List err: ", err)
			continue
		}
		for _, reqStr := range listRequest {
			req, err := ParseRequest(reqStr)
			if err != nil {
				log.Println("ParseRequest err", err)
				continue
			}
			s.scheduler.Push(*req)
		}
	}
}

func (s *Spider) saveRequest(req *Request, state RequestState) {
	for _, store := range s.RequestsStore {
		reqStr, err := RequestStringify(*req)
		if err != nil {
			log.Println(err)
			return
		}
		store.Add(fmt.Sprintf("%s-%s", state, req.Id), reqStr)
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
					for req := range task {
						s.saveRequest(&req, RequestNormal) //"NO" not downloaded
					}
					for s.scheduler.Len() > 0 {
						req := s.scheduler.Poll()
						s.saveRequest(&req, RequestNormal) //"NO" not downloaded
					}
					log.Printf("协程 %d 程序结束\n", index)
					return
				case <-interruptChan:
					for req := range task {
						s.saveRequest(&req, RequestNormal) //"NO" not downloaded
					}
					for s.scheduler.Len() > 0 {
						req := s.scheduler.Poll()
						s.saveRequest(&req, RequestNormal) //"NO" not downloaded
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
					if err != nil {
						log.Println("handle request err: ", err)
					}
					if s.listListener != nil && len(s.listListener) > 0 {
						for _, listener := range s.listListener {
							if err != nil {
								listener.OnError(req, err, ctx)
							} else {
								listener.OnSuccess(req, ctx)
							}
						}
					}
					time.Sleep(s.sleepTime)
				case <-time.After(s.timeOut):
					if s.isTimeOut {
						cancel()
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
	//处理字节
	if s.byteHandler != nil {
		handleByte, err := s.byteHandler.Handle(resp.Body)
		if err != nil {
			log.Println("byteHandler.Handle err: ", err)
		} else {
			resp.Body = handleByte
		}
	}

	lenHandle := len(s.listHandler)
	lenPipeline := len(s.listPipeline)
	for i := 0; i < lenHandle; i++ {
		result := &Result{make([]Request, 0), make(map[string]interface{})}
		//处理响应
		err = s.listHandler[i].Handle(*resp, result, ctx)
		if err != nil {
			if err == ErrorSkip {
				continue
			}
			return
		}
		if result.TargetRequests != nil && len(result.TargetRequests) > 0 {
			reqs := s.requestFilter.Filter(result.TargetRequests...)
			s.scheduler.Push(reqs...)
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
	if s.isSaveHtml {
		go func(p string, r Request, sp Response, c context.Context) {
			select {
			case <-c.Done():
				return
			default:
				doman, _ := s.getDoman(req.Url)
				absPath, errr := filepath.Abs(filepath.Clean(p))
				if errr != nil {
					return
				}
				absPath = filepath.Join(absPath, doman)
				os.MkdirAll(absPath, 0777)

				if s.suffixGenerate != nil {
					baseName := path.Base(r.Url)
					ext := path.Ext(baseName)
					if ext != "" {
						index := strings.LastIndex(baseName, ext)
						baseName = baseName[0:index]
					}
					suffix := s.suffixGenerate()
					absPath = path.Join(absPath, fmt.Sprintf("%s%s", baseName, suffix))
				} else {
					absPath = path.Join(absPath, path.Base(r.Url))
				}

				f, errr := os.Create(absPath)
				if err != nil {
					log.Println("save html Create file err: ", errr)
					return
				}
				defer f.Close()
				_, errr = f.Write(sp.Body)
				if err != nil {
					log.Println("save html WriteString err: ", errr)
					return
				}
			}
		}(s.saveHtmlPath, *req, *resp, ctx)
	}
	return
}

//SetStoreDB 存储器 存储请求数据
func (s *Spider) AddRequestStore(store Store) {
	s.RequestsStore = append(s.RequestsStore, store)
}

//Clear 清楚存储的数据
func (s *Spider) ClearRequestStore() {
	s.isClearStoreDB = true
}

//SaveHtml 是否保存html 默认false不保存
//savepath保存地址
//也可以在自定义的Handler处理器中自行实现保存逻辑
//suffixGenerate 名字后缀函数,html存储名字和生成的后缀拼接
func (s *Spider) SaveHtml(savepath string, suffixGenerate func() string) {
	s.isSaveHtml = true
	s.saveHtmlPath = savepath
	s.suffixGenerate = suffixGenerate
}
