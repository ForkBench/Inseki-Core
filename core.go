package main

import (
	_ "embed"
	"os"
	"sync"

	"inseki-core/tools"
)

//go:embed config.json
var configJson string

func main() {
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

	ch := make(chan string)
	var wg sync.WaitGroup

	for !stack.IsEmpty() {
		value := stack.Pop()

		wg.Add(1)
		go func(value tools.StackValue, ch chan string) {
			defer wg.Done()

			// TODO: Do something with the file

			// Add the path to the stack
			ch <- value.Filepath
		}(value, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for path := range ch {
		println(path)
	}
}
