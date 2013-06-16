// TODO: accept initially empty strings for search, insertion and removal
package gotree

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"runtime"
)

var _ = spew.Dump
var _ = fmt.Println

const (
//maxLen int = 2<<14 - 1 // two byte max
)

var listContainerMax int = 2

func init() {
	//TODO Optimize listContainerMax prior to use
}

// All methods must handle empty string check too!!
type container interface {
	search(data []byte) (found Byte)
	remove(suffix []byte) (old Byte)
	// must replace this containers parent if newParent != nil, this is because this method
	// might add to the tree depth if it feels the need to burst
	insert(suffix []byte, item Byte) (old Byte, newParent *accessContainer)

	hasElem() bool
}

type listContainer struct {
	*list.List
	single Byte // empty byte holder
}

type listElem struct {
	key  []byte
	item Byte
}

func (l *listContainer) search(data []byte) (found Byte) {
	// take care of empty string case
	if len(data) == 0 {
		return l.single
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if bytes.Equal(data, e.Value.(*listElem).key) {
			l.MoveToFront(e)
			return e.Value.(*listElem).item
		}
	}
	return
}

func (l *listContainer) hasElem() bool {
	if l.single != nil || l.Len() > 0 {
		return true
	}
	return false
}

func (l *listContainer) insert(suffix []byte, item Byte) (old Byte, newParent *accessContainer) {
	if len(suffix) == 0 {
		// empty string case
		old = l.single
		l.single = item
		return
	}
	// search for previous old entry to return
	for e := l.Front(); e != nil; e = e.Next() {
		if bytes.Equal(suffix, e.Value.(*listElem).key) {
			l.MoveToFront(e)
			old = e.Value.(*listElem).item
			return
		}
	}
	// not found so add it in
	l.PushFront(&listElem{suffix, item})

	// check if we need to burst
	if l.Len() > listContainerMax {
		// add more depth to tree
		newParent = &accessContainer{}
		// transfer empty string
		newParent.single = l.single
		// transfer the rest
		for e := l.Front(); e != nil; e = e.Next() {
			elem := e.Value.(*listElem)
			// byte to be removed
			index := elem.key[0]
			// remove byte
			elem.key = elem.key[1:]
			// if we have not created a new child yet create new child
			// first check for empty string case
			if newParent.records[index] == nil {
				newContainer := &listContainer{list.New(), nil}
				if len(elem.key) == 0 {
					newContainer.single = elem.item
				} else {
					newContainer.PushBack(elem)

				}
				// set new child
				newParent.records[index] = newContainer
			} else {
				if len(elem.key) == 0 {
					newParent.records[index].(*listContainer).single = elem.item
				} else {
					newParent.records[index].(*listContainer).PushBack(elem)

				}
			}

		}
		l = nil
		// remove our dead mem now, hopefully runtime will compact the scattered memory
		runtime.GC()
	}

	return
}

func (l *listContainer) remove(suffix []byte) (old Byte) {
	if len(suffix) == 0 {
		// empty string case
		old = l.single
		l.single = nil
		return
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if bytes.Equal(suffix, e.Value.(*listElem).key) {
			old = l.Remove(e).(*listElem).item
			return
		}
	}
	return
}

// doesn't follow container interface for it's not a container but a "trie node"
type accessContainer struct {
	single  Byte             // empty string case
	records [256]interface{} // may be a accessContainer or container
}

type BurstTree struct {
	root     interface{}
	size     int
	height   int
	iterNext func() Byte
}

func (burst *BurstTree) Clear() {
	burst.root = nil
	burst.size = 0
	burst.height = 0
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

	c := burst.root // current object
	parent := burst.root.(*accessContainer)
	// traverse down tree
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

	c := burst.root // current object
	// we need the parent for we may burst and add an access tree which needs to be linked with it's proper parent
	parent := burst.root.(*accessContainer)
	for i := 0; i <= n; i++ {
		switch cOld := c.(type) {
		case *accessContainer:
			// use our current byte as index to next level of trie
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
		case nil: // need to add new
			suffix := query[i:]
			// TODO: Try other concrete types of containers
			newContainer := &listContainer{list.New(), nil}
			old, _ = newContainer.insert(suffix, item)
			parent.records[query[i-1]] = newContainer
			burst.size++
			return
		}
	}
	// assumes that the container took care of the empty string case of insertion
	// we have exhausted our string insert into empty string record for our last accessContainer
	old = parent.single
	if old == nil {
		burst.size++
	}
	parent.single = item
	return
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
					//fmt.Println("empty String removing", item, "for", query[i-1])
					burst.size--
					cOld.single = nil
					goto CheckEmpty
				}
				return // found nothing
			}
			parent = cOld
			parents = append(parents, cOld)
			// use our current byte as index to next level of trie
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			old = cOld.remove(suffix)
			if old != nil {
				burst.size--
				if !cOld.hasElem() {
					//fmt.Println("removing empty string", item, "in", query[i-1])
					// remove empty container
					parent.records[query[i-1]] = nil
					//parents[len(parents)-1].records[query[i-1]] = nil
				}
				goto CheckEmpty
			}
			return // found nothing

		case nil:
			//fmt.Println("Found Nil for", item)
			// if present only possible place is last access container
			old = parent.single
			//parent := parents[len(parents)-1]
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
