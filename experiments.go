package gotree

/* Experiments: not used */
/*
import (
	"fmt"
)


func (t *RBTree) searchRecurse(key interface{}) (value interface{}, ok bool) {

	if key == nil {
		return
	}
	return t.get(t.root, key)
}

func (t *RBTree) get(node *Node, key interface{}) (value interface{}, ok bool) {
	if node == nil {
		return nil, false
	}

	switch t.cmp(node.Key, key) {
	case EQ:
		return node.Value, true
	case LT:
		return t.get(node.left, key)
	case GT:
		return t.get(node.right, key)
	default:
		panic("Compare result of undefined")
	}
}
func (t *RBTree) initIterRecurse(order TravOrder) <-chan *Node {

	n := t.root
	t.iterChan = make(chan *Node)
	switch order {
	case InOrder:
		var inorder func(node *Node)
		inorder = func(node *Node) {
			if node == nil {
				return
			}
			inorder(node.left)
			t.iterChan <- node
			inorder(node.right)
		}
		go func() {
			inorder(n)
			close(t.iterChan)
		}()
	case PreOrder:
		var preorder func(node *Node)
		preorder = func(node *Node) {
			if node == nil {
				return
			}
			t.iterChan <- node
			preorder(node.left)
			preorder(node.right)
		}
		go func() {
			preorder(n)
			close(t.iterChan)
		}()
	case PostOrder:
		var postorder func(node *Node)
		postorder = func(node *Node) {
			if node == nil {
				return
			}
			postorder(node.left)
			postorder(node.right)
			t.iterChan <- node
		}
		go func() {
			postorder(n)
			close(t.iterChan)
		}()
	default:
		s := fmt.Sprintf("rbTree has not implemented %s.", order)
		panic(s)
	}
	return t.iterChan

}

func (t *RBTree) insertIter(h *Node, key interface{}, value interface{}) *Node {

	// empty tree
	if h == nil {
		return &Node{color: Red, Key: key, Value: value}
	}

	// setup our own stack and helpers
	var (
		stack     = []*Node{}
		count int = 0
		prior *Node
	)

L:
	for {
		switch t.cmp(h.Key, key) {
		case EQ:
			h.Value = value
			return t.root // no need for rest of the fix code
		case LT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.left
			if h == nil {
				h = &Node{color: Red, Key: key, Value: value}
				prior.left = h
				break L
			}
		case GT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.right
			if h == nil {
				h = &Node{color: Red, Key: key, Value: value}
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

			switch t.cmp(stack[count-1].Key, h.Key) {
			case LT:
				stack[count-1].left = h
			case GT:
				stack[count-1].right = h
			}
			h = stack[count-1]
		}

	}

	return h
}
*/

// Experinmental
/*

func TestSearchRecursive(t *testing.T) {

	tree := New(testCmpInt)
	iters := 10000
	for i := 0; i < iters; i++ {
		tree.Insert(i, i)
	}
	_, ok := tree.searchRecurse(nil)
	if ok {
		t.Errorf("Not minding nil key's")
	}

	tree.Traverse(InOrder, inc(t))
	for i := 0; i < iters; i++ {
		value, ok := tree.searchRecurse(i)
		if !ok {
			t.Errorf("All these values should be present")
		}
		if value != i {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}

	for i := iters; i < iters*2; i++ {
		value, ok := tree.searchRecurse(i)
		if ok {
			t.Errorf("values should not be present")
		}
		if value != nil {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}
}

func TestIterRecurse(t *testing.T) {

	tree := New(testCmpInt)
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	preOrder := []string{"d", "b", "a", "c", "h", "f", "e", "g", "i"}
	postOrder := []string{"a", "c", "b", "e", "g", "f", "i", "h", "d"}
	for i, v := range items {
		tree.Insert(i, v)
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}

	count := 0
	c := tree.initIterRecurse(InOrder)
	for i, n := 0, <-c; n != nil; i, n = i+1, <-c {

		count++
		if items[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, items[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}

	count = 0
	c = tree.initIterRecurse(PreOrder)
	for i, n := 0, <-c; n != nil; i, n = i+1, <-c {

		count++
		if preOrder[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, preOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	c = tree.initIterRecurse(PostOrder)
	for i, n := 0, <-c; n != nil; i, n = i+1, <-c {

		count++
		if postOrder[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, postOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}

	//tree.Traverse(PreOrder, printNode)
	//scs := spew.ConfigState{Indent: "\t"}
	//scs.Dump(tree.root)

}

func BenchmarkSearchRecurse(b *testing.B) {


	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	tree := New(testCmpInt)
	for i := 0; i < searchTotal; i++ {
		a := r.Intn(searchSpace)
		m[a] = a
		tree.Insert(a, a)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := r.Intn(searchSpace)
		tree.searchRecurse(a)
	}
}

func BenchmarkRecurseIterInOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		c := tree.initIterRecurse(InOrder)
		for i, n := 0, <-c; n != nil; i, n = i+1, <-c {
			sum += n.Value.(int)

		}
	}

}

func BenchmarkRecurseIterPreOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		c := tree.initIterRecurse(PreOrder)
		for i, n := 0, <-c; n != nil; i, n = i+1, <-c {
			sum += n.Value.(int)

		}
	}

}

func BenchmarkRecurseIterPostOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		c := tree.initIterRecurse(PostOrder)
		for i, n := 0, <-c; n != nil; i, n = i+1, <-c {
			sum += n.Value.(int)

		}
	}

}
*/

/*

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
*/
