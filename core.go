package main

import (
	_ "embed"
	"os"
	"sync"

	"inseki-core/tools"
)

//go:embed config.json
var configJson string

func process() {
	// Read the config file
	config := tools.ReadEmbedConfigFile(configJson)

	// Check if the folder exists
	tools.CheckIfConfigFolderExists(config)

	insekiignore := tools.ReadInsekiIgnore(config)

	structures := tools.ImportStructure(config, insekiignore)

	if len(structures) == 0 {
		println("No structures found")
		os.Exit(0)
	}

	patterns := tools.ExtractNames(structures, false)
	associations := tools.StringNodeToAssociation(patterns)

	stack := &tools.Stack{}

	tools.ExploreFolder("~/Documents/", insekiignore, tools.FilterWithPatternMap(&associations, stack))

	ch := make(chan tools.Target)
	var wg sync.WaitGroup

	for !stack.IsEmpty() {
		value := stack.Pop()

		wg.Add(1)
		go func(value tools.Target, ch chan tools.Target) {
			defer wg.Done()

			target := tools.Target{
				Filepath:    value.Filepath,
				Association: value.Association,
			}

			// TODO: Do something with the file

			// Add the path to the stack
			ch <- target
		}(value, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func main() {
	process()
}
