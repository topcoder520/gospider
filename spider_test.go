package gospider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func TestExample_1(t *testing.T) {
	fmt.Println("start spider....")
	spider := NewSpider("https://www.w3school.com.cn/tags/tag_html.asp")
	spider.ClearStoreDB()
	//spider.SaveHtml(true, "./data/html/", nil)
	spider.SaveHtml("./data/html/", func() string {
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
	spider.ClearStoreDB()
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

type Book struct {
	Id             string
	Url            string
	Title          string
	Author         string
	LastUpdateTime string
	Intro          string
}

type BookHandler struct {
}

func (h *BookHandler) Handle(resp Response, handleResult *Result, ctx context.Context) error {
	//resp.Request.Url
	isMatch, err := regexp.MatchString(`https://www.bqg74.com/0_(\d+)/`, resp.Request.Url)
	if err != nil {
		return err
	}
	if !isMatch {
		return ErrorSkip
	}
	if strings.HasSuffix(resp.Request.Url, ".html") {
		return ErrorSkip
	}
	fmt.Println(string(resp.Request.Url))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
	if err != nil {
		log.Fatal(err)
	}
	book := Book{}
	book.Title = doc.Find("div#info>h1").Text()
	book.Url = resp.Request.Url
	book.Id = fmt.Sprintf("%d", time.Now().Unix())
	book.Author = ""
	doc.Find("div#info>p").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			book.Author = s.Text()
		} else if i == 2 {
			book.LastUpdateTime = s.Text()
		}
	})
	book.Intro = doc.Find("div#intro").Text()
	//fmt.Println(book)
	handleResult.AddItem("book", book)

	listRequest := make([]Request, 0)
	doc.Find("div#list>dl>dd>a").Each(func(i int, s *goquery.Selection) {
		chapter := s.Text()
		chapterUri, exits := s.Attr("href")
		if exits {
			request := NewRequest()
			request.Url = book.Url + chapterUri
			request.Extras["title"] = chapter
			request.Extras["bookTitle"] = book.Title
			request.Extras["bookId"] = book.Id
			listRequest = append(listRequest, request)
		}
	})
	handleResult.TargetRequests = listRequest
	return nil
}

type Chapter struct {
	BookTitle string
	BookId    string
	Title     string
	Url       string
	PageIndex string
	TotalPage string
	Content   string
}

type ChapterHandler struct{}

func (h *ChapterHandler) Handle(resp Response, handleResult *Result, ctx context.Context) error {
	isMatch, err := regexp.MatchString(`https://www.bqg74.com/0_(\d+)/`, resp.Request.Url)
	if err != nil {
		return err
	}
	if !isMatch {
		return ErrorSkip
	}
	if !strings.HasSuffix(resp.Request.Url, ".html") {
		return ErrorSkip
	}
	fmt.Println(string(resp.Request.Url))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
	if err != nil {
		log.Fatal(err)
	}
	chapter := Chapter{}
	chapter.Title = doc.Find("div.bookname>h1").Text()
	chapter.Url = resp.Request.Url
	chapter.Content = doc.Find("div#content").Text()
	rex := regexp.MustCompile(`(第(\d+)/(\d+))页`)
	arr := rex.FindStringSubmatch(chapter.Title)
	if len(arr) == 0 {
		return ErrorSkip
	}
	chapter.PageIndex = arr[2]
	chapter.TotalPage = arr[3]
	bookTitleInf := resp.Request.GetExtras("bookTitle")
	if bookTitleInf != nil {
		chapter.BookTitle = bookTitleInf.(string)
	}
	bookIdInf := resp.Request.GetExtras("bookId")
	if bookIdInf != nil {
		chapter.BookId = bookIdInf.(string)
	}
	handleResult.AddItem("chapter", chapter)

	intTotalPage, _ := strconv.Atoi(chapter.TotalPage)
	intPageIndex, _ := strconv.Atoi(chapter.PageIndex)
	if (intPageIndex + 1) <= intTotalPage {

		baseName := path.Base(chapter.Url)
		ext := path.Ext(baseName)
		index := strings.LastIndex(baseName, ext)
		baseName = baseName[0 : index-1]
		baseName = fmt.Sprintf("%s_%d%s", baseName, (intPageIndex + 1), ext)
		curl := chapter.Url[0:(strings.LastIndex(chapter.Url, "/")+1)] + baseName
		request := NewRequest()
		request.Extras["title"] = chapter.Title
		request.Extras["bookTitle"] = chapter.BookTitle
		request.Extras["bookId"] = chapter.BookId
		request.Url = curl
		handleResult.AddTargetRequest(request)
	}

	return nil
}

