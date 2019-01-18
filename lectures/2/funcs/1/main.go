package main

import (
	"fmt"
	"math"
)

// TODO: Реализовать вычисление Квадратного корня
func Sqrt(x float64) float64 {
	eps := 1e-15
	x1 := 1.0
	for i := 0; Abs(x1-0.5*(x1+x/x1)) > eps; i++ {
		x1 = 0.5 * (x1 + x/x1)
	}
	return x1
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println(Sqrt(-1) - math.Sqrt(-1))
	fmt.Println(math.Sqrt(-2))
}
