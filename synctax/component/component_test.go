package component

import "testing"

type Inner struct {
}

func (i Inner) DoSomething() {

}

type Outter struct {
	Inner
}

func TestCom(t *testing.T) {
	var o Outter
	o.DoSomething()
	var op *Outter
	op.DoSomething()
}
