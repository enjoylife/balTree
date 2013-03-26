package gotree

import (
	"fmt"
	"sync"
	"unsafe" // for sizeof
)

type avlNode struct {
	parent, left, right *avlNode
	key                 interface{}
	item                interface{}
	balance             int
}

func (n *avlNode) Next() *avlNode {
	var p *avlNode
	if n.right != nil {
		return n.right.minChild()
	}
	p = n.parent
	for p != nil && p.right == n {
		n = p
		p = n.parent
	}
	return p
}

func (n *avlNode) Prev() *avlNode {
	var p *avlNode
	if n.left != nil {
		return n.left.maxChild()
	}
	p = n.parent
	for p != nil && p.right == n {
		n = p
		p = n.parent
	}
	return p
}
func (n *avlNode) Key() interface{} {
	return n.key
}

func (n *avlNode) Value() interface{} {
	return n.item
}

func (n *avlNode) minChild() *avlNode {
	for n.left != nil {
		n = n.left
	}
	return n

}

func (n *avlNode) maxChild() *avlNode {
	for n.right != nil {
		n = n.right
	}
	return n
}

type AvlTree struct {
	Height int
	Size   int

	root, first, last *avlNode

	cmp  CompareFunc
	lock sync.RWMutex
}

func New(cmp CompareFunc) *AvlTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &AvlTree{root: nil, first: nil, last: nil, cmp: cmp}
}

func (t *AvlTree) First() (interface{}, bool) {
	if t.first != nil {
		return t.first, true
	}
	return nil, false
}
func (t *AvlTree) Last() (interface{}, bool) {
	if t.last != nil {
		return t.last, true
	}

	return nil, false
}

func (t *AvlTree) rotateLeft(n *avlNode) {
	var (
		p      = n
		q      = n.right
		parent = p.parent
	)
	if p.parent != nil {
		if parent.left == p {
			parent.left = q
		} else {
			parent.right = q
		}
	} else {
		t.root = q
	}
	q.parent = parent
	p.parent = q
	p.right = q.left
	if p.right != nil {
		p.right.parent = p
	}
	q.left = p
}

func (t *AvlTree) rotateRight(n *avlNode) {
	var (
		p      = n
		q      = n.left
		parent = p.parent
	)
	if p.parent != nil {
		if parent.left == p {
			parent.left = q
		} else {
			parent.right = q
		}
	} else {
		t.root = q
	}
	q.parent = parent
	p.parent = q
	p.left = q.right
	if p.left != nil {
		p.left.parent = p
	}
	q.right = p

}

func (t *AvlTree) Space(format string) float64 {

	const prior = unsafe.Sizeof(*t)
	const nodeSize = unsafe.Sizeof(t.root)
	bytes := float64(prior + nodeSize*uintptr(t.Size))
	switch format {
	case "KiB":
		return bytes / (2 << 9)
	case "kB":
		return bytes / 1000
	case "MiB":
		return bytes / (2 << 19)
	case "MB":
		return bytes / 1000000
	default:
		return -1.0
	}
	return 0
}

func (t *AvlTree) Search(key interface{}) (item interface{}, ok bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	var node *avlNode = t.root
	for node != nil {
		resp := t.cmp(node.key, key)
		if resp == 0 {
			return node.item, true
		}
		if resp > 0 {
			node = node.left
		} else {
			node = node.right
		}

	}
	return nil, false
}

