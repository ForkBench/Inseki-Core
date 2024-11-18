package main

import (
	_ "embed"
	"fmt"
	"time"
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

	s := Structure{
		Children: []Node{
			{
				Name:  "Test",
				IsDir: true,
				Depth: 1,
			},
		},
		RootPath: "/",
	}

	start := time.Now()
	fmt.Printf("Hash : %08x\n", s.Hash())
	elapsed := time.Since(start)
	snapshot := time.Now()
	fmt.Printf("First hash : %d\n", elapsed)

	// The next one are quicker (around 6 times by memorizing the hash)
	fmt.Printf("Hash : %08x\n", s.Hash())
	fmt.Printf("Hash : %08x\n", s.Hash())
	elapsed = time.Since(snapshot)
	fmt.Printf("Last hash : %d\n", elapsed)

}
