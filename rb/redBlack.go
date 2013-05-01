/*
Possible iteration ideas
for e := someList.Front(); e != nil; e = e.Next() {
    v := e.Value.(T)
}

for n := tree.InitIter(gotree.InOrder); n != nil; n = tree.Next() {

    n.key ....
    n.value ....
}
*/
package gorbtree

import (
	"fmt"
	"gotree"
	"sync"
)

type color bool

const (
	Red   color = true
	Black color = false
)

type rbNode struct {
	left, right *rbNode
	key, value  interface{}
	color       color
}
type rbIter struct {
	current *rbNode
	stack   []*rbNode
	next    func() *rbNode
}

type RbTree struct {
	Height      int
	Size        int
	first, last *rbNode
	iter        rbIter
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

/*
iterativeInorder(node)
  parentStack = empty stack
  while not parentStack.isEmpty() or node != null
    if node != null then
      parentStack.push(node)
      node = node.left
    else
      node = parentStack.pop()
      visit(node)
      node = node.right
*/

func (t *RbTree) Next() *rbNode {
	return t.iter.next()
}

func (t *RbTree) InitIter(order gotree.TravOrder) *rbNode {
	current := t.root
	stack := []*rbNode{}
	switch order {
	case gotree.InOrder:
		t.iter.next = func() (out *rbNode) {

			if len(stack) > 0 || current != nil {
				if current != nil {
					stack = append(stack, current)
					current = current.left
				} else {
					stackIndex := len(stack) - 1
					out = stack[stackIndex]
					stack = stack[0:stackIndex]
					current = current.right
				}
				return out
			} else {
				return nil
			}
		}

		//case gotree.PreOrder:
		//case gotree.PostOrder:
	default:
		s := fmt.Sprintf("rbTree has not implemented %s for iteration.", order)
		panic(s)
	}
	return t.root

}
func (t *RbTree) Traverse(order gotree.TravOrder, f gotree.IterFunc) {

	n := t.root
	switch order {
	case gotree.InOrder:
		var inorder func(node *rbNode)
		inorder = func(node *rbNode) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node.key, node.value)
			inorder(node.right)
		}
		inorder(n)
	case gotree.PreOrder:
		var preorder func(node *rbNode)
		preorder = func(node *rbNode) {
			if node == nil {
				return
			}
			f(node.key, node.value)
			preorder(node.left)
			preorder(node.right)
		}

		preorder(n)

	default:
		s := fmt.Sprintf("rbTree has not implemented %s.", order)
		panic(s)
	}

}

func (t *RbTree) TraverseIter(f gotree.IterFunc) {
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

	if t.root == nil {
		t.root = &rbNode{color: Red, key: key, value: value, left: nil, right: nil}
	} else {
		//	t.root = t.insert(t.root, key, value)
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

	// setup our own stack and helpers
	var (
		stack     = []*rbNode{}
		count int = 0
		prior *rbNode
	)

L:
	for {
		switch t.cmp(h.key, key) {
		case gotree.EQ:
			h.value = value
			return t.root // no need for rest of the fix code
		case gotree.LT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.left
			if h == nil {
				h = &rbNode{color: Red, key: key, value: value}
				prior.left = h
				break L
			}
		case gotree.GT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.right
			if h == nil {
				h = &rbNode{color: Red, key: key, value: value}
				prior.right = h
				break L
			}
		default:
			panic("Compare result undefined")
		}

		if prior == h {
			panic("Shouldn't be equal last check")
		}

	}

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

		if count == 0 {
			break L2
		}

		if count > 0 {

			switch t.cmp(stack[count-1].key, h.key) {
			case gotree.LT:
				stack[count-1].left = h
			case gotree.GT:
				stack[count-1].right = h
			}
			h = stack[count-1]
		}

	}

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
