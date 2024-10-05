package interface_test

import (
	"fmt"
	"testing"
)

type Integer int

func TestInterface(t *testing.T) {
	var i int = 10
	i2 := Integer(i)

	i3 := int(i2)
	fmt.Println(i, i2, i3)
}

type User struct {
	Name string
	Age  int
}

type List interface {
	Add(index int, value any)
	Append(value any)
	Delete(index int)
}
type LinkedList struct {
	head *node
	tail *node
	Len  int
}

func (l *LinkedList) Add(index int, value any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Append(value any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(index int) {
	//TODO implement me
	panic("implement me")
}

type node struct {
}

type ArrayList struct {
}

func (a *ArrayList) Add(index int, value any) {
	//TODO implement me
	panic("implement me")
}

func (a *ArrayList) Append(value any) {
	//TODO implement me
	panic("implement me")
}

func (a *ArrayList) Delete(index int) {
	//TODO implement me
	panic("implement me")
}
