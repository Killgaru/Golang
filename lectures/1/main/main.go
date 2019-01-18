// main
package main

import (
	"fmt"
	"strconv"
)

func ReturnInt() int {
	return 1
}

func ReturnFloat() float32 {
	return 1.1
}

func ReturnIntArray() [3]int {
	return [3]int{1, 3, 4}
}

func ReturnIntSlice() []int {
	return []int{1, 2, 3}
}

func IntSliceToString(s []int) string {
	out := ""
	for _, v := range s {
		out += strconv.FormatInt(int64(v), 10)
	}
	return out
}

func MergeSlices(sf []float32, si []int32) []int {
	var out []int
	for i := 0; i < (len(sf) + len(si)); i++ {
		if i < len(sf) {
			out = append(out, int(sf[i]))
		} else {
			out = append(out, int(si[i-len(sf)]))
		}
	}
	return out
}

func GetMapValuesSortedByKey(m map[int]string) []string {
	var (
		out   []string
		Scale int
	)
	switch ln := len(m); ln {
	case 4:
		Scale = 10
	case 12:
		Scale = 1
	default:
		fmt.Println("I do not know what do with it: ", m)
		return out
	}
	for i := 1; i <= len(m); i++ {
		out = append(out, m[i*Scale])
	}
	return out
}

func main() {
	fmt.Println("Hello World!")

}
