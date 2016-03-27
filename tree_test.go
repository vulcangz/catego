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

func TestGetAncestors(t *testing.T) {

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
		p, err = tree.GetAncestors(test[i].in)
		if !reflect.DeepEqual(p, test[i].out) {
			t.Fatalf("expected %v received %v", test[i].out, p)
		}
		if err != nil && err.Error() != test[i].err.Error() {
			t.Fatalf("expected error value %q received %q", test[i].err, err)
		}

		t.Logf("given %v received %v and error value %v", test[i].in, p, err)

	}
}

func TestGetDescendants(t *testing.T) {

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
			in:  7,
			out: []ID{8},
			err: nil,
		},
		{
			in:  0,
			out: []ID{1, 2, 4, 5, 6, 7, 8, 3},
			err: nil,
		},
		{
			in:  11,
			out: nil,
			err: errors.New("not found"),
		},
		{
			in:  2,
			out: []ID{4, 5, 6, 7, 8},
			err: nil,
		},
	}
	for i := range test {
		var p []ID
		p, err = tree.GetDescendants(test[i].in)
		if !reflect.DeepEqual(p, test[i].out) {
			t.Fatalf("expected %v received %v", test[i].out, p)
		}
		if err != nil && err.Error() != test[i].err.Error() {
			t.Fatalf("expected error value %q received %q", test[i].err, err)
		}

		t.Logf("given %v received %v and error value %v", test[i].in, p, err)

	}
}

func TestExclude(t *testing.T) {

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
		in  []ID
		out []ID
		err error
	}{
		{
			in:  []ID{7},
			out: []ID{1, 2, 4, 5, 6, 3},
			err: nil,
		},
		{
			in:  []ID{0},
			out: nil,
			err: errors.New("root node cant be excluded, its all the tree"),
		},
		{
			in:  []ID{11},
			out: []ID{1, 2, 4, 5, 6, 7, 8, 3},
			err: nil,
		},
		{
			in:  []ID{2},
			out: []ID{1, 3},
			err: nil,
		},
	}
	for i := range test {
		var p []ID
		p, err = tree.Exclude(test[i].in)
		if !reflect.DeepEqual(p, test[i].out) {
			t.Fatalf("given %v expected %v received %v", test[i].in, test[i].out, p)
		}
		if err != nil && err.Error() != test[i].err.Error() {
			t.Fatalf("given %v expected error value %q received %q", test[i].in, test[i].err, err)
		}

		t.Logf("given %v received %v and error value %v", test[i].in, p, err)

	}
}

func TestBlacklist(t *testing.T) {
	tree, err := NewTree(newLoaderTest([][2]ID{
		{1, 0},
		{2, 1},
		{3, 0},
		{4, 2},
		{5, 2},
		{6, 2},
		{7, 6},
		{8, 7},
		{10234, 8},
	}))
	if err != nil {
		t.Fatal(err)
	}

	test := []struct {
		blacklist          []ID
		whitelist          []ID
		testBlacklister    []ID
		testBlacklisterOut []bool
		err                error
	}{
		{
			blacklist:          []ID{7},
			whitelist:          nil,
			testBlacklister:    []ID{6, 7, 8},
			testBlacklisterOut: []bool{false, true, true},
			err:                nil,
		},
		{
			blacklist:          nil,
			whitelist:          []ID{7},
			testBlacklister:    []ID{1, 2, 3, 4, 5, 6, 7, 8},
			testBlacklisterOut: []bool{true, true, true, true, true, true, false, false},
			err:                nil,
		},
		{
			blacklist:          []ID{7},
			whitelist:          []ID{2},
			testBlacklister:    []ID{1, 2, 3, 4, 5, 6, 7, 8},
			testBlacklisterOut: []bool{true, false, true, false, false, false, true, true},
			err:                nil,
		},
	}

	for i := range test {

		var b *Blacklister
		b, err = tree.GetBlackLister(test[i].blacklist, test[i].whitelist)
		if err != nil && err.Error() != test[i].err.Error() {
			t.Fatalf("received error %v different than expected %v", err, test[i].err)
		}

		for j := range test[i].testBlacklister {
			if res := b.Is(test[i].testBlacklister[j]); res != test[i].testBlacklisterOut[j] {
				t.Fatalf(
					"given %v, expected %v and received %v - blacklist contains %v",
					test[i].testBlacklister[j], test[i].testBlacklisterOut[j], res, b.store.ToNums(),
				)
			}
		}

	}
}

func BenchmarkGetAncestors(b *testing.B) {
	var source [][2]ID

	size := 100
	for i := 0; i < size; i++ {
		source = append(source, [2]ID{ID(i + 1), ID(i)})
	}

	tree, err := NewTree(newLoaderTest(source))

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = tree.GetAncestors(ID(size))
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkGetDescendants(b *testing.B) {
	var source [][2]ID

	size := 1000
	for i := 0; i < size; i++ {
		source = append(source, [2]ID{ID(i + 1), ID(i)})
	}

	tree, err := NewTree(newLoaderTest(source))

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = tree.GetDescendants(0)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkBlacklist(b *testing.B) {
	var source [][2]ID

	size := 1000
	for i := 0; i < size; i++ {
		source = append(source, [2]ID{ID(i + 1), ID(i)})
	}

	tree, err := NewTree(newLoaderTest(source))

	if err != nil {
		b.Fatal(err)
	}

	var bl *Blacklister
	bl, err = tree.GetBlackLister([]ID{500}, nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		bl.Is(ID(i))
	}

}
