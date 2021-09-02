package gospider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestExample1(t *testing.T) {
	spider := NewSpider("https://www.bqg74.com/0_11/")
	spider.AddHeader("Host", "www.bqg74.com")
	spider.AddHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	spider.AddHeader("Cookie", "Hm_lvt_15ab8d1f6950da8c0dd49a424be8cbdb=1630327189; Hm_lpvt_15ab8d1f6950da8c0dd49a424be8cbdb=1630339530")
	//spider.ClearRequestStore()
	spider.SetByteHandler(&GBKByteHandler{})
	spider.SetGoroutines(5)
	spider.SetSleepTime(2 * time.Second)
	//spider.SaveHtml("./data/html", nil)
	spider.AddHandler(&BookHandler{})
	spider.AddHandler(&ChapterHandler{})
	spider.AddPipeline(&BookPipeline{spider: spider, bookStore: &DataStore{dataDB: CreateDataDB("./data/db/bookdb/")}})
	spider.Run()
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
	spider    *Spider
	bookStore Store
}

func (p *BookPipeline) Process(handleResult *Result, ctx context.Context) error {
	bookInf, ok := handleResult.TargetItems["book"]
	if ok {
		book := bookInf.(Book)
		bookstr, err := BookStringify(book)
		if err != nil {
			return err
		}
		p.bookStore.Add(fmt.Sprintf("book-%s", book.Id), bookstr)
		return nil
	}
	chapterInf, ok := handleResult.TargetItems["chapter"]
	if ok {
		chapter := chapterInf.(Chapter)
		chapterstr, err := ChapterStringify(chapter)
		if err != nil {
			return err
		}
		p.bookStore.Add(fmt.Sprintf("book-chapter-%s-%s-%s", chapter.BookId, chapter.Title, chapter.PageIndex), chapterstr)
	}
	return nil
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

func ParseBook(str string) (*Book, error) {
	book := &Book{}
	err := json.Unmarshal([]byte(str), book)
	return book, err
}
