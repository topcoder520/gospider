package main

import "github.com/topcoder520/gospider"

func main() {
	spider := gospider.NewSpider("https://studygolang.com/pkgdoc")
	spider.SetGoroutines(2)
	spider.Run()
}
