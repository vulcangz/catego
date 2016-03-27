package catego

// ID is the identifier of all node in the tree
type ID uint64

// Node describe a node into the Tree structure
// You can access to the parent or the direct children by getting it
type Node struct {
	ID       ID
	Parent   *Node
	Children []*Node
}
