package inseki

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Node Structure to represent a file system node
type Node struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Optional    bool   `json:"optional,omitempty"`
	Children    []Node `json:"children,omitempty"`
	HashValue   uint64 `json:"hash,omitempty"`
}

type Structure struct {
	Root Node `json:"root"`
	Hash uint64
	Name string
}

/*
String method to represent a Node as a string
*/
func (n Node) String(depth ...int) string {

	str := ""
	indent := strings.Repeat("\t", len(depth))

	for _, child := range n.Children {
		str += child.String(append(depth, 1)...)
	}

	if n.IsDirectory {
		return fmt.Sprintf("%sDirectory: %s (%s)\n%s", indent, n.Name, strconv.FormatBool(n.Optional), str)
	} else {
		return fmt.Sprintf("%sFile: %s (%s)\n", indent, n.Name, strconv.FormatBool(n.Optional))
	}
}

/*
String method to represent a Structure as a string
*/
func (s Structure) String() string {
	return s.Root.String()
}

/*
JSONToStructure method to read a JSON file and return a Structure
*/
func JSONToStructure(jsonPath string) (error, Structure) {
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return err, Structure{}
	}

	var rootNode Node
	err = json.Unmarshal(jsonData, &rootNode)
	if err != nil {
		return err, Structure{}
	}

	rootNode.HashValue = rootNode.Hash()

	structure := rootNode.NodeToStructure()
	structure.Name = filepath.Base(jsonPath)

	return nil, structure
}

/*
ImportStructure method to import all structures from a folder
*/
func ImportStructure(config Config, insekiIgnore []string, numberFilesAnalysed *int) (error, map[uint64]Structure) {
	nodes := make(map[uint64]Structure)

	path := TranslateDir(config.StructurePath)

	// Read all .json
	err := ExploreFolder(path, insekiIgnore, func(path string, info os.FileInfo) error {
		if strings.HasSuffix(path, ".json") {
			err, structure := JSONToStructure(path)
			if err != nil {
				return err
			}

			// Check if the hash is not in the map, add it
			if _, ok := nodes[structure.Hash]; !ok {
				nodes[structure.Hash] = structure
			} else {
				// If the hash is already in the map, check if the node is equal
				// If it is equal, then it is a duplicate
				if nodes[structure.Hash].Equal(structure, false) {
					return errors.New(fmt.Sprintf("Duplicate: %s\n", path))
				} else {
					// If it is not equal, then it is a conflict
					return errors.New(fmt.Sprintf("Conflict: %s\n", path))
				}
			}
		}
		return nil
	}, numberFilesAnalysed)

	if err != nil {
		return err, nil
	}

	return nil, nodes
}

/*
ExportStructure method to export a Node to a JSON file
*/
func ExportStructure(structure Structure, path string) error {
	jsonData, err := json.MarshalIndent(structure.Root, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path,
		jsonData,
		0644)
	if err != nil {
		return err
	}

	return nil
}

func SortNodes(structures *[]Structure) {
	sort.Slice(*structures, func(i, j int) bool {
		return (*structures)[i].Contains((*structures)[j])
	})
}

/*
ExtractNames
Extract all the names from a list of nodes
*/
func ExtractNames(nodes map[uint64]Structure, extractOptional bool) map[string][]Structure {
	names := make(map[string][]Structure)
	for _, structure := range nodes {
		structure.ExtractNames(extractOptional, names)
	}

	// For each name, sort the nodes
	for _, nodes := range names {
		// Sort the nodes
		SortNodes(&nodes)
	}

	return names
}

// ----------------------------- Node -----------------------------

/*
NodeToStructure method to convert a Node to a list of files, canBeOptional is used to include optional files
*/
func (n Node) NodeToStructure() Structure {
	return Structure{
		Root: n,
		Hash: n.Hash(),
		Name: "Undefined",
	}
}

/*
NodeToString method to convert a Node to a list of files, canBeOptional is used to include optional files
*/
func (n Node) NodeToString(canBeOptional bool) []string {
	var files []string
	for _, child := range n.Children {
		if child.IsDirectory {
			for _, file := range child.NodeToString(canBeOptional) {
				if !child.Optional || canBeOptional {
					files = append(files, fmt.Sprintf("%s/%s", n.Name, file))
				}
			}
		} else {
			if !child.Optional || canBeOptional {
				files = append(files, fmt.Sprintf("%s/%s", n.Name, child.Name))
			}
		}
	}

	return files
}

/*
Equal
See if a node is equal to another node (using hash) :
*/
func (n Node) Equal(other Node, canBeOptional bool) bool {
	// TODO: Check canBeOptional
	return n.Hash() == other.Hash() && n.Contains(other) && other.Contains(n)
}

/*
Contains
See if a node contains another node :

# If all the children of A are in B (where children of B can be optional), then A is in B

A.Contains(B) -> A is in B
*/
func (n Node) Contains(other Node) bool {

	// If the name is different, return false
	if n.Name != other.Name {
		return false
	}

	// If the node is a directory
	if n.IsDirectory {
		// For each child
		for _, child := range n.Children {
			// Check if the child is in the other node
			found := false
			for _, otherChild := range other.Children {
				if child.Contains(otherChild) {
					found = true
					break
				}
			}

			// If the child is not found and it is not optional, return false
			if !found && !child.Optional {
				return false
			}
		}
	}

	return true

}

