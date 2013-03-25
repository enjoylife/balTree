package gotree

import (
	"fmt"
	"math/rand"
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

	r := rand.New(rand.NewSource(int64(5)))
	mem := make(map[int]int)
	tree := New(testCmp)
	iterations := 10000000

	for i := 0; i < iterations; i++ {
		item := r.Int()
		which := r.NormFloat64()
		if which < 0 {
			old, ok := mem[item]
			if ok {
				// we have already inputed the value
				check := tree.Insert(item, item)
				fn(check)
				eq(check, old)
			} else {
				mem[item] = item
				// not in yet
				check := tree.Insert(item, item)
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
	tree.Traverse(inc(t))

	//tree.Print()
	fmt.Printf("%+v", tree.Height)

	//tree.Print()
	//fmt.Printf("%+v", tree)
}
