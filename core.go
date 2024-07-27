package main

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"inseki-core/tools"
)

//go:embed config.json
var configJson string

func process() {
	// ----------------------------- Read the configuration -----------------------------
	config := tools.ReadEmbedConfigFile(configJson)

	// Check if the folder exists
	tools.CheckIfConfigFolderExists(config)

	insekiignore := tools.ReadInsekiIgnore(config)

	// ----------------------------- Read the structures -----------------------------
	numberStructuresAnalysed := 0

	structures := tools.ImportStructure(config, insekiignore, &numberStructuresAnalysed)

	if len(structures) == 0 {
		println("No structures found")
		os.Exit(0)
	}

	patterns := tools.ExtractNames(structures, false)
	associations := tools.StringNodeToAssociation(patterns)

	stack := &tools.Stack{}

	fmt.Printf("Number of structures analysed: %d\n", numberStructuresAnalysed)

	// ----------------------------- Explore the folder -----------------------------
	numberFilesAnalysed := 0

	tools.ExploreFolder("~/Documents/", insekiignore, tools.FilterWithPatternMap(&associations, stack), &numberFilesAnalysed)

	fmt.Printf("Number of files analysed: %d\n", numberFilesAnalysed)

	ch := make(chan tools.Target)
	var wg sync.WaitGroup

	for !stack.IsEmpty() {
		value := stack.Pop()

		wg.Add(1)
		go func(value tools.Target, ch chan tools.Target) {
			defer wg.Done()

			// TODO: Do something with the file

			// Add the path to the stack
			ch <- value
		}(value, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// for value := range ch {
	// 	println(value.String())
	// }
}

func main() {
	process()
}
