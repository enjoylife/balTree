package gotree

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"runtime"
	"testing"
)

// TODO remove
var _ = fmt.Println
var _ = spew.Dump

// container is the type which allows us to switch out leaf node containers for a burst tree.
// internal for it has alot of burst tree specific corner cases and shouldnt be considered a full dictionary stucture.
// all methods must consider empty suffix parameter
type container interface {
	search(suffix []byte) (found Byte)
	remove(suffix []byte) (old Byte)
	// must replace this containers parent if newParent != nil, this is because this method
	// might add to the tree depth if it feels the need to burst
	insert(suffix []byte, item Byte) (old Byte, newParent *accessContainer)
	// key ,value
	iter(TravOrder) func() ([]byte, Byte)
	isEmpty() bool
}

const (
	maxLen    int = 2<<14 - 1
	lenOffset int = 2
)

type compactArray struct {
	single Byte
	items  []Byte
	// length prefixed, logically seperated byte strings
	// a compact reprsentation of strings.
	records []byte
}

func (c *compactArray) sort(order TravOrder) {
	if len(c.items) == 0 {
		return
	}

	switch order {
	case AnyOrder:
		return
	}
	type sorter struct {
		suffix []byte
		item   Byte
	}
	//newRec := make([]byte, 0, len(c.records))
	tempRec := make([]sorter, len(c.items))
	dend, dstart, suffixCount, recLen := 0, 0, 0, len(c.records)

	// write into stuct array
	for {
		// compute offsets
		dlen := int((c.records[dend]) | (c.records[dend+1])<<8)
		skip := lenOffset + dlen

		tempRec[suffixCount] = sorter{c.records[dstart+lenOffset : (dend + skip)], c.items[suffixCount]}
		dend += skip
		dstart += skip //move our indexs
		suffixCount++  // keep suffix index in sync

		// search failed
		if dend == recLen {
			break
		}
	}
	//TODO  avoid sort if already sorted by checking on the prior pass of data
	switch order {
	case InOrder:

		for i := 1; i < len(c.items); i++ {
			for j := i; j > 0 && bytes.Compare(tempRec[j].suffix, tempRec[j-1].suffix) < 0; j-- {
				tempRec[j-1], tempRec[j] = tempRec[j], tempRec[j-1]
			}
		}
	case RevOrder:

		for i := 1; i < len(c.items); i++ {
			for j := i; j > 0 && bytes.Compare(tempRec[j].suffix, tempRec[j-1].suffix) > 0; j-- {
				tempRec[j], tempRec[j-1] = tempRec[j-1], tempRec[j]
			}
		}
	}

	c.records = nil
	c.items = nil
	for _, x := range tempRec {

		c.extend(x.suffix, x.item)
	}

}

func (c *compactArray) iter(order TravOrder) (fn func() ([]byte, Byte)) {
	//TODO RevOrder

	c.sort(order)
	dend, dstart, suffixCount, notDone, recLen, singleOut := 0, 0, 0, true, len(c.records), false
	return func() (key []byte, found Byte) {

		if !singleOut {
			singleOut = true
			if c.single != nil {
				return []byte{}, c.single
			}
		}
		if len(c.records) == 0 {
			return nil, nil
		}

		for notDone {
			// compute offsets
			dlen := int((c.records[dend]) | (c.records[dend+1])<<8)
			skip := lenOffset + dlen
			key = c.records[dstart+lenOffset : (dend + skip)] // get string
			found = c.items[suffixCount]
			dend += skip
			dstart += skip //move our indexs
			suffixCount++  // keep suffix index in sync

			// cant break, we need to output nil to signal were done iterating
			if dend == recLen {
				notDone = false
			}
			return
		}
		return nil, nil
	}
}

func (c *compactArray) isEmpty() bool {
	if c.single != nil || len(c.records) > 0 {
		return false
	}
	return true
}

func (c *compactArray) extend(suffix []byte, item Byte) {

	checkLen := len(suffix)
	c.records = append(c.records, byte(checkLen), byte(checkLen>>8))
	c.records = append(c.records, suffix...)
	c.items = append(c.items, item)
}

func (c *compactArray) insert(suffix []byte, item Byte) (old Byte, newParent *accessContainer) {

	// empty string case
	if len(suffix) == 0 {
		old = c.single
		c.single = item
		return
	}

	checkLen := len(suffix)
	recLen := len(c.records)

	// too big for whats inside
	if checkLen+lenOffset > recLen {
		c.extend(suffix, item)
	} else {
		dend, dstart, suffixCount := 0, 0, 0
		for {
			// compute offsets
			dlen := int((c.records[dend]) | (c.records[dend+1])<<8)
			skip := lenOffset + dlen
			strRemain := c.records[dstart+lenOffset : (dend + skip)] // get string

			if len(strRemain) == len(suffix) {
				dtest := bytes.Equal(strRemain, suffix)
				if dtest {
					old = c.items[suffixCount]
					c.items[suffixCount] = item
					return
				}
			}
			dend += skip
			dstart += skip
			// keep suffix index in sync
			suffixCount++

			if dend == recLen {
				// search failed, insert at end
				c.extend(suffix, item)
				break
			}
		}
	}
	// check if we need to burst
	if len(c.items) > containerMax {

		// add more depth to tree
		newParent = &accessContainer{}
		// Begin transfering to new depth
		var newContainer *compactArray

		// transfer empty string
		newParent.single = c.single
		// we need new lenth since we inserted prior
		recLen := len(c.records)
		dend, dstart, suffixCount := 0, 0, 0
		for {

			// compute offsets
			dlen := int((c.records[dend]) | (c.records[dend+1])<<8)
			skip := lenOffset + dlen
			elem := c.records[dstart+lenOffset : (dend + skip)] // get string

			// byte to be removed
			index := elem[0]
			// remove byte
			elem = elem[1:]

			// if we have not created a new child yet create new child
			// first check for empty string case
			if newParent.records[index] == nil {
				newContainer = &compactArray{}
				// set new child
				newParent.records[index] = newContainer
			} else {
				newContainer = newParent.records[index].(*compactArray)
			}

			if len(elem) == 0 {
				newContainer.single = c.items[suffixCount]
			} else {
				newContainer.extend(elem, c.items[suffixCount])

			}
			dend += skip
			dstart += skip
			// keep suffix index in sync
			suffixCount++
			// done
			if dend == recLen {
				break

			}
		}
		c = nil
		// remove our dead mem now, hopefully runtime will compact the scattered memory
		runtime.GC()

	}
	return
}

