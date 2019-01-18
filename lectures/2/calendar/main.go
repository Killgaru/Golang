// main
package main

import (
	"fmt"
	"time"
)

type calendar struct {
	month   string
	quarter int
}

func NewCalendar(t time.Time) (c calendar) {
	a := t.Month()
	b := 1
	for i := b; a != time.Month(i); i++ {
		b = i + 1
	}
	if b%3 == 0 {
		c = calendar{a.String(), (b / 3)}
	} else {
		c = calendar{a.String(), (b/3 + 1)}
	}
	//fmt.Printf(" val: %v\n type: %T\n", c, c)
	return
}

func (c calendar) CurrentQuarter() int {
	return c.quarter
}

func main() {
	tp, _ := time.Parse("2006-01-02", "2015-11-15")
	fmt.Printf(" val: %v\n type: %T\n", tp, tp)
	c := NewCalendar(tp)
	fmt.Println(c.CurrentQuarter())

}
