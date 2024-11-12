package main

import (
	"fmt"
	"log"
	"sync"
)

func analyze(path string, associations []Association, stack *Stack, insekiIgnore []string) (error, []Response) {
	// ----------------------------- Explore the folder -----------------------------
	numberFilesAnalysed := 0

	err := ExploreFolder(path,
		insekiIgnore,
		FilterWithPatternMap(&associations, stack),
		&numberFilesAnalysed)
	if err != nil {
		return err, nil
	}

	log.Printf("Number of files analysed: %d\n", numberFilesAnalysed)

	ch := make(chan Response)
	var wg sync.WaitGroup

	for !stack.IsEmpty() {
		value := stack.Pop()

		wg.Add(1)
		go func(value Target, ch chan Response) {
			defer wg.Done()

			// For each structure, check if the file is a match
			for _, structure := range value.Association.Structures {
				matched, root := structure.Matches(value.Filepath)
				if matched {
					ch <- Response{
						Filepath:  value.Filepath,
						Structure: structure,
						Root:      root,
					}
				}
			}
		}(value, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return nil, processResult(ch)
}

func processResult(ch chan Response) []Response {
	// map[root] = []Response
	// Helps to delete duplicates and intricated structures
	sorted := make(map[string][]Response)

	for response := range ch {
		// If the root is not in the map, create a new entry
		if _, ok := sorted[response.Root]; !ok {
			sorted[response.Root] = make([]Response, 0)
		}

		isDuplicate := false

		// If it's a duplicate (same structure and root) or intricated structure, skip
		for index, value := range sorted[response.Root] {

			// If the structure is the same, skip
			if value.Structure.Equal(response.Structure, true) {
				isDuplicate = true
				break
			}

			// If the structure current structure is contained in the stored structure, skip
			if value.Structure.Contains(response.Structure) {
				isDuplicate = true
				break
			}

			// If the stored structure is contained in the current structure, remove the stored structure
			if response.Structure.Contains(value.Structure) {
				sorted[response.Root] = append(sorted[response.Root][:index], sorted[response.Root][index+1:]...)
			}

		}

		if isDuplicate {
			continue
		}

		// Append the response to the root
		sorted[response.Root] = append(sorted[response.Root], response)
	}

	// ----------------------------- Return the results -----------------------------

	results := make([]Response, 0)

	for _, value := range sorted {
		for _, response := range value {
			results = append(results, response)
		}
	}

	return results
}

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

	// ----------------------------- Analyze the folder -----------------------------

	err, val := analyze(path, associations, stack, insekiIgnore)

	// ----------------------------- Process the results -----------------------------

	return nil, val
}
