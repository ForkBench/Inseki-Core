package tools

import (
	"fmt"
	"log"
	"sync"
)

// Process : Run a complete disk analysis
func Process(path string, config Config, insekiIgnore []string) (error, []Response) {

	// ----------------------------- Read the structures -----------------------------
	numberStructuresAnalysed := 0

	err, structures := ImportStructure(config, insekiIgnore, &numberStructuresAnalysed)
	if err != nil {
		return err, nil
	}

	if len(structures) == 0 {
		return fmt.Errorf("no structure found"), nil
	}

	patterns := ExtractNames(structures, false)
	associations := StringNodeToAssociation(patterns)

	stack := &Stack{}

	log.Printf("Number of structures analysed: %d\n", numberStructuresAnalysed)

	// ----------------------------- Explore the folder -----------------------------
	numberFilesAnalysed := 0

	err = ExploreFolder(path,
		insekiIgnore,
		FilterWithPatternMap(&associations, stack),
		&numberFilesAnalysed)
	if err != nil {
		return err, nil
	}

	fmt.Printf("Number of files analysed: %d\n", numberFilesAnalysed)

	ch := make(chan Response)
	var wg sync.WaitGroup

	for !stack.IsEmpty() {
		value := stack.Pop()

		wg.Add(1)
		go func(value Target, ch chan Response) {
			defer wg.Done()

			// For each structure, check if the file is a match
			for _, structure := range value.Association.Structures {
				if structure.Matches(value.Filepath) {
					ch <- Response{
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

	// ----------------------------- Process the results -----------------------------

	results := make([]Response, 0)

	for value := range ch {
		results = append(results, value)
	}

	return nil, results
}
