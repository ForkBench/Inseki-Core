package main

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"github.com/ForkBench/Inseki-Core/tools"
)

//go:embed config.json
var configJson string

func process() {
	// ----------------------------- Read the configuration -----------------------------
	config := tools.ReadEmbedConfigFile(configJson)

	// Check if the folder exists
	tools.CheckIfConfigFolderExists(config)

	insekiIgnore := tools.ReadInsekiIgnore(config)

	// ----------------------------- Read the structures -----------------------------
	numberStructuresAnalysed := 0

	structures := tools.ImportStructure(config, insekiIgnore, &numberStructuresAnalysed)

	if len(structures) == 0 {
		println("No structures found")
		os.Exit(-1)
	}

	patterns := tools.ExtractNames(structures, false)
	associations := tools.StringNodeToAssociation(patterns)

	stack := &tools.Stack{}

	fmt.Printf("Number of structures analysed: %d\n", numberStructuresAnalysed)

	// ----------------------------- Explore the folder -----------------------------
	numberFilesAnalysed := 0

	err := tools.ExploreFolder("~/Documents/",
		insekiIgnore,
		tools.FilterWithPatternMap(&associations, stack),
		&numberFilesAnalysed)
	if err != nil {
		return
	}

	fmt.Printf("Number of files analysed: %d\n", numberFilesAnalysed)

	ch := make(chan tools.Response)
	var wg sync.WaitGroup

	for !stack.IsEmpty() {
		value := stack.Pop()

		wg.Add(1)
		go func(value tools.Target, ch chan tools.Response) {
			defer wg.Done()

			// For each structure, check if the file is a match
			for _, structure := range value.Association.Structures {
				if structure.Matches(value.Filepath) {
					ch <- tools.Response{
						Filepath:  value.Filepath,
						Structure: structure,
					}
				}
			}
		}(value, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for value := range ch {
		println(value.String())
	}
}

func main() {
	process()
}
