package gotree

// TODO: accept initially empty strings for search, insertion and removal

import (
	//"container/list"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"runtime"
)

var _ = spew.Dump
var _ = fmt.Println
var containerMax int

func init() {
	containerMax = listContainerMax
}

// doesn't follow container interface for it's not a container but a "trie node"
type accessContainer struct {
	single  Byte             // empty string case
	records [256]interface{} // may be a accessContainer or container
}

type BurstTree struct {
	root     interface{}
	size     int
	iterNext func() Byte
}

func (burst *BurstTree) Clear() {
	burst.root = nil
	burst.size = 0
	burst.iterNext = nil
	runtime.GC()
}

func (burst *BurstTree) Search(item Byte) (found Byte) {

	// preconditions
	if item == nil || burst.root == nil {
		return
	}
	query := item.ToBytes()
	if query == nil {
		// bad method call from Byte type
		return
	}
	n := len(query)
	if n == 0 {
		// don't accept empty case for first access container
		return
	}

	c := burst.root                // interface
	parent := c.(*accessContainer) // must always be non nil
	for i := 0; ; i++ {
		switch cOld := c.(type) {
		case *accessContainer:
			// empty string case
			if i == n {
				return cOld.single
			}
			parent = cOld
			// use our current byte as index to next level of trie
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			// needs to handle suffix being an empty string case!!
			return cOld.search(suffix)
		case nil:
			// only possible place is last access container
			return parent.single
		}
	}
}

func (burst *BurstTree) Insert(item Byte) (old Byte) {

	// preconditions
	if item == nil {
		return
	}
	query := item.ToBytes()

	if query == nil {
		return
	}

	// need a non nil parent for traversal
	if burst.root == nil {
		burst.root = &accessContainer{}
	}

	n := len(query)
	if n == 0 {
		return
	}

	c := burst.root

	// We need the parent for we may burst and add an
	// access tree which needs to be linked with it's proper parent.
	parent := burst.root.(*accessContainer)

	for i := 0; ; i++ {
		switch cOld := c.(type) {
		case *accessContainer:

			// empty string case
			if i == n {
				old = cOld.single
				cOld.single = item
				if old == nil {
					burst.size++
				}
				return
			}
			parent = cOld
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			found, newParent := cOld.insert(suffix, item)
			if newParent != nil {
				parent.records[query[i-1]] = newParent
			}
			if found == nil {
				burst.size++
			}
			return found
		case nil:
			var newContainer container
			suffix := query[i:]
			newContainer = &compactArray{}
			//newContainer := &listContainer{list.New(), nil} // TODO: Try other concrete types of containers
			old, _ /*Should never burst,or else it's just a simple trie */ = newContainer.insert(suffix, item)
			parent.records[query[i-1]] = newContainer
			burst.size++
			return
		}
	}
}

func (burst *BurstTree) Size() int {
	return burst.size
}

func (burst *BurstTree) Remove(item Byte) (old Byte) {

	// preconditions
	if item == nil || burst.root == nil {
		return
	}
	query := item.ToBytes()
	if query == nil {
		return
	}
	n := len(query)
	if n == 0 {
		return
	}

	c := burst.root // current object
	// we need the parents for we may fully empty access containers which may trigger more removes in prior depths
	parents := []*accessContainer{}
	parent := burst.root.(*accessContainer)
	for i := 0; ; i++ {
		switch cOld := c.(type) {
		case *accessContainer:
			// empty string case
			if i == n {
				old = cOld.single
				if old != nil {
					burst.size--
					cOld.single = nil
					goto CheckEmpty
				}
				return // found nothing
			}
			parent = cOld
			parents = append(parents, cOld)
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			old = cOld.remove(suffix)
			if old != nil {
				burst.size--
				if cOld.isEmpty() {
					// remove empty container
					parent.records[query[i-1]] = nil
				}
				goto CheckEmpty
			}
			return // found nothing

		case nil:
			// if present only possible place is last access container
			old = parent.single
			if old != nil {
				burst.size--
				parent.single = nil
				goto CheckEmpty
			}
			return // found nothing

		}
	}
	panic("Shouldn't be here")

CheckEmpty:
	// continually check for empty access containers
	for last := len(parents) - 1; last > 0; last-- {
		parent := parents[last]
		if parent.single != nil { // not empty, stop
			return
		}
		for _, v := range parent.records {
			if v != nil { // not empty, stop
				return
			}
		}
		// remove
		parents[last-1].records[query[last-1]] = nil

	}
	return
}

func (burst *BurstTree) Next() (next Byte) {
	return burst.iterNext()
}

type iter struct {
	index int
	it    *accessContainer
}

func (burst *BurstTree) IterInit(order TravOrder) (start Byte) {
	//TODO: test and corner case elmination
	//TODO: output key as well
	if burst.root == nil {
		start = nil
		return
	}
	var cIter func() ([]byte, Byte)
	// should we output from a container
	isC := false

	current := burst.root.(*accessContainer)
	stack := []iter{}

	index := -1
	switch order {
	case InOrder:
		burst.iterNext = func() (out Byte) {
			// we need to keep trying to go down levels, once we hit either a nill or container,
			// we need to either output all the containers items in order or ignore the nil.
			// Then continue traversing the rest of the record array.
			// We have to pay attention to the current record index as we continue going down the levels
			// so when we come back up we are at the spot we were when we started traversing down.
		Dive:
			for {
				// output containers items first
				if isC {
					// TODO stop ignoring key
					_, out = cIter()
					if out != nil {
						return
					} else {
						isC = false
					}

				}

				// output an empty string before more traversal
				if index == -1 {
					index++
					if current.single != nil {
						out = current.single
						break
					}
				}
				for index < 255 {
					switch cur := current.records[index].(type) {
					case *accessContainer:

						index++
						stack = append(stack, iter{index, current})
						current = cur // go down one more level
						index = -1
						goto Dive
					case container:
						cIter = cur.iter(order)
						if cIter != nil {
							isC = true
						}
						index++
						goto Dive
					case nil:
						index++
					}
				}
				if len(stack) > 0 {
					// pop
					stackIndex := len(stack) - 1
					s := stack[stackIndex]
					current, index = s.it, s.index
					stack = stack[0:stackIndex]
				} else {
					out = nil
					// last node, reset
					burst.iterNext = nil
					return

				}
			}
			return out
		}
		return burst.iterNext()
	case RevOrder:
		//TODO
	}
	return
}
func (burst *BurstTree) Map(order TravOrder, f ByteIterFunc) {
	//TODO
}
