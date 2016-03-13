package catego

import "sync"

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
	defer func() { _ = t.Unlock() }()
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
