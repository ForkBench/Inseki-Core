package main

import (
	"crypto/sha256"
	"path/filepath"
)

type Structure struct {
	RootPath string
	Children []*Node
	HashVal  []byte
}

type Node struct {
	Children []*Node
	IsDir    bool
	Depth    uint8
	Name     string
}

// ------------------------ Common ------------------------
func (n *Node) hash() []byte {
	hashed := sha256.New224()

	// Add the name of the node.
	hashed.Write([]byte(n.Name))

	// Add hashes of all children recursively.
	for _, child := range n.Children {
		childHash := child.hash()
		hashed.Write(childHash)
	}

	// Return the final hash.
	return hashed.Sum(nil)
}

func (s *Structure) Hash() []byte {
	// Return cached hash if already computed.
	if s.HashVal != nil {
		return s.HashVal
	}

	hashed := sha256.New224()

	// Add the root path of the structure.
	hashed.Write([]byte(s.RootPath))

	// Add hashes of all children nodes.
	for _, child := range s.Children {
		childHash := child.hash()
		hashed.Write(childHash)
	}

	// Cache the hash and return it.
	s.HashVal = hashed.Sum(nil)
	return s.HashVal
}

func structureEquals(s1 *Structure, s2 *Structure) bool {
	return bytesCmp(s1.Hash(), s2.Hash())
}

func bytesCmp(bytes1 []byte, bytes2 []byte) bool {
	if len(bytes1) != len(bytes2) {
		return false
	}

	for i := 0; i < len(bytes1); i++ {
		if bytes1[i] != bytes2[i] {
			return false
		}
	}

	return true
}

func (n1 *Node) Contains(n2 *Node) bool {

	// Check the name
	if matched, _ := filepath.Match(n1.Name, n2.Name); !matched {
		return false
	}

	// Check if it's a dir
	if n1.IsDir != n2.IsDir {
		return false
	}

	// TODO: Refactor
	for _, childS2 := range n2.Children {
		hasAMatch := false

		for _, childS1 := range n1.Children {
			// TODO: Reduce as the time goes

			// Note : We check if the smaller is contained in the other one
			if childS2.Contains(childS1) {
				hasAMatch = true
				break
			}
		}

		if !hasAMatch {
			return false
		}
	}

	return true
}

func (s1 *Structure) Contains(s2 *Structure) bool {

	// For efficiency
	if structureEquals(s1, s2) {
		return true
	}

	for _, childS2 := range s2.Children {
		hasAMatch := false

		for _, childS1 := range s1.Children {
			// TODO: Reduce as the time goes

			// Note : We check if the smaller is contained in the other one
			if childS2.Contains(childS1) {
				hasAMatch = true
				break
			}
		}

		if !hasAMatch {
			return false
		}
	}

	return true
}
