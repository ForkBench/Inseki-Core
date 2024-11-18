package inseki

import (
	"crypto/sha256"
	"encoding/json"
	"os"
	"path/filepath"
)

type Structure struct {
	Root    *Node `json:"root"`
	HashVal []byte
	// Map a filename with all its possible depth
	Depths map[string][]uint8
	Name   string
}

type Node struct {
	Children []*Node `json:"children,omitempty"`
	IsDir    bool    `json:"isDirectory"`
	Pattern  string  `json:"pattern,omitempty"`
}

func ImportStructures(structuresPath string) []Structure {
	structures := make([]Structure, 0)

	// Read recursively for jsons, and import them
	err := filepath.Walk(structuresPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			err, structure := readJSON(path)
			if err != nil {
				panic(err)
			}

			structures = append(structures, structure)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return structures
}

func readJSON(jsonPath string) (error, Structure) {
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return err, Structure{}
	}

	var rootNode Node
	err = json.Unmarshal(jsonData, &rootNode)
	if err != nil {
		return err, Structure{}
	}

	structure := rootNode.NodeToStructure()
	structure.Name = filepath.Base(jsonPath)

	return nil, structure
}

func (n *Node) NodeToStructure() Structure {
	return Structure{
		Root:    n,
		HashVal: n.hash(),
		Depths:  nil,
	}
}

func (n *Node) getDepths(depths map[string][]uint8, depth uint8) {

	if depths[n.Pattern] == nil {
		depths[n.Pattern] = make([]uint8, 0)
	}

	depths[n.Pattern] = append(depths[n.Pattern], depth)

	for _, child := range n.Children {
		child.getDepths(depths, depth+1)
	}
}

func (s *Structure) getDepths() map[string][]uint8 {

	if s.Depths != nil {
		// Optimization
		return s.Depths
	}

	depths := make(map[string][]uint8)

	s.Root.getDepths(depths, 0)

	return depths
}

func (n *Node) hash() []byte {
	hashed := sha256.New224()

	// Add the name of the node.
	hashed.Write([]byte(n.Pattern))

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

	hash := s.Root.hash()

	s.HashVal = hash

	return hash
}

func structureEquals(s1 *Structure, s2 *Structure) bool {
	return bytesCmp(s1.Hash(), s2.Hash())
}

func bytesCmp(bytes1 []byte, bytes2 []byte) bool {
	if len(bytes1) != len(bytes2) {
		return false
	}

	// Compare each byte
	for i := 0; i < len(bytes1); i++ {
		if bytes1[i] != bytes2[i] {
			return false
		}
	}

	return true
}

func (n1 *Node) contains(n2 *Node) bool {

	// Check the name
	if matched, _ := filepath.Match(n1.Pattern, n2.Pattern); !matched {
		return false
	}

	// Check if it's a dir
	if n1.IsDir != n2.IsDir {
		return false
	}

	for _, childS2 := range n2.Children {
		hasAMatch := false

		for _, childS1 := range n1.Children {
			// TODO: Reduce as the time goes

			// Note : We check if the smaller is contained in the other one
			if childS2.contains(childS1) {
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

	return s1.Root.contains(s2.Root)
}

func (n *Node) correspondsTo(path string) bool {

	// Check if a file / folder exists (with glob for "*.c")
	res, err := filepath.Glob(path)
	if err != nil || len(res) == 0 {
		return false
	}

	if n.IsDir {
		for _, child := range n.Children {
			if !child.correspondsTo(filepath.Join(path, child.Pattern)) {
				return false
			}
		}
	}

	return true
}

func (s *Structure) correspondsTo(root string) bool {
	return s.Root.correspondsTo(root)
}

func (s *Structure) IsAValidStructure(path string) bool {

	// Get the depths
	depths := s.getDepths()

	for name, associatedDepths := range depths {
		if matched, _ := filepath.Match(name, filepath.Base(path)); !matched {
			continue
		}

		for _, depth := range associatedDepths {
			root := getRoot(path, depth)
			if s.correspondsTo(root) {
				return true
			}
		}
	}

	return false

}

func getRoot(path string, depth uint8) string {
	for i := uint8(0); i < depth; i++ {
		path = filepath.Dir(path)
	}

	return path
}