// Matches : Check if a Structure matches a file with a specific depth
func (n Node) Matches(root string) bool {
	// Has to match from the root

	// If the current node is a file
	if !n.IsDirectory {
		// Check if the file exists with os.Glob
		files, _ := filepath.Glob(filepath.Join(root, n.Name))
		return len(files) > 0
	}

	// If the current node is a directory
	// Check if the root is the same as the name
	if matched, _ := filepath.Match(n.Name, filepath.Base(root)); matched {
		// Check if the children match (all non optional children need to be present)
		hasAllChildren := true

		for _, child := range n.Children {
			// If the child is optional, skip
			if child.Optional {
				continue
			}

			if !child.IsDirectory {
				// Check if the child is in the root
				if !child.Matches(root) {
					hasAllChildren = false
					break
				}
			} else {
				// Check if the child is in the root
				if !child.Matches(filepath.Join(root, child.Name)) {
					hasAllChildren = false
					break
				}
			}
		}

		return hasAllChildren
	}

	return false

}

func (n Node) GetDepths(filename string, depths *[]uint8, depth int) {

	base := filepath.Base(filename)

	for _, child := range n.Children {

		// Check if the node name is the same as the filename (matches, because it could be *.c)
		if matched, _ := filepath.Match(child.Name, base); matched {
			*depths = append(*depths, uint8(depth))
		}

		if child.IsDirectory {
			child.GetDepths(filename, depths, depth+1)
		}
	}
}

/*
Hash a node using merkle tree
*/
func (n Node) Hash(depth ...int) uint64 {
	if n.IsDirectory {
		var hash uint64
		for _, child := range n.Children {
			hash += child.Hash(append(depth, 1)...)
		}
		return hash
	} else {
		if len(n.Name) >= 2 {
			return uint64(n.Name[1]) + uint64(n.Name[len(n.Name)-1])<<len(depth)
		} else {
			return uint64(len(depth))
		}
	}
}

// ----------------------------- Structure -----------------------------

/*
ExtractNames
For a node, if its root isn't "*", then add it to the map and return
*/
func (s Structure) ExtractNames(extractOptional bool, names map[string][]Structure) {
	addToNames := func(name string, s Structure) {
		if _, ok := names[name]; !ok {
			names[name] = []Structure{s}
		} else {
			names[name] = append(names[name], s)
		}
	}

	if s.Root.Name != "*" {
		addToNames(s.Root.Name, s)
	}

	for _, child := range s.Root.Children {
		addToNames(child.Name, s)
	}
}

/*
StructureToString method to convert a Structure to a list of files, canBeOptional is used to include optional files
*/
func (s Structure) StructureToString(canBeOptional bool) []string {
	return s.Root.NodeToString(canBeOptional)
}

/*
StringToStructure method to convert a list of files to a Structure
*/
func StringToStructure(files []string) Structure {
	rootNode := Node{
		Name:        ".",
		IsDirectory: true,
	}

	for _, file := range files {
		path := strings.Split(file, "/")
		currentNode := &rootNode

		for i, name := range path {
			if i == len(path)-1 {
				currentNode.Children = append(currentNode.Children, Node{
					Name:        name,
					IsDirectory: false,
				})
			} else {
				found := false
				for _, child := range currentNode.Children {
					if child.Name == name {
						currentNode = &child
						found = true
						break
					}
				}

				if !found {
					newNode := Node{
						Name:        name,
						IsDirectory: true,
					}
					currentNode.Children = append(currentNode.Children, newNode)
					currentNode = &newNode
				}
			}
		}
	}

	return rootNode.NodeToStructure()
}

/*
StringNodeToAssociation
String-Node map to Association
*/
func StringNodeToAssociation(stringNode map[string][]Structure) []Association {
	var associations []Association
	for pattern, nodes := range stringNode {
		associations = append(associations, Association{
			Pattern:    pattern,
			Structures: nodes,
		})
	}

	return associations
}

/*
Contains
See if a structure contains another structure using Contains
*/
func (s Structure) Contains(other Structure) bool {
	return s.Root.Contains(other.Root)
}

func (s Structure) GetDepths(filename string) []uint8 {
	depths := make([]uint8, 0)

	s.Root.GetDepths(filename, &depths, 1)

	return depths
}

// Matches : Check if a Structure matches a file
// Returns the root of the structure
func (s Structure) Matches(path string) (bool, string) {
	/*
		The idea is the following :

		We have a Structure, and a match (for example, ~/home/dev/main.c)

		We need to check if **all** folders and files that are non-optional in the structure are contained around this file.

		1. First, we need to determine where is the detected file/folder in the structure (which depth)

		Ex : {
			README
			src {
				main.c
			}
		}

		-> main.c is at depth 1 (we need one "../" to go back to the root)

		2. Second, we need to go to the possible root folder
		3. See if the structure matches from the root

		/!\ There could be multiple depth :

		{
			README
			src {
				main.c
			}
			main.c
		}

		-> main.c is at two different depths
	*/

	path = filepath.Clean(path)

	// Get the depth of the file
	depths := s.GetDepths(filepath.Base(path))

	for _, depth := range depths {
		root := GoUp(path, depth)

		if s.Root.Matches(root) {
			return true, root
		}
	}

	return false, ""
}

/*
Equal
See if a structure is equal to another structure (using hash) :
*/
func (s Structure) Equal(other Structure, canBeOptional bool) bool {
	return s.Hash == other.Hash && s.Contains(other) && other.Contains(s)
}

// ----------------------------- Useful -----------------------------

// GoUp : Go up in the path
func GoUp(path string, n uint8) string {
	for i := 0; i < int(n); i++ {
		path = filepath.Dir(path)
	}
	return path
}

func isDirectory(path string) (bool, error) {
	// Get the file info
	info, err := os.Stat(path)
	if err != nil {
		return false, err // If there is an error, return it
	}
	// Check if it's a directory
	return info.IsDir(), nil
}
