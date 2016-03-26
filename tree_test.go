package catego

import (
	"testing"

	"reflect"

	"github.com/juju/errgo/errors"
)

type loaderTest struct {
	offset int
	data   [][2]ID
}

func newLoaderTest(data [][2]ID) *loaderTest {
	return &loaderTest{
		offset: -1,
		data:   data,
	}
}

func (l *loaderTest) Next() bool {
	if len(l.data) == l.offset+1 {
		return false
	}
	l.offset++
	return true
}

func (l *loaderTest) Get() (current ID, parent ID, err error) {
	return l.data[l.offset][0], l.data[l.offset][1], nil
}

func TestGetParents(t *testing.T) {

	tree, err := NewTree(newLoaderTest([][2]ID{
		{1, 0},
		{2, 1},
		{3, 0},
		{4, 2},
		{5, 2},
		{6, 2},
		{7, 6},
		{8, 7},
	}))
	if err != nil {
		t.Fatal(err)
	}

	test := []struct {
		in  ID
		out []ID
		err error
	}{
		{
			in:  8,
			out: []ID{7, 6, 2, 1, 0},
			err: nil,
		},
		{
			in:  9,
			out: nil,
			err: errors.New("not found"),
		},
		{
			in:  3,
			out: []ID{0},
			err: nil,
		},
		{
			in:  5,
			out: []ID{2, 1, 0},
			err: nil,
		},
	}
	for i := range test {
		var p []ID
		p, err = tree.GetParents(test[i].in)
		if !reflect.DeepEqual(p, test[i].out) {
			t.Fatalf("expected %v received %v", test[i].out, p)
		}
		if err != nil && err.Error() != test[i].err.Error() {
			t.Fatalf("expected error value %q received %q", test[i].err, err)
		}

		t.Logf("given %v received %v and error value %v", test[i].in, p, err)

	}
}

func BenchmarkGetParents(b *testing.B) {
	var source [][2]ID

	for i := 0; i < 100; i++ {
		source = append(source, [2]ID{ID(i + 1), ID(i)})
	}

	tree, err := NewTree(newLoaderTest(source))

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = tree.GetParents(100)
		if err != nil {
			b.Fatal(err)
		}
	}

}
