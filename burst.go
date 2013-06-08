package gotree

import (
	"bytes"
	"container/list"
)

const (
	maxLen           int = 2<<14 - 1 // two byte max
	lenOffset        int = 2
	listContainerMax int = 100
)

type container interface {
	search(data []byte) (found Byte)
	record() Byte
	setRec(b Byte) // TODO: need to check if nil that we remove our self if container is completely nil
	insert(suffix []byte, item Byte) (old Byte, burst bool, newparent *accessContainer)
	remove(suffix []byte) (old Byte)
}

type compactContainer struct {
	record  Byte // empty string case
	data    []Byte
	records []byte
}

func (c *compactContainer) search(data []byte) (found Byte) {
	if data == nil { // TODO Does this corraspond with the empty string case??
		return
	}

	checkLen := len(data)
	recLen := len(c.records)
	if checkLen == 0 { //TODO or does this mean empty string??
		return
	}

	var dend, dstart, dataCount int
	// linear scan
	for {
		// recover our 2 byte len
		dlen := c.records[dend] | c.records[dend+1]<<8
		skip := lenOffset + int(dlen)
		strRemain := c.records[dstart+lenOffset : (dend + skip)]
		dtest := bytes.Equal(strRemain, data)
		if dtest {
			found = c.data[dataCount]
			return
		}

		dend += skip
		dstart += skip
		dataCount++
		if dend == recLen {
			return
		}
	}
}

type listContainer struct {
	*list.List
	single Byte
}

type listElem struct {
	key  []byte
	item Byte
}

func (l *listContainer) search(data []byte) (found Byte) {
	for e := l.Front(); e != nil; e = e.Next() {
		if bytes.Equal(data, e.Value.(listElem).key) {
			l.MoveToFront(e)
			return e.Value.(listElem).item
		}
	}
	return
}

func (l *listContainer) record() Byte {
	return l.single
}

func (l *listContainer) setRec(b Byte) {
	if b == nil {
		l.single = b
		if l.Len() == 0 {
			l = nil
		}
		return
	}
	l.single = b
	return
}

func (l *listContainer) insert(suffix []byte, item Byte) (old Byte, burst bool, newparent *accessContainer) {
	for e := l.Front(); e != nil; e = e.Next() {
		if bytes.Equal(suffix, e.Value.(listElem).key) {
			l.MoveToFront(e)
			old = e.Value.(listElem).item
			return
		}
	}
	if l.Len() > listContainerMax {
		burst = true
		// burst
	} else {
		l.PushFront(&listElem{suffix, item})
		return
	}

	return
}
func (l *listContainer) remove(suffix []byte) (old Byte) {
	//TODO
	return
}

type accessContainer struct {
	single  Byte // empty string case
	records [256]interface{}
}

func (c *accessContainer) search(data []byte) (found interface{}) {

	return
}

type BurstTree struct {
	root     interface{}
	size     int
	height   int
	iterNext func() Byte
}

func (burst *BurstTree) Search(item Byte) (found Byte) {

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
	for i := 0; i < n; i++ {
		switch cOld := c.(type) {
		case *accessContainer:
			// use our current byte as index to next level of trie
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			return cOld.search(suffix)
		default: //nil
			return nil
		}
	}
	// we have exhausted our string empty string records
	switch cOld := c.(type) {
	case container:
		return cOld.record()
	case accessContainer:
		return cOld.single
	default:
		panic("NO!")
	}
	return
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
	n := len(query)
	if n == 0 {
		return
	}

	c := burst.root // current object
	// we need the parent for we may burst and add an access tree which needs to be linked with it's proper parent
	var parent *accessContainer
	for i := 0; i < n; i++ {
		switch cOld := c.(type) {
		case *accessContainer:
			// use our current byte as index to next level of trie
			parent = cOld
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			found, didBurst, newParent := cOld.insert(suffix, item)
			if didBurst {
				parent.records[query[i]] = newParent
			}
			old = found
		case nil: //nil
			suffix := query[i:]
			newContainer := &listContainer{list.New(), nil}
			found, didBurst, newParent := newContainer.insert(suffix, item)
			old = found
			if didBurst { // if we bursted on this it would be a binary tree
				parent.records[query[i]] = newParent
				return
			}
			parent.records[query[i]] = newContainer
		}
	}
	// we have exhausted our string insert into empty string record
	switch cOld := c.(type) {
	case container:
		old = cOld.record()
		cOld.setRec(item)
		return
	case accessContainer:
		old = cOld.single
		cOld.single = item
		return
	default:
		panic("NO!")
	}
	return
}

func (burst *BurstTree) Remove(item Byte) (old Byte) {

	// preconditions
	if item == nil {
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
	var parents []*accessContainer
	for i := 0; i < n; i++ {
		switch cOld := c.(type) {
		case *accessContainer:
			// use our current byte as index to next level of trie
			parents = append(parents, cOld)
			c = cOld.records[query[i]]
		case container:
			suffix := query[i:]
			old = cOld.remove(suffix)
			last := len(parents) - 1
		Empty:
			// continuelly check for empty access containers
			for ; last >= 0; last-- {
				parent := parents[last]
				if parent.single != nil { // not empty
					break Empty
				}
				for _, v := range parent.records {
					if v != nil {
						break Empty
					}
				}
				parent = nil // remove empty

			}
			return

		default: //nil; we need to create a new container and potentially link with parent
		}
	}
	// we have exhausted our string insert into empty string record
	switch cOld := c.(type) {
	case container:
		old = cOld.record()
		cOld.setRec(nil)
		return
	case accessContainer:
		old = cOld.single
		cOld.single = nil

		last := len(parents) - 1
	Empty2:
		// continuelly check for empty access containers
		for ; last >= 0; last-- {
			parent := parents[last]
			if parent.single != nil { // not empty
				break Empty2
			}
			for _, v := range parent.records {
				if v != nil {
					break Empty2
				}
			}
			parent = nil // remove empty

		}
		return
	default:
		panic("NO!")
	}
	return
}
