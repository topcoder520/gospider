package gospider

import (
	"fmt"
	"testing"
)

func TestSCR(t *testing.T) {
	scr := &RequestScheduler{}
	scr.Push(NewRequest())
	scr.Push(NewRequest())
	fmt.Println(scr.Poll())
	fmt.Println(scr.PollN(2))
}
