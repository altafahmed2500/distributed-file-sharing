package merkle

import (
	"crypto/sha256"
	"distributed-file-sharing/services/chunker"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Node struct {
	Hash  string
	Left  *Node
	Right *Node
}

// Build Merkle Tree from chunk hashes
func BuildMerkleTree(chunkHashes []string) *Node {
	var nodes []*Node

	// Create leaf nodes
	for _, h := range chunkHashes {
		nodes = append(nodes, &Node{Hash: h})
	}

	// Build upwards
	for len(nodes) > 1 {
		var newLevel []*Node
		for i := 0; i < len(nodes); i += 2 {
			if i+1 == len(nodes) {
				// Odd number of nodes, promote directly
				newLevel = append(newLevel, nodes[i])
			} else {
				combined := nodes[i].Hash + nodes[i+1].Hash
				newHash := sha256.Sum256([]byte(combined))
				newNode := &Node{
					Hash:  hex.EncodeToString(newHash[:]),
					Left:  nodes[i],
					Right: nodes[i+1],
				}
				newLevel = append(newLevel, newNode)
			}
		}
		nodes = newLevel
	}

	return nodes[0]
}

// Load Meta.json and return chunk hashes
func LoadChunkHashes(metaPath string) ([]string, error) {
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	var meta chunker.FileMeta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	var hashes []string
	for _, chunk := range meta.Chunks {
		hashes = append(hashes, chunk.Hash)
	}
	return hashes, nil
}

// Print Merkle Tree
func PrintMerkleTree(node *Node, level int) {
	if node == nil {
		return
	}
	PrintMerkleTree(node.Right, level+1)
	fmt.Printf("%s%s\n", spaces(level*4), shortHash(node.Hash))
	PrintMerkleTree(node.Left, level+1)
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

func shortHash(hash string) string {
	if len(hash) > 8 {
		return hash[:8] + "..."
	}
	return hash
}
