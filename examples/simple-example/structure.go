package main

import (
	"crypto/sha256"
)

type Structure struct {
	RootPath string
	Children []Node
	HashVal  []byte
}

type Node struct {
	Children []Node
	IsDir    bool
	Depth    uint8
	Name     string
}

// ------------------------ Common ------------------------

func (n *Node) hash() []byte {
	hashed := sha256.New224()
	hashed.Write([]byte(n.Name))

	for _, child := range n.Children {
		hashed.Sum(child.hash())
	}

	return hashed.Sum([]byte(n.Name))
}

func (s *Structure) Hash() []byte {

	if s.HashVal != nil {
		return s.HashVal
	}

	hashed := sha256.New224()
	hashed.Write([]byte(s.RootPath))

	for _, child := range s.Children {
		hashed.Sum(child.hash())
	}

	return hashed.Sum([]byte(s.RootPath))
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
