package inseki

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const SET_WAIT_SIZE = 10

func SortByName(entry string, otherEntry string) int8 {
	return int8(strings.Compare(entry, otherEntry))
}

func ResearchAndAnalyse(root string, structures []Structure) {
	set := Set[string]{}

	var wg sync.WaitGroup // Wait for all goroutines to finish

	// Start the research
	SpreadResearch(root, set, &wg, structures)

	wg.Wait()

}

func SpreadResearch(root string, set Set[string], wg *sync.WaitGroup, structures []Structure) {

	root = translateDir(root)

	info, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic("bad path")
	}

	entries, _ := os.ReadDir(root)

	for _, entry := range entries {
		absolutePath := filepath.Join(root, entry.Name())

		set.Add(absolutePath, SortByName)

		if entry.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					if set.Size() > SET_WAIT_SIZE {
						break
					}
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					if ValidateStructures(structures, set) {
					}
				}()

				// Dispatch recursively
				SpreadResearch(absolutePath, set, wg, structures)
			}()
		}
	}
}

func ValidateStructures(structures []Structure, set Set[string]) bool {
	for _, structure := range structures {
		val, err := set.Get()
		if err != nil {
			panic(err)
		}

		if structure.IsAValidStructure(val) {
			return true
		}
	}

	return false
}

func translateDir(dir string) string {
	// Replace the ~ with the home directory
	if dir[:1] == "~" {
		dir = os.Getenv("HOME") + dir[1:]
	}

	return dir
}
