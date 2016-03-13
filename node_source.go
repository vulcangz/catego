package catego

type NodeSource interface {
	// Will be used as a condition for a loop
	Next() bool
	// Return the current node id. If parent is 0 then it is a root node
	Get() (current ID, parent ID, err error)
}
