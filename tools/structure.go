package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Structure to represent a file system node
type Node struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Optional    bool   `json:"optional,omitempty"`
	Children    []Node `json:"children,omitempty"`
	HashValue   uint64 `json:"hash,omitempty"`
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
JSONToNode method to read a JSON file and return a Node
*/
func JSONToNode(jsonPath string) Node {
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	var rootNode Node
	err = json.Unmarshal(jsonData, &rootNode)
	if err != nil {
		panic(err)
	}

	rootNode.HashValue = rootNode.Hash()

	return rootNode
}

/*
ImportStructure method to import all structures from a folder
*/
func ImportStructure(structuresRoot string) map[uint64]Node {
	nodes := make(map[uint64]Node)

	// Read all .json
	err := ExploreFolder(structuresRoot, func(path string, info os.FileInfo) error {
		if strings.HasSuffix(path, ".json") {
			node := JSONToNode(path)

			// Check if the hash is not in the map, add it
			if _, ok := nodes[node.Hash()]; !ok {
				nodes[node.Hash()] = node
			} else {
				// If the hash is already in the map, check if the node is equal
				// If it is equal, then it is a duplicate
				if nodes[node.Hash()].Equal(node, false) {
					panic(fmt.Sprintf("Duplicate: %s\n", path))
				} else {
					// If it is not equal, then it is a conflict
					panic(fmt.Sprintf("Conflict: %s\n", path))
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return nodes
}

/*
ExportStructure method to export a Node to a JSON file
*/
func ExportStructure(node Node, path string) {
	jsonData, err := json.MarshalIndent(node, "", "    ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path,
		jsonData,
		0644)
	if err != nil {
		panic(err)
	}
}

/*
See if a node is equal to another node (using hash) :
*/
func (n Node) Equal(other Node, canBeOptional bool) bool {
	return n.Hash() == other.Hash() && n.Contains(other) && other.Contains(n)
}

/*
See if a node contains another node :

# If all the children of A are in B (where children of B can be optional), then A is in B

A.Contains(B) -> A is in B
*/
func (n Node) Contains(other Node) bool {

	if n.Name != other.Name || n.IsDirectory != other.IsDirectory || n.Optional != other.Optional {
		return false
	}

	for _, child := range n.Children {
		found := false
		for _, otherChild := range other.Children {
			if child.Equal(otherChild, true) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

/*
Hash a node using merkle tree
*/
func (n Node) Hash() uint64 {
	if n.IsDirectory {
		var hash uint64
		for _, child := range n.Children {
			hash += child.Hash()
		}
		return hash
	} else {
		if len(n.Name) >= 2 {
			return uint64(n.Name[1]) + uint64(n.Name[len(n.Name)-1])<<len(n.Name)
		} else {
			return 0
		}
	}
}
