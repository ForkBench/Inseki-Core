package main

import (
	_ "embed"
	"fmt"
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

	s1 := &Structure{
		Children: []*Node{
			{
				Name:  "*",
				IsDir: true,
				Depth: 1,
			},
		},
		RootPath: "/",
	}

	s2 := &Structure{
		Children: []*Node{
			{
				Name:  "Test",
				IsDir: true,
				Depth: 1,
			},
		},
		RootPath: "/",
	}

	fmt.Printf("%0x\n", s1.Hash())
	fmt.Printf("%0x\n", s2.Hash())
	println(s1.Contains(s2))
	println(s2.Contains(s1))
	println(s1.Contains(s1))
	println(s2.Contains(s2))
	println(structureEquals(s1, s2))
	println(structureEquals(s2, s1))

}
