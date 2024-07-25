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

	return rootNode
}

/*
ImportStructure method to import all structures from a folder
*/
func ImportStructure(structuresRoot string) []Node {
	var nodes []Node

	// Read all .json
	err := ExploreFolder(structuresRoot, func(path string, info os.FileInfo) error {
		if strings.HasSuffix(path, ".json") {
			nodes = append(nodes, JSONToNode(path))
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