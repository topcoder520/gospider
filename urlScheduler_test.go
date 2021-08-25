package gospider

import (
	"fmt"
	"testing"
)

func TestSCR(t *testing.T) {
	scr := &UrlScheduler{}
	scr.Push("12345")
	scr.Push("sssss")
	fmt.Println(scr.Poll())
	fmt.Println(scr.PollN(2))
}
