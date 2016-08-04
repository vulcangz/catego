# Catego

[![Build Status](https://travis-ci.org/mrsinham/catego.png?branch=master)](https://travis-ci.org/mrsinham/catego)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrsinham/catego)](https://goreportcard.com/report/github.com/mrsinham/catego)

Catego is a simple library to manage categories (universe, subcat, etc) tree or any tree structure that has the same implementation. The datastructure is threadsafe and has those specs :

* One node can have only one parent
* One node can have many children

Assuming this with this library you can :

* Get any node in the tree with an almost O(1) method (map+mutex)
* Get all children of a node (descending the tree)
* Get all parents list (climbing the tree to root node)
* Get all node of the tree except the given nodes (exclusion)
* Create a blacklister that will says instantly if a node is banned from the list you provided (if you submit a blacklisted node, all its children will be banned too)


## Example use-case

This library has been made to represent a tree of categories in RAM and to know if a category is banned very quickly. If you ban an upper category, every category beneath will be banned as well.

## Documentation

https://godoc.org/github.com/mrsinham/catego

## Installation

This library uses cool libraries as (thanks to the authors) :

* github.com/Workiva/go-datastructures
* github.com/juju/errgo

So include them in your path/vendor directories, or use directly :

```bash
$ go get gopkg.in/mrsinham/catego.v1
```

## Usage

### Using the tree

To use it, you just need to create a NodeSource structure that will satisfies this :

```go
type NodeSource interface {
	// Will be used as a condition for a loop
	Next() bool
	// Return the current node id. If parent is 0 then it is a root node
	Get() (current ID, parent ID, err error)
}
```

And now you can just do :

```go
var cs NodeSource
cs = MyCustomNodeSourceCreator()

if t, err := NewTree(cs); err != nil{
	// catch err
}

var ch []ID
ch, err = t.GetDescendants(22)
if err != nil {
	// catch er
}

// ch has now all the children of node 22.
```

#### Notes :

The root node is 0. If you want to tells the library that the node is the highest node on the tree, you need to says that its parent is 0 (but you can change it with the TreeOption).

### Blacklist

And for blacklist you can do :

```go 
var cs NodeSource
cs = MyCustomNodeSourceCreator()

if t, err := NewTree(cs); err != nil{
	// catch err
}

var b *Blacklister
var err error
b, err = t.GetBlacklister([]ID{22,8,17}, nil)
if err != nil {
	// catch error
}

// if 272 is a child (of any level) of 22,8, 17, banned = true
banned := b.Is(272)
```

Because Blacklister is a top level structure on a bitarray, once it has been created, the lookup to know if a node is banned is close to O(1).

#### Notes :

Whitelist means "ban everything but", it is not equal as a safe list because blacklisted categories are always strongest than whitelisted categories.


## Contribute

PR and issues are welcome :)
