package gotree

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

var _ = fmt.Println

func testCmp(a interface{}, b interface{}) int {
	return a.(int) - b.(int)
}

func inc(t *testing.T) func(key interface{}, value interface{}) {
	var prior int = -1
	return func(key interface{}, value interface{}) {
		if prior < value.(int) {
			//fmt.Printf("Prior: %d, Current: %d\n", prior, value.(int))
			prior = value.(int)
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, value.(int))
		}
	}
}

func TestInsertAndDelete(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	eq := func(a interface{}, b interface{}) {
		if a.(int) != b.(int) {
			t.Errorf("%d is not equal to %d", a, b)
		}

	}
	tr := func(a bool) {
		if !a {
			t.Errorf("SHould be true")
		}
	}
	f := func(a bool) {
		if a {
			t.Errorf("SHould be false")
		}
	}
	fn := func(a interface{}) {
		if a == nil {
			t.Errorf("Should not be nil")
		}
	}
	tn := func(a interface{}) {
		if a != nil {
			t.Errorf("Should  be nil")
		}
	}
	_ = eq
	_ = tn
	_ = fn
	_ = f
	_ = tr

	r := rand.New(rand.NewSource(int64(5)))
	mem := make(map[int]int)
	tree := New(testCmp)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		item := r.Int()
		which := r.NormFloat64()
		if which > 0 {
			old, ok := mem[item]
			if ok {
				// we have already inputed the value
				check, _ := tree.Insert(item, item)
				fn(check)
				eq(check, old)
			} else {
				mem[item] = item
				// not in yet
				check, _ := tree.Insert(item, item)
				tn(check)
			}
		} else {
			old, ok := mem[item]
			if ok {
				mem[item] = 0
				check, ok := tree.Remove(item)
				tr(ok)
				fn(check)
				eq(check, old)
			} else {
				// havent inserted yet
				check, ok := tree.Remove(item)
				f(ok)
				tn(check)
			}
		}
		if i%200000 == 0 {
			fmt.Println("Iter", i)
		}
	}
	tree.checkBalance(inc(t), t)
	tree.Traverse(inc(t))

	//tree.Print()
	fmt.Printf("Height: %+v\n", tree.Height)
	fmt.Printf("Size %+v\n", tree.Size)
	fmt.Printf("Space %+v\n", tree.Space("KiB"))

	//tree.Print()
	//fmt.Printf("%+v", tree)
}

func (t *AvlTree) checkBalance(f IterFunc, t2 *testing.T) {

	var node *avlNode = t.first
	if t.root == nil {
		return
	}
	skew := make(map[int]int)
	for {
		f(node.key, node.item)
		switch node.balance {
		case -1, 1, 0:
		default:
			skew[node.balance]++
		}
		if node == t.last {
			if len(skew) != 0 {
				t2.Errorf("Should not have any weird balances")
			}
			break
		}
		node = node.Next()
	}
}

func BenchmarkInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmp)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		old, ok := tree.Insert(r.Int(), r.Int())
		if !ok {
			_ = old
		}
	}

}

func BenchmarkDelete(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmp)
	for i := 0; i < b.N; i++ {
		old, ok := tree.Insert(r.Int(), r.Int())
		if !ok {
			_ = old
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		old, ok := tree.Remove(r.Int())
		if !ok {
			_ = old
		}
	}
}
