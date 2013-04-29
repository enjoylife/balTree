package gorbtree

import (
	"../"
	"fmt"
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
			node = node.left
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.key, node.value)
			stack = stack[0:stackIndex]
			//fmt.Println("stack size", len(stack))
			node = node.right
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
	t.root = t.insert(t.root, key, value)
	t.root.color = Black
	return
}

func (t *RbTree) insert(h *rbNode, key interface{}, value interface{}) *rbNode {

	if h == nil {
		return &rbNode{color: Red, key: key, value: value}
	}
	switch cmp := t.cmp(h.key, key); {
	case cmp == 0:
		h.value = value
	case cmp > 0:
		h.left = t.insert(h.left, key, value)
	default:
		h.right = t.insert(h.right, key, value)
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

func (t *RbTree) InsertIter(key interface{}, value interface{}) (old interface{}, ok bool) {

	if key == nil {
		ok = false
		return
	}
	t.lock.Lock() // for writing
	defer t.lock.Unlock()
	h := t.root
	type path struct {
		node *rbNode
		dir  bool
	}
	stack := []path{}
	for h != nil {
		if h.left.isRed() && h.right.isRed() {
			h.colorFlip()
		}
		cmp := t.cmp(h.key, key)
		if cmp == 0 {
			//fmt.Println("Old value")
			old = h.value
			h.value = value
			ok = true
			break

		} else if cmp > 0 {
			stack = append(stack, path{h, Left})
			//fmt.Println("Going left")
			h = h.left
		} else {
			stack = append(stack, path{h, Right})
			fmt.Println("Going right")
			h = h.right
		}
	}
	//fmt.Printf("Stack: %v\n", stack)
	if h == nil {
		h = &rbNode{color: Red, key: key, value: value, left: nil, right: nil}
		//stack = append(stack, h)
	} else {
		fmt.Println("Have old")
	}

	stackHeight := len(stack) - 1
	//fmt.Println("Height", stackHeight)
	for i := stackHeight; i >= 0; i-- {
		//fmt.Println("Iter", i, "value", h.value)
		parent := stack[i]
		if h == nil || parent.node == nil {
			s := fmt.Sprintf("I: %d", i)
			panic(s)
		}
		if h.right.isRed() && !h.left.isRed() {
			//fmt.Println("Rotate Left")
			h = h.rotateLeft()
		}
		if h.left.isRed() && h.left.left.isRed() {
			//fmt.Println("Rotate Right")
			h = h.rotateRight()
		}
		if parent.dir == Left {
			parent.node.left = h
		} else {
			parent.node.right = h
		}
		h = parent.node
	}
	t.root = h
	t.root.color = Black

	return
}

func (h *rbNode) isRed() bool {
	return h != nil && h.color == Red
}

func (h *rbNode) colorFlip() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
