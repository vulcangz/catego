package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/mrsinham/catego.v1"
)

// ID is the identifier of all node in the tree
// type ID uint64

type category struct {
	ID     catego.ID
	Parent catego.ID
	Name   string
}

type categories []category

type loaderTest struct {
	offset int
	data   categories
}

func newLoaderTest() *loaderTest {
	return &loaderTest{
		offset: -1,
		data:   treeInput,
	}
}

func (l *loaderTest) Next() bool {
	if len(l.data) == l.offset+1 {
		return false
	}
	l.offset++
	return true
}

func (l *loaderTest) Get() (current catego.ID, parent catego.ID, err error) {
	return l.data[l.offset].ID, l.data[l.offset].Parent, nil
}

var treeInput = []category{
	{
		ID:     1,
		Parent: 0,
		Name:   "V1",
	},
	{
		ID:     2,
		Parent: 0,
		Name:   "V2",
	},
	{
		ID:     3,
		Parent: 0,
		Name:   "V3",
	},
	{
		ID:     4,
		Parent: 1,
		Name:   "V4",
	},
	{
		ID:     5,
		Parent: 1,
		Name:   "V5",
	},
	{
		ID:     6,
		Parent: 1,
		Name:   "V6",
	},
	{
		ID:     7,
		Parent: 2,
		Name:   "V7",
	},
	{
		ID:     8,
		Parent: 2,
		Name:   "V8",
	},
	{
		ID:     9,
		Parent: 2,
		Name:   "V9",
	},
	{
		ID:     10,
		Parent: 3,
		Name:   "V10",
	},
	{
		ID:     11,
		Parent: 3,
		Name:   "V11",
	},
	{
		ID:     13,
		Parent: 3,
		Name:   "V13",
	},
}

type ID2Category map[catego.ID]int

func main() {

	id2cats := make(ID2Category, 0)
	for k, v := range treeInput {
		id2cats[v.ID] = k
	}
	tree, err := catego.NewTree(newLoaderTest())
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\n------test 1------")
	var p []catego.ID
	p, err = tree.GetAncestors(8)
	if err != nil {
		fmt.Printf("1) expected error %v", err)
	}
	fmt.Printf("1) GetAncestors of 8 = %v\n", p)

	fmt.Println("\n------test 2------")
	p, err = tree.GetDescendants(8)
	if err != nil {
		fmt.Printf("2) expected error %v", err)
	}
	fmt.Printf("2) GetDescendants of 8 = %v\n", p)

	fmt.Println("\n------test 3------")
	p, err = tree.GetDescendants(1)
	if err != nil {
		fmt.Printf("3.1 expected error %v", err)
	}
	fmt.Printf("3.1 GetDescendants of 1 = %v\n", p)

	fmt.Printf("3.2 Get category struct\n")
	for k, v := range p {
		spew.Dump(treeInput[id2cats[v]])
		spew.Printf("3.2%d category struct of %d = %v\n", k, v, treeInput[id2cats[v]])
		fmt.Println("------------------")
	}

	fmt.Println("\n------test 4------")
	p, err = tree.GetDescendants(10)
	if err != nil {
		fmt.Printf("4) expected error %v", err)
	}
	fmt.Printf("4) GetDescendants of 10 = %v\n", p)

	fmt.Println("\n------test 5------")
	p, err = tree.GetDescendants(11)
	if err != nil {
		fmt.Printf("5) expected error %v", err)
	}
	fmt.Printf("5) GetDescendants of 11 = %v\n", p)

	var (
		b         *catego.Blacklister
		blacklist = []catego.ID{1, 10}
		whitelist = []catego.ID{11}
	)
	fmt.Println("\n------test 6------")
	b, err = tree.GetBlackLister(blacklist, whitelist)
	if err != nil {
		fmt.Printf("6) expected error %v", err)
	}
	fmt.Printf("6.1 2 is banned=%v\n", b.Is(2))
	fmt.Printf("6.2 11 is banned=%v\n", b.Is(11))
	fmt.Printf("6.3 GetBlackLister blacklist contains=%v\n", b.GetStorage().ToNums())

}

/*
output:

------test 1------
1) GetAncestors of 8 = [2 0]

------test 2------
2) GetDescendants of 8 = []

------test 3------
3.1 GetDescendants of 1 = [4 5 6]
3.2 Get category struct
(main.category) {
 ID: (catego.ID) 4,
 Parent: (catego.ID) 1,
 Name: (string) (len=2) "V4"
}
3.20 category struct of 4 = {4 1 V4}
------------------
(main.category) {
 ID: (catego.ID) 5,
 Parent: (catego.ID) 1,
 Name: (string) (len=2) "V5"
}
3.21 category struct of 5 = {5 1 V5}
------------------
(main.category) {
 ID: (catego.ID) 6,
 Parent: (catego.ID) 1,
 Name: (string) (len=2) "V6"
}
3.22 category struct of 6 = {6 1 V6}
------------------

------test 4------
4) GetDescendants of 10 = []

------test 5------
5) GetDescendants of 11 = []

------test 6------
6.1 2 is banned=true
6.2 11 is banned=false
6.3 GetBlackLister blacklist contains=[1 2 3 4 5 6 7 8 9 10 13]
*/
