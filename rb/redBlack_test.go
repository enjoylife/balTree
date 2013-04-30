package gorbtree

import (
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"gotree"
	"math/rand"
	"testing"
)

var _ = fmt.Printf

func testCmp(a interface{}, b interface{}) gotree.Direction {
	switch result := (a.(int) - b.(int)); {
	case result > 0:
		return gotree.LT
	case result < 0:
		return gotree.GT
	case result == 0:
		return gotree.EQ
	default:
		panic("Invalid Compare function Result")
	}
}

/* Helpers for tree traversal and testing tree properties */

func printNode(key interface{}, value interface{}) {
	fmt.Println("VALUE: ", value.(int))
}

func isBalanced(t *RbTree) bool {
	if t == nil {
		return true
	}
	var black int // number of black links on path from root to min
	for x := t.root; x != nil; x = x.left {
		if x.color == Black {
			black++
		}
	}
	fmt.Println("Black count", black)
	return nodeIsBalanced(t.root, black)
}

func nodeIsBalanced(n *rbNode, black int) bool {
	if n == nil && black == 0 {
		return true
	} else if n == nil && black != 0 {
		return false
	}
	if n.color == Black {
		black--
	}
	return nodeIsBalanced(n.left, black) && nodeIsBalanced(n.right, black)
}

func inc(t *testing.T) func(key interface{}, value interface{}) {
	var prior int = -1
	return func(key interface{}, value interface{}) {
		if prior < value.(int) {
			//fmt.Println("VALUE: ", value.(int))
			prior = value.(int)
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, value.(int))
		}
	}
}

func TestInsert(t *testing.T) {

	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmp)
	iters := 10
	for i := 0; i < iters; i++ {
		item := r.Int()
		item = i
		tree.Insert(item, item)
	}
	fmt.Println(isBalanced(tree))
	fmt.Println(tree.Height)
	order := tree.TraverseTest(inc(t))
	fmt.Printf("Len: %d \n", len(order))
	for i := 0; i < len(order); i++ {
		//fmt.Println(order[i].value.(int))
	}

	var black int // number of black links on path from root to min
	for x := tree.root; x != nil; x = x.left {
		if x.color == Black {
			black++
		}
	}
	if black != tree.Height {
		t.Errorf("Height is not correct got %d, should be", tree.Height, black)
	}

	tree.Traverse(gotree.PreOrder, printNode)
}

func (t *RbTree) TraverseTest(f gotree.IterFunc) []*rbNode {
	node := t.root
	order := []*rbNode{nil}
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			order = append(order, node)
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
	return order
}
func BenchmarkMapInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[r.Int()] = r.Int()
	}

}

func BenchmarkInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmp)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(r.Int(), r.Int())
	}

}
