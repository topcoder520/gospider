package gospider

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func TestExample_1(t *testing.T) {
	fmt.Println("start spider....")
	spider := NewSpider("https://www.w3school.com.cn/tags/tag_html.asp")
	spider.ClearStoreDB()
	//spider.SaveHtml(true, "./data/html/", nil)
	spider.SaveHtml(true, "./data/html/", func() string {
		return ".html"
	})
	spider.Run()
}

func TestExample_path(t *testing.T) {
	p := "www/kl.html"
	fmt.Println(path.Ext(p))
	fmt.Println(path.Base(p))
}

func TestExample_dbname(t *testing.T) {
	db, err := leveldb.OpenFile("./data/db/studygolang.com/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	data, err := db.Get([]byte(StoreKey), nil)
	fmt.Println(err)
	fmt.Printf("[studygolang.com]:%s\n", string(data))
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

	// 存入数据
	db.Put([]byte("1"), []byte("6"), nil)
	db.Put([]byte("2"), []byte("7"), nil)
	db.Put([]byte("3"), []byte("8"), nil)
	db.Put([]byte("foo-4"), []byte("9"), nil)
	db.Put([]byte("5"), []byte("10"), nil)
	db.Put([]byte("6"), []byte("11"), nil)
	db.Put([]byte("moo-7"), []byte("12"), nil)
	db.Put([]byte("8"), []byte("13"), nil)

	// 遍历数据库内容
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("[%s]:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		panic(err)

	}

	fmt.Println("***************************************************")

	// 删除某条数据
	err = db.Delete([]byte("2"), nil)

	// 读取某条数据
	data, err := db.Get([]byte("2"), nil)
	fmt.Printf("[2]:%s:%s\n", data, err)

	// 根据前缀遍历数据库内容
	fmt.Println("***************************************************")
	iter = db.NewIterator(util.BytesPrefix([]byte("foo-")), nil)
	for iter.Next() {
		fmt.Printf("[%s]:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
	err = iter.Error()

	// 遍历从指定 key
	fmt.Println("***************************************************")
	iter = db.NewIterator(nil, nil)
	for ok := iter.Seek([]byte("5")); ok; ok = iter.Next() {
		fmt.Printf("[%s]:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
	err = iter.Error()

	// 遍历子集范围
	fmt.Println("***************************************************")
	iter = db.NewIterator(&util.Range{Start: []byte("foo"), Limit: []byte("loo")}, nil)
	for iter.Next() {
		fmt.Printf("[%s]:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
	err = iter.Error()

	// 批量操作
	fmt.Println("***************************************************")
	batch := new(leveldb.Batch)
	batch.Put([]byte("foo"), []byte("value"))
	batch.Put([]byte("bar"), []byte("another value"))
	batch.Delete([]byte("baz"))
	err = db.Write(batch, nil)

	// 遍历数据库内容
	iter = db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("[%s]:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
	err = iter.Error()
}

func TestLevelDB2(t *testing.T) {
	db, err := leveldb.OpenFile("./data/db", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// 存入数据
	db.Put([]byte("1"), []byte("6"), nil)
	db.Put([]byte("2"), []byte("7"), nil)
	db.Put([]byte("3"), []byte("8"), nil)
	db.Put([]byte("foo-4"), []byte("9"), nil)
	db.Put([]byte("foo-4loo"), []byte("19"), nil)
	db.Put([]byte("foo-4loo1"), []byte("619"), nil)
	db.Put([]byte("5"), []byte("10"), nil)
	db.Put([]byte("6"), []byte("11"), nil)
	db.Put([]byte("moo-7loo"), []byte("12l00"), nil)
	db.Put([]byte("8"), []byte("13"), nil)
	// 遍历子集范围
	fmt.Println("***************************************************")
	iter := db.NewIterator(&util.Range{Start: []byte(""), Limit: []byte("loo")}, nil)
	for iter.Next() {
		fmt.Printf("[%s]:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()

}

func TestLevelDB3(t *testing.T) {
	req := Request{
		Url:        "www.baidu.com",
		Method:     "GET",
		Header:     map[string][]string{"k1": []string{"123", "7845"}, "k2": []string{"123", "7845"}},
		Downloader: nil,
		Extras:     map[string]interface{}{"k3": "va4545"},
		Skip:       false,
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	} else {
		s := string(b)
		fmt.Println(s)
		newReq := &Request{}
		err := json.Unmarshal([]byte(s), newReq)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(*newReq)
		}
	}

}

func TestExample_url(t *testing.T) {
	u, err := url.Parse("http://192.168.1.32:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u.Hostname())
}
