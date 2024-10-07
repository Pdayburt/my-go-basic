package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {

	//仅仅声明没有初始化，读写都会报错
	/*var ch chan int

	var chV1 chan struct{}*/

	/*ints := make(chan int)
	c := make(chan int, 2)
	*/

	/*ch2 := make(chan int, 2)
	ch2 <- 123
	val := <-ch2
	fmt.Println(val)*/

}

type MyStruct struct {
	ch        chan struct{}
	closeOnce sync.Once
}

func (m *MyStruct) Close() error {
	m.closeOnce.Do(func() {
		close(m.ch)
	})
	return nil
}

func TestChannelBlock(t *testing.T) {

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	go func() {
		time.Sleep(time.Second)
		ch1 <- 1
	}()
	select {
	case val := <-ch2:
		fmt.Println("ch2", val)
	case val := <-ch1:
		fmt.Println("ch1", val)
	}

}
