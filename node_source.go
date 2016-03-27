package catego

// NodeSource is the interface used to put data into the tree
// The loader will call Next() and Get() to fetch the categories
// until Next() returns false
type NodeSource interface {
	// Will be used as a condition for a loop
	Next() bool
	// Return the current node id. If parent is 0 then it is a root node
	Get() (current ID, parent ID, err error)
}
