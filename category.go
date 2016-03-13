package catego

type ID uint64

type Node struct {
	Id       ID
	Parent   *Node
	Children []*Node
}
