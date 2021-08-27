package gospider

import (
	"fmt"
	"testing"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestExample_1(t *testing.T) {
	fmt.Println("start spider....")
	spider := NewSpider("https://studygolang.com/pkgdoc")
	spider.Run()
}

func TestExample_2(t *testing.T) {
	fmt.Println("start spider....")
	spider := NewSpider("https://studygolang.com/pkgdoc")
	spider.SetTimeOut(10 * time.Second) //10s 后退出
	spider.AddSeedUrl("https://studygolang.com/pkgdoc")
	spider.SetSleepTime(1 * time.Second)
	spider.AddHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	spider.SetGoroutines(1)
	spider.PreHandleRequest = func(req *Request) {
		//req.Skip = true
		req.Extras["info"] = "8888888"
	}
	//spider.AddListener()  //添加监听器，监听器必须实现 Listener 接口,可添加多个
	//spider.SetDownloader() //设置下载器,下载器必须实现 Downloader 接口，不设置则使用默认下载器
	//spider.AddHandler() //添加处理响应的处理器,处理器必须实现 Handler 接口，不设置则使用默认的,可添加多个
	//spider.AddPipeline() //添加结果处理器,处理器必须实现 Pipeline 接口，负责处理Handler返回的数据，不设置则使用默认的,可添加多个
	//spider.SetScheduler() //设置调度器,调度器必须实现 Scheduler 接口，不设置则使用默认的
	spider.Run()
}

func TestLevelDB(t *testing.T) {
	db, err := leveldb.OpenFile("./data/db", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	err = db.Put([]byte("key"), []byte("value"), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := db.Get([]byte("key"), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("--------------------result-----------------")
	fmt.Println(string(data))
}
