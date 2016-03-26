package catego

import (
	"sync"

	"github.com/juju/errgo/errors"
)

type Tree struct {
	sync.RWMutex
	registry IdToCat
	rootNode *Node
}

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

func (t *Tree) Add(current ID, parent ID) {

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

func (t *Tree) Get(id ID) (*Node, error) {

	t.RLock()
	defer t.RUnlock()

	var ok bool
	var node *Node

	if node, ok = t.registry[id]; ok {
		return node, nil
	}
	return nil, errors.New("not found")
}

func (t *Tree) GetParents(id ID) ([]ID, error) {

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

func (t *Tree) GetChildren(id ID) ([]ID, error) {
	return nil, errors.New("to implement")
}

func (t *Tree) GetSiblings(id ID) ([]ID, error) {
	return nil, errors.New("to implement")
}