func (c *compactArray) search(suffix []byte) (found Byte) {
	// take care of empty string case
	if len(suffix) == 0 {
		return c.single
	}
	if len(suffix) > maxLen {
		return
	}

	checkLen := len(suffix)
	recLen := len(c.records)

	// too big for whats inside
	if checkLen+lenOffset > recLen {
		return
	}

	var dend, dstart, suffixCount int
	for {
		// compute offsets
		dlen := int((c.records[dend]) | (c.records[dend+1])<<8)
		skip := lenOffset + dlen
		strRemain := c.records[dstart+lenOffset : (dend + skip)] // get string
		// TODO is this duplicating work by bytes.Equal?
		if len(strRemain) == len(suffix) {
			dtest := bytes.Equal(strRemain, suffix)
			if dtest {
				found = c.items[suffixCount]
				return
			}
		}
		dend += skip
		dstart += skip //move our indexs
		suffixCount++  // keep suffix index in sync

		// search failed
		if dend == recLen {
			return
		}
	}
}

func (c *compactArray) remove(suffix []byte) (found Byte) {
	// take care of empty string case
	if len(suffix) == 0 {
		found = c.single
		c.single = nil
		return
	}

	checkLen := len(suffix)
	recLen := len(c.records)

	// too big for whats inside
	if checkLen+lenOffset > recLen {
		return
	}

	var dend, dstart, suffixCount int
	for {
		// compute offsets
		dlen := int((c.records[dend]) | (c.records[dend+1])<<8)
		skip := lenOffset + dlen
		strRemain := c.records[dstart+lenOffset : (dend + skip)] // get string
		// TODO is this duplicating work by bytes.Equal?
		if len(strRemain) == len(suffix) {
			dtest := bytes.Equal(strRemain, suffix)
			if dtest {
				found = c.items[suffixCount]

				// remove
				c.records = append(
					c.records[:dstart],
					c.records[(dend+skip):]...,
				)
				c.items = append(
					c.items[:suffixCount],
					c.items[suffixCount+1:]...,
				)
				return
			}
		}
		dend += skip
		dstart += skip //move our indexs
		suffixCount++  // keep suffix index in sync

		// search failed
		if dend == recLen {
			return
		}
	}
}

var listContainerMax int = 150

func optimize() func(b *testing.B) {

	l := list.New()
	for i := 0; i < listContainerMax; i++ {
		l.PushFront(i)
	}
	origMax := listContainerMax
	return func(b *testing.B) {
		b.StopTimer()
		b.ResetTimer()
		// add upto new listContainerMax
		for i := origMax; i < listContainerMax; i++ {
			l.PushFront(i)
		}
		origMax = listContainerMax
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for e := l.Front(); e != nil; e = e.Next() {

			}
		}

	}
}

func init() {
	//TODO Optimize listContainerMax prior to use
	optimized := false
	benchFunc := optimize()
	for !optimized {
		// we look for a linear search that will take around 200 ns
		result := testing.Benchmark(benchFunc).NsPerOp()
		if result > 200 {
			optimized = true
		} else {
			listContainerMax += 50
			//fmt.Println("result:", result, "size:", listContainerMax)
		}
	}
}

type listContainer struct {
	*list.List
	single Byte // empty byte holder
}

type listElem struct {
	key  []byte
	item Byte
}
type listElemSlice []listElem

func (p listElemSlice) Len() int           { return len(p) }
func (p listElemSlice) Less(i, j int) bool { return bytes.Compare(p[i].key, p[j].key) <= 0 }
func (p listElemSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (l *listContainer) search(suffix []byte) (found Byte) {
	// take care of empty string case
	if len(suffix) == 0 {
		return l.single
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if bytes.Equal(suffix, e.Value.(*listElem).key) {
			l.MoveToFront(e)
			return e.Value.(*listElem).item
		}
	}
	return
}

func (l *listContainer) isEmpty() bool {
	if l.single != nil || l.Len() > 0 {
		return false
	}
	return true
}

func (l *listContainer) iter(order TravOrder) (fn func() ([]byte, Byte)) {
	// TODO Sort
	for e := l.Front(); e != nil; e = e.Next() {
	}

	return
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
	if l.Len() > containerMax {
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
