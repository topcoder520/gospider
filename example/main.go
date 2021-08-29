package main

import "github.com/topcoder520/gospider"

func main() {
	spider := gospider.NewSpider("https://www.duquanben.com/")
	spider.SetGoroutines(2)
	spider.Run()
}