func (t *AvlTree) Insert(key interface{}, item interface{}) (old interface{}, ok bool) {
	if key == nil {
		ok = false
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()

	var (
		node       *avlNode = t.root
		parent     *avlNode = nil
		unbalanced *avlNode = node
		isLeft     bool     = false
	)
	old = nil

	// search
	for node != nil {
		if node.balance != 0 {
			unbalanced = node
		}
		resp := t.cmp(node.key, key)
		if resp == 0 {
			old = node.item
			break
		}
		parent = node
		if isLeft = (resp > 0); isLeft {
			node = node.left
		} else {
			node = node.right
		}

	}
	// add if only a new insertion
	if old == nil {
		t.Size++
		node = &avlNode{
			key: key, item: item,
			parent: parent, left: nil, right: nil,
			balance: 0,
		}
	} else {
		node.item = item
		return
	}

	// base case
	if parent == nil {
		t.root = node
		t.first = node
		t.last = node
		t.Height++
		return old, true
	}

	// link to tree
	if isLeft {
		parent.left = node
	} else {
		parent.right = node
	}

	// Maintain first and last pointers
	if isLeft {
		if parent == t.first {
			t.first = node
		}
	} else {
		if parent == t.last {
			t.last = node
		}
	}

	// fix balances on our way up to unbalanced node
	for {
		switch node {
		case parent.left:
			parent.balance--
		case parent.right:
			parent.balance++
		default:
			panic("SHOULDNT")
		}
		if parent == unbalanced {
			break
		}
		node = parent
		parent = parent.parent
	}

	switch unbalanced.balance {
	case 0:
	case 1, -1:
		t.Height++
	case 2:
		right := unbalanced.right
		if right.balance == 1 {
			unbalanced.balance = 0
			right.balance = 0
			//	} else if right.left == nil {
		} else {
			if right.left == nil {
				panic("WE GOT A NIL")
			}
			switch right.left.balance {
			case 1:
				unbalanced.balance = -1
				right.balance = 0
			case 0:
				unbalanced.balance = 0
				right.balance = 0
			case -1:
				unbalanced.balance = 0
				right.balance = 1
			}
			right.left.balance = 0
			t.rotateRight(right)
		}
		t.rotateLeft(unbalanced)
	case -2:
		left := unbalanced.left
		if left.balance == -1 {
			unbalanced.balance = 0
			left.balance = 0
		} else if left.right == nil {
		} else {
			if left.right == nil {
				panic("WE GOT A NIL")
			}
			switch left.right.balance {
			case 1:
				unbalanced.balance = 0
				left.balance = -1
			case 0:
				unbalanced.balance = 0
				left.balance = 0
			case -1:
				unbalanced.balance = 1
				left.balance = 0
			}
			left.right.balance = 0
			t.rotateLeft(left)
		}
		t.rotateRight(unbalanced)

	default:
		fmt.Println(unbalanced.balance)
	}

	return old, true
}

func (t *AvlTree) Remove(key interface{}) (old interface{}, ok bool) {
	t.lock.Lock()
	defer t.lock.Unlock()

	var (
		node   *avlNode = t.root
		parent *avlNode = nil
		isLeft bool     = false
	)
	old = nil

	// search
	for node != nil {
		resp := t.cmp(node.key, key)
		if resp == 0 {
			old = node.item
			break
		}
		parent = node
		if isLeft = (resp > 0); isLeft {
			node = node.left
		} else {
			node = node.right
		}

	}

	// didnt fint it
	if old == nil {
		return nil, false
	}

	t.Size--
	var (
		left  *avlNode = node.left
		right *avlNode = node.right
		next  *avlNode
	)
	if node == t.first {
		t.first = node.Next()
	}
	if node == t.last {
		t.last = node.Prev()
	}

	if left == nil {
		next = right
	} else if right == nil {
		next = left
	} else {
		next = right.minChild()
	}
	if parent != nil {
		isLeft = (parent.left == node)
		if isLeft {
			parent.left = next
		} else {
			parent.right = next
		}
	} else {
		t.root = next
	}

	if left != nil && right != nil {
		next.balance = node.balance
		next.left = left
		left.parent = next

		if next != right {
			parent = next.parent
			next.parent = node.parent

			node = next.right
			parent.left = node
			isLeft = true

			next.right = right
			right.parent = next
		} else {
			next.parent = parent
			parent = next
			node = parent.right
			isLeft = false
		}
	} else {
		node = next
	}
	if node != nil {
		node.parent = parent
	}

	for parent != nil {
		balance := 0
		node = parent
		parent = parent.parent
		if isLeft {
			isLeft = (parent != nil && parent.left == node)
			node.balance++
			balance = node.balance
			if balance == 0 {
				continue
			}
			if balance == 1 {
				return old, true
			}
			right = node.right
			switch right.balance {
			case 0:
				node.balance = 1
				right.balance = -1
				t.rotateLeft(node)
				return old, true
			case 1:
				node.balance = 0
				right.balance = 0
			case -1:
				switch right.left.balance {
				case 1:
					node.balance = -1
					right.balance = 0
				case 0:
					node.balance = 0
					right.balance = 0
				case -1:
					node.balance = 0
					right.balance = 1
				}
				right.left.balance = 0
				t.rotateRight(right)
			}
			t.rotateLeft(node)
		} else {
			isLeft = (parent != nil && (parent.left == node))
			node.balance--
			balance = node.balance
			if balance == 0 {
				continue
			}
			if balance == -1 {
				return old, true
			}
			left = node.left
			switch left.balance {
			case 0:
				node.balance = -1
				left.balance = 1
				t.rotateRight(node)
				return old, true
			case -1:
				node.balance = 0
				left.balance = 0
			case 1:
				switch left.right.balance {
				case 1:
					node.balance = 0
					left.balance = -1
				case 0:
					node.balance = 0
					left.balance = 0
				case -1:
					node.balance = 1
					left.balance = 0
				}
				left.right.balance = 0
				t.rotateLeft(left)
			}
			t.rotateRight(node)
		}
	}
	t.Height--

	return old, true
}

func (t *AvlTree) Traverse(f IterFunc) {
	var node *avlNode = t.first
	if t.root == nil {
		return
	}
	for {
		f(node.key, node.item)
		if node == t.last {
			break
		}
		node = node.Next()
	}
}

func (t *AvlTree) Print() {
	var node *avlNode = t.first
	fmt.Println("First", node)
	if t.root == nil {
		return
	}
	for {
		fmt.Printf("Key: %v Value: %v \n", node.key, node.item)
		if node == t.last {
			break
		}
		node = node.Next()
	}
}
