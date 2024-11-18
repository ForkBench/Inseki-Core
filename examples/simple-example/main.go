package main

import (
	_ "embed"
	inseki "github.com/ForkBench/Inseki-Core"
	"math/rand/v2"
)

func cmpInt(a int, b int) int8 {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

func main() {

	set := inseki.NewSet[int]()

	for i := 0; i < 10; i++ {
		// Random number
		set.Add(rand.Int()%100, cmpInt)
	}

	for !set.IsEmpty() {
		element, _ := set.Get()
		println(element)
	}

}
