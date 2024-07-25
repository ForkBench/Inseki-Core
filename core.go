package main

import "inseki-core/tools"

func main() {
	structures := tools.ImportStructure("./structures")

	for _, structure := range structures {
		files := structure.NodeToString(false)
		for _, file := range files {
			println(file)
		}
	}
}
