package control

import (
	"fmt"
	"testing"
)

func TestDeferV1(t *testing.T) {
	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Printf("%p %v\n", &i, i)
		}()
	}
}

func TestDeferV2(t *testing.T) {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			fmt.Printf("%p %v\n", &val, val)
		}(i)
	}
}

func TestDeferV3(t *testing.T) {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			fmt.Printf("%p %v\n", &j, j)
		}()
	}
}

func TestName(t *testing.T) {
	fmt.Print(1 + 1)
}
