package structures

import (
	"fmt"
	"path/filepath"
)

// Node Structure to represent a file system node
type Node struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Optional    bool   `json:"optional,omitempty"`
	Children    []Node `json:"children,omitempty"`
	HashValue   uint64 `json:"hash,omitempty"`
}

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
