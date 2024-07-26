package main

import "inseki-core/tools"

func main() {
	insekiignore := tools.ReadInsekiIgnore()

	structures := tools.ImportStructure("./structures", insekiignore)

	patterns := tools.ExtractNames(structures, false)

	tools.ExploreFolder("~/Documents/", insekiignore, tools.FilterWithPatternMap(patterns, func(filepath string, nodes []tools.Node) {
		println(filepath)
	}))
}
