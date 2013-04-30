package gorbtree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gotree"
	"sync"
)

var _ = fmt.Println

const (
	Red   = true
	Black = false
	Left  = true
	Right = false
)

type rbNode struct {
	left, right *rbNode
	key, value  interface{}
	color       bool
}

type RbTree struct {
	Height      int
	Size        int
	first, last *rbNode
	root        *rbNode
	cmp         gotree.CompareFunc
	lock        sync.RWMutex
}

func (n *rbNode) Key() interface{} {
	return n.key
}

func (n *rbNode) Value() interface{} {
	return n.value
}

func (n *rbNode) MinChild() *rbNode {
	for n.left != nil {
		n = n.left
	}
	return n

}

func (n *rbNode) MaxChild() *rbNode {
	for n.right != nil {
		n = n.right
	}
	return n
}

// Cant modify the tree as we go down
// for if the stack gets out of sync as we head
// back up it will be all bad

func (t *RbTree) Traverse(f gotree.IterFunc) {
	node := t.root
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.right
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.key, node.value)
			stack = stack[0:stackIndex]
			//fmt.Println("stack size", len(stack))
			node = node.left
		}
	}
}
func (h *rbNode) rotateLeft() (x *rbNode) {
	x = h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = Red
	return
}
func (h *rbNode) rotateRight() (x *rbNode) {
	x = h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = Red
	return
}

func New(cmp gotree.CompareFunc) *RbTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &RbTree{root: nil, first: nil, last: nil, cmp: cmp}
}

func (t *RbTree) Search(key interface{}) (value interface{}, ok bool) {

	if key == nil {
		return

	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	x := t.root
	for x != nil {
		cmp := t.cmp(x.key, key)
		if cmp == 0 {
			return x.value, true
		} else if cmp > 0 {
			x = x.left
		} else {
			x = x.right
		}
	}
	return
}

func (t *RbTree) Insert(key interface{}, value interface{}) (old interface{}, ok bool) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.root == nil {
		t.root = &rbNode{color: Red, key: key, value: value, left: nil, right: nil}
	} else {
		//t.root = t.insert(t.root, key, value)
		t.root = t.insertIter(t.root, key, value)
	}
	if t.root.color == Red {
		t.Height++
	}
	t.root.color = Black
	return
}

func (t *RbTree) insert(h *rbNode, key interface{}, value interface{}) *rbNode {

	// empty tree
	if h == nil {
		return &rbNode{color: Red, key: key, value: value}
	}

	switch t.cmp(h.key, key) {
	case gotree.EQ:
		h.value = value
	case gotree.LT:
		h.left = t.insert(h.left, key, value)
	case gotree.GT:
		h.right = t.insert(h.right, key, value)
	default:
		panic("Compare result of undefined")
	}

	if h.right.isRed() && !(h.left.isRed()) {
		h = h.rotateLeft()
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = h.rotateRight()
	}

	if h.left.isRed() && h.right.isRed() {
		h.colorFlip()
	}
	return h
}

func (t *RbTree) insertIter(h *rbNode, key interface{}, value interface{}) *rbNode {

	// empty tree
	if h == nil {
		return &rbNode{color: Red, key: key, value: value}
	}

	// we need to store which way we took on the way down,
	// so we can reconnect nodes parents
	type path struct {
		node *rbNode
		dir  gotree.Direction
	}

	// setup our own stack and helpers
	var (
		stack        = []path{}
		count    int = 0
		whichWay gotree.Direction
		prior    *rbNode
	)

L:
	for {
		if h == nil {
			panic("FOUND NULL AT START")
		}
		fmt.Printf("Node %d, Key %d\n", h.key, key)
		fmt.Println("Count", count)
		switch t.cmp(h.key, key) {
		case gotree.EQ:
			h.value = value
			fmt.Println("RETURN ROOT EQUAL")
			return t.root // no need for rest of the fix code
		case gotree.LT:
			whichWay = gotree.LT
			prior = h
			stack = append(stack, path{node: prior, dir: whichWay})
			count++
			h = h.left
			if prior == h {
				panic("Shouldn't be equal")
			}
		case gotree.GT:
			whichWay = gotree.GT
			prior = h
			stack = append(stack, path{node: prior, dir: whichWay})
			count++
			h = h.right
			if prior == h {
				panic("Shouldn't be equal")
			}
		default:
			panic("Compare result undefined")
		}

		// we found our spot to insert, create our node and linke to parent
		if h == nil {
			h = &rbNode{color: Red, key: key, value: value, left: nil, right: nil}
			switch whichWay {
			case gotree.LT:
				prior.left = h
			case gotree.GT:
				prior.right = h
			default:
				panic("Unknown direction")
			}
			fmt.Println("LEAVING NIL")
			break L
		}
		if prior == h {
			panic("Shouldn't be equal last check")
		}

	}

	fmt.Println("Before Fix")
	spew.Dump(t.root)

	// h is parent of new node at this point
	h = prior
L2:
	for {
		count--

		if h.right.isRed() && !(h.left.isRed()) {
			h = h.rotateLeft()
		}
		if h.left.isRed() && h.left.left.isRed() {
			h = h.rotateRight()
		}
		if h.left.isRed() && h.right.isRed() {
			h.colorFlip()
		}
		if prior == nil {
			panic("WE GOT A NIL")
		}

		if count == 0 {
			break L2
		}

		if count > 0 {
			switch stack[count-1].dir {
			case gotree.LT:
				stack[count-1].node.left = h
			case gotree.GT:
				stack[count-1].node.right = h
			}
			h = stack[count-1].node
		}

	}

	fmt.Println("After Fix")
	spew.Dump(t.root)
	fmt.Println("DONE")
	return h
}

func (h *rbNode) isRed() bool {
	return h != nil && h.color == Red
}

func (h *rbNode) colorFlip() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
