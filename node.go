package radix

// Node represents a node in the Radix Tree. In most ways it acts as a Trie
// node with the difference that Label can be an arbitrary length depending on
// what other words are added to the tree. Prefix is stored so we don't need to
// determine a given node's prefix at query time when gathering suggestions for
// a query.
type Node struct {
	Label          string
	Prefix         string
	Children       []*Node
	IsWordBoundary bool
	IsRootNode     bool
}

// InitNode initializes an empty Node
func InitNode(label string) *Node {
	var children []*Node

	return &Node{
		Label:          label,
		Prefix:         "",
		Children:       children,
		IsWordBoundary: false,
		IsRootNode:     false,
	}
}

// IsWord checks if a Node marks the end of a word
func (n *Node) IsWord() bool {
	return n.IsWordBoundary
}

// IsRoot checks if a node is the root node of the Tree
func (n *Node) IsRoot() bool {
	return n.IsRootNode
}

// IsLeaf checks if a node has any children
func (n *Node) IsLeaf() bool {
	return len(n.Children) == 0
}
