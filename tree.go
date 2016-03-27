package catego

import (
	"sync"

	"github.com/juju/errgo/errors"
)

// Tree is the structure that will allow to search nodes in it
type Tree struct {
	sync.RWMutex
	registry IdToCat
	rootNode *Node
}

// NewTree creates a tree using a NodeSource. It will loop on the source until
// Next() returns false and add every node to the tree
func NewTree(loader NodeSource) (*Tree, error) {

	t := &Tree{
		rootNode: &Node{
			Id: 0,
		},
		registry: make(IdToCat),
	}

	t.registry[0] = t.rootNode

	var c ID
	var p ID
	var err error

	for loader.Next() {
		c, p, err = loader.Get()
		if err != nil {
			return nil, err
		}
		t.Add(c, p)
	}
	return t, nil
}

// Add adds a new node to the tree. If the node has no parent, use 0. It will be attached to the node root
func (t *Tree) Add(current ID, parent ID) {

	// TODO: move locker in public method and creates add() that will adds without lock
	t.Lock()
	defer func() { t.Unlock() }()
	var currentPresent bool
	var parentPresent bool
	_, parentPresent = t.registry[parent]
	_, currentPresent = t.registry[current]

	if !currentPresent {
		t.registry[current] = &Node{
			Id: current,
		}
	}

	if !parentPresent {
		t.registry[parent] = &Node{
			Id: parent,
		}
	}

	t.registry[current].Parent = t.registry[parent]

	t.registry[parent].Children = append(t.registry[parent].Children, t.registry[current])

}

// Get returns the wanted node
// complexity is O(1)
func (t *Tree) Get(id ID) (*Node, error) {

	// TODO: create t.get() to get without lock
	t.RLock()
	defer t.RUnlock()

	var ok bool

	if _, ok = t.registry[id]; ok {
		return t.registry[id], nil
	}
	return nil, errors.New("not found")
}

// GetAncestors returns all the parent to the root down
// it walks the tree to the top. Complexity is O(n) where n is the distance to the top.
func (t *Tree) GetAncestors(id ID) ([]ID, error) {

	t.RLock()
	defer t.RUnlock()
	p, err := t.Get(id)
	if err != nil {
		return nil, err
	}

	var parents []ID
	var currentParent *Node

	currentParent = p.Parent

	for {
		parents = append(parents, currentParent.Id)
		if currentParent.Id == 0 {
			break
		}
		currentParent = currentParent.Parent
	}

	return parents, nil
}

// GetDescendants returns all the children of the ID pass in param
// If the children have children all the hierarchy will be returned.
// Complexity is O(n) where n is the number of children
func (t *Tree) GetDescendants(id ID) ([]ID, error) {

	t.RLock()
	defer t.RUnlock()
	var current *Node
	var err error
	current, err = t.Get(id)
	if err != nil {
		return nil, err
	}

	var allChild []ID
	var getChilds func(n *Node)
	getChilds = func(n *Node) {
		if len(n.Children) == 0 {
			return
		}
		for i := range n.Children {
			allChild = append(allChild, n.Children[i].Id)
			getChilds(n.Children[i])
		}
	}

	getChilds(current)

	return allChild, nil
}

func (t *Tree) GetSiblings(id ID) ([]ID, error) {
	return nil, errors.New("to implement")
}