type BookPipeline struct {
	spider *Spider
}

func BookStringify(book Book) (string, error) {
	b, err := json.Marshal(book)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ChapterStringify(book Chapter) (string, error) {
	b, err := json.Marshal(book)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *BookPipeline) Process(handleResult *Result, ctx context.Context) error {
	bookInf, ok := handleResult.TargetItems["book"]
	if ok {
		book := bookInf.(Book)
		bookstr, err := BookStringify(book)
		if err != nil {
			return err
		}
		for _, sotre := range p.spider.RequestsStore {
			sotre.Add(fmt.Sprintf("book-%s", book.Id), bookstr)
		}
		return nil
	}
	chapterInf, ok := handleResult.TargetItems["chapter"]
	if ok {
		chapter := chapterInf.(Chapter)
		chapterstr, err := ChapterStringify(chapter)
		if err != nil {
			return err
		}
		for _, sotre := range p.spider.RequestsStore {
			sotre.Add(fmt.Sprintf("book-chapter-%s-%s-%s", chapter.BookId, chapter.Title, chapter.PageIndex), chapterstr)
		}
	}
	return nil
}

func TestBiquge(t *testing.T) {

	spider := NewSpider("https://www.bqg74.com/0_11/")
	spider.AddHeader("Host", "www.bqg74.com")
	spider.AddHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	spider.AddHeader("Cookie", "Hm_lvt_15ab8d1f6950da8c0dd49a424be8cbdb=1630327189; Hm_lpvt_15ab8d1f6950da8c0dd49a424be8cbdb=1630339530")
	//spider.ClearStoreDB()
	spider.SetByteHandler(&GBKByteHandler{})
	spider.SetGoroutines(5)
	spider.SetSleepTime(2 * time.Second)
	//spider.SaveHtml("./data/html", nil)
	spider.AddHandler(&BookHandler{})
	spider.AddHandler(&ChapterHandler{})
	spider.AddPipeline(&BookPipeline{spider: spider})
	spider.Run()
}

func ParseBook(str string) (*Book, error) {
	book := &Book{}
	err := json.Unmarshal([]byte(str), book)
	return book, err
}

func TestBiqugeStore(t *testing.T) {
	levelStore := RequestStore{}
	path := "./data/db/www.bqg74.com"
	levelStore.dataDB = CreateDataDB(path)
	listBook, err := levelStore.List("book-")
	if err != nil {
		log.Println(err)
		return
	}
	for _, str := range listBook {
		fmt.Println(str)
	}
}

func TestRegex(t *testing.T) {
	isMatch, err := regexp.MatchString(`https://www.bqg74.com/0_(\d+)/`, "https://www.bqg74.com/0_0100/4545.html")
	if err == nil && isMatch {
		fmt.Println("匹配成功")
	}
	fmt.Println(isMatch, err)

	r, _ := regexp.Compile("p([a-z]+)ch")
	b := r.MatchString("peach.html")
	fmt.Println(b) //结果为true
}

func TestRegex2(t *testing.T) {
	title := "第000章 天空一声雷响，老子闪亮登场！(第2/08页)"
	isMatch, err := regexp.MatchString(`(第(\d+)/(\d+)页)`, title)
	if err == nil && isMatch {
		fmt.Println("匹配成功")
	}
	rex := regexp.MustCompile(`(第(\d+)/(\d+))页`)
	b := rex.FindStringSubmatch(title)
	for _, c := range b {
		fmt.Println(c)
	}
}
