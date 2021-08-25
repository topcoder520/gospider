package gospider

import (
	"fmt"
	"testing"
	"time"
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
	//spider.SetDownloader() //设置自己实现的下载器,下载器必须实现Downloader接口，不设置则使用默认下载器
	//spider.AddHandler() //添加自己实现的处理html文本的处理器,处理器必须实现Handler接口，不设置则使用默认的
	//spider.AddPipeline() //添加自己实现的结果处理器,处理器必须实现Pipeline接口，负责处理Handler返回的数据，不设置则使用默认的
	//spider.SetScheduler() //设置自己实现的调度器,调度器必须实现Scheduler接口，不设置则使用默认的
	spider.Run()
}
