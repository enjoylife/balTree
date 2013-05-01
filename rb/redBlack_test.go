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

func testCmpString(c interface{}, d interface{}) gotree.Direction {
	a := c.(string)
	b := d.(string)
	min := len(b)
	if len(a) < len(b) {
		min = len(a)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(a[i]) - int(b[i])
	}
	if diff == 0 {
		diff = len(a) - len(b)
	}

	switch result := diff; {
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
	fmt.Println("VALUE: ", value)
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

/*func accumulate(t *testing.T,store, []string) func(key interface{}, value interface{}){
    return func( key interface{}, value interface{}){
    }
}*/

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
	var black int
	for x := tree.root; x != nil; x = x.left {
		if x.color == Black {
			black++
		}
	}
	if black != tree.Height {
		t.Errorf("Height is not correct got %d, should be", tree.Height, black)
	}

	tree.Traverse(gotree.InOrder, inc(t))
}

func TestTraversal(t *testing.T) {
	tree := New(testCmpString)
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(v, v)
	}
	tree.Traverse(gotree.InOrder, printNode)

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
