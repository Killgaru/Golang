package main

import (
	"fmt"
	"reflect"
	"runtime"
)

type memoizeFunction func(int, ...int) interface{}

// TODO реализовать
var (
	fibonacci       memoizeFunction
	romanForDecimal memoizeFunction
	m               = map[string]interface{}{}
)

//TODO Write memoization function

func (f memoizeFunction) memoize(vi int) interface{} {
	key := fmt.Sprintf("%v_%v", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), vi)
	//fmt.Println("Entered to memoize" + " " + key)
	if v, ok := m[key]; ok {
		return v
	}
	m[key] = f(vi)
	return m[key]
}

// TODO обернуть функции fibonacci и roman в memoize
func init() {

	fibonacci = func(fb int, fbx ...int) interface{} {
		fb0 := 0
		fb1 := 1
		switch {
		case fb == fb0:
			return fb0
		case fb == fb1:
			return fb1
		}
		for i := 2; i <= fb; i++ {
			fb1 += fb0
			fb0 = fb1 - fb0
		}
		return fb1

	}
	romanForDecimal = func(rd int, rdx ...int) interface{} {
		const (
			d1000 = 1000
			r1000 = "M"
			d900  = 900
			r900  = "CM"
			d500  = 500
			r500  = "D"
			d400  = 400
			r400  = "CD"
			d100  = 100
			r100  = "C"
			d90   = 90
			r90   = "XC"
			d50   = 50
			r50   = "L"
			d40   = 40
			r40   = "XL"
			d10   = 10
			r10   = "X"
			d9    = 9
			r9    = "IX"
			d5    = 5
			r5    = "V"
			d4    = 4
			r4    = "IV"
			d1    = 1
			r1    = "I"
		)
		var output string
		if rd >= 4000 {
			fmt.Println("I can't convert it: ", rd)
			return output
		}
		for i := 0; rd > 0; i++ {
			switch {
			case rd-d1000 >= 0:
				rd -= d1000
				output += r1000
			case rd-d900 >= 0:
				rd -= d900
				output += r900
			case rd-d500 >= 0:
				rd -= d500
				output += r500
			case rd-d400 >= 0:
				rd -= d400
				output += r400
			case rd-d100 >= 0:
				rd -= d100
				output += r100
			case rd-d90 >= 0:
				rd -= d90
				output += r90
			case rd-d50 >= 0:
				rd -= d50
				output += r50
			case rd-d40 >= 0:
				rd -= d40
				output += r40
			case rd-d10 >= 0:
				rd -= d10
				output += r10
			case rd-d9 >= 0:
				rd -= d9
				output += r9
			case rd-d5 >= 0:
				rd -= d5
				output += r5
			case rd-d4 >= 0:
				rd -= d4
				output += r4
			case rd-d1 >= 0:
				rd -= d1
				output += r1
			default:
				fmt.Println("I can't convert it: ", rd)
			}
		}
		return output

	}
}

func main() {
	memFibonacci := fibonacci.memoize
	memRomanForDecimal := romanForDecimal.memoize
	fmt.Println("Fibonacci(45) =", memFibonacci(45).(int))
	for _, x := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 16, 17, 18, 19, 20, 25, 30, 40, 50, 60, 69, 70, 80,
		90, 99, 100, 200, 300, 400, 500, 600, 666, 700, 800, 900,
		1000, 1009, 1444, 1666, 1945, 1997, 1999, 2000, 2008, 2010,
		2012, 2500, 3000, 3999} {
		fmt.Printf("%4d = %s\n", x, memRomanForDecimal(x).(string))
	}
}
