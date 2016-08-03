package catego

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Workiva/go-datastructures/bitarray"
)

const (
	DefaultRootNodeID ID = 0
)

// TreeOptions me
type TreeOptions struct {
	// RootNodeID is the id used for the root
	RootNodeID ID
	// NoIDSpecialID is the ID meaning : NoID declared for it
	NoIDSpecialID *ID
}

// Tree is the structure that will allow to search nodes in it
type Tree struct {
	sync.RWMutex
	registry IDToCat
	rootNode *Node
	maxID    ID
	option   *TreeOptions
}

// NewTree creates a tree with all options at default
// See TreeOption object for all options available
func NewTree(loader NodeSource) (*Tree, error) {
	return NewTreeWithOptions(loader, nil)
}

// NewTreeWithOptions creates a tree using a NodeSource. It will loop on the source until
// Next() returns false and add every node to the tree
func NewTreeWithOptions(loader NodeSource, opts *TreeOptions) (*Tree, error) {

	// changing the root node ID
	rootID := DefaultRootNodeID
	if opts != nil {
		rootID = opts.RootNodeID
	}

	t := &Tree{
		rootNode: &Node{
			ID: rootID,
		},
		registry: make(IDToCat),
		option:   opts,
	}

	t.registry[rootID] = t.rootNode

	var c ID
	var p ID
	var err error

	for loader.Next() {
		c, p, err = loader.Get()
		if err != nil {
			return nil, err
		}
		t.add(c, p)
	}
	return t, nil
}

// Add adds a new node to the tree. If the node has no parent, use 0. It will be attached to the node root
func (t *Tree) Add(current ID, parent ID) {

	t.Lock()
	defer t.Unlock()

	t.Add(current, parent)
}

func (t *Tree) add(current ID, parent ID) {

	if t.option != nil && t.option.NoIDSpecialID != nil && current == *t.option.NoIDSpecialID {
		// this ID is to ignore, no need to store it
		return
	}

	var currentPresent bool
	var parentPresent bool
	_, parentPresent = t.registry[parent]
	_, currentPresent = t.registry[current]

	if !currentPresent {
		t.registry[current] = &Node{
			ID: current,
		}
	}

	if !parentPresent {
		t.registry[parent] = &Node{
			ID:     parent,
			Parent: t.rootNode,
		}
	}

	if current > t.maxID {
		t.maxID = current
	}

	if parent > t.maxID {
		t.maxID = parent
	}

	t.registry[current].Parent = t.registry[parent]

	t.registry[parent].Children = append(t.registry[parent].Children, t.registry[current])

}

// Get returns the wanted node
// complexity is O(1)
func (t *Tree) Get(id ID) (*Node, error) {
	t.RLock()
	defer t.RUnlock()
	return t.get(id)
}

func (t *Tree) get(id ID) (*Node, error) {
	var ok bool

	if _, ok = t.registry[id]; ok {
		return t.registry[id], nil
	}
	return nil, fmt.Errorf("get: id not found %d", id)

}

// GetAncestors returns all the parent to the root down
// it walks the tree to the top. Complexity is O(n) where n is the distance to the top.
func (t *Tree) GetAncestors(id ID) ([]ID, error) {

	t.RLock()
	defer t.RUnlock()
	p, err := t.get(id)
	if err != nil {
		return nil, err
	}

	var parents []ID
	var currentParent *Node

	currentParent = p.Parent

	for {
		parents = append(parents, currentParent.ID)
		if currentParent.ID == t.rootNode.ID {
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
	return t.getChildren(id, nil)
}

// Exclude returns all node except the provided ones. It is an "heavy"
// operation.
func (t *Tree) Exclude(id []ID) ([]ID, error) {
	t.RLock()
	defer t.RUnlock()
	m := make(map[ID]bool, len(id))
	for i := range id {
		if id[i] == t.rootNode.ID {
			return nil, errors.New("root node cant be excluded, its all the tree")
		}
		m[id[i]] = true
	}
	return t.getChildren(t.rootNode.ID, m)
}

func (t *Tree) getChildren(id ID, exclude map[ID]bool) ([]ID, error) {
	var current *Node
	var err error
	current, err = t.get(id)
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
			if exclude != nil {
				// Do some exclusion
				var ok bool
				if _, ok = exclude[n.Children[i].ID]; ok {
					continue
				}
			}
			allChild = append(allChild, n.Children[i].ID)
			getChilds(n.Children[i])
		}
	}

	getChilds(current)

	return allChild, nil
}

// GetSiblings should returns all the node that are at the same level than the current node
// TODO: implement this
func (t *Tree) GetSiblings(id ID) ([]ID, error) {
	return nil, errors.New("to implement")
}

// GetBlackLister returns a Blacklister object, given the blacklisted and whitelisted nodes
// Blacklist : node and all children of it are banned
// Whitelist : all node but those ones and children are banned
func (t *Tree) GetBlackLister(blacklist []ID, whitelist []ID) (*Blacklister, error) {

	var blacklistedCategory []ID
	var err error

	if len(whitelist) > 0 {
		blacklistedCategory, err = t.Exclude(whitelist)
		if err != nil {
			return nil, err
		}
	}

	for i := range blacklist {
		// NoIDSpecialID means that this special ID means
		blacklistedCategory = append(blacklistedCategory, blacklist[i])
		if t.option != nil &&
			t.option.NoIDSpecialID != nil &&
			blacklist[i] == *t.option.NoIDSpecialID {
			continue
		}

		var descendant []ID
		descendant, err = t.GetDescendants(blacklist[i])
		if err != nil {
			return nil, err
		}
		blacklistedCategory = append(blacklistedCategory, descendant...)
	}

	b := bitarray.NewBitArray(uint64(t.maxID))

	for i := range blacklistedCategory {
		err = b.SetBit(uint64(blacklistedCategory[i]))
		if err != nil {
			return nil, err
		}
	}

	return &Blacklister{
		store: b,
	}, nil

}
