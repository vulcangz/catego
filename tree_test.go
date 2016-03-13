package catego

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type loaderTest struct {
	offset int
	data   [][2]ID
}

func newLoaderTest() *loaderTest {
	return &loaderTest{
		offset: -1,
		data: [][2]ID{
			{1, 0},
			{2, 0},
			{3, 0},
			{4, 2},
			{5, 2},
			{6, 1},
		},
	}
}

func (l *loaderTest) Next() bool {
	l.offset++
	if len(l.data)-1 == l.offset {
		return false
	}
	return true
}

func (l *loaderTest) Get() (current ID, parent ID, err error) {
	return l.data[l.offset][0], l.data[l.offset][1], nil
}

func TestTree(t *testing.T) {

	tree, err := NewTree(newLoaderTest())
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(tree)
}
