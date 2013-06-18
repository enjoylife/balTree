package gotree

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

const (
	searchTotal = 1000000
	searchSpace = searchTotal / 2
	iters       = 10000
)

var trees = []Tree{&RBTree{}, &SplayTree{}}

type exInt int

func (this exInt) Compare(b Interface) Balance {
	var out Balance
	switch that := b.(type) {
	case exInt:
		switch result := int(this - that); {
		case result > 0:
			out = GT
		case result < 0:
			out = LT
		case result == 0:
			out = EQ
		}
	case exStruct:
		switch result := int(this) - that.M; {
		case result > 0:
			out = GT
		case result < 0:
			out = GT
		case result == 0:
			out = EQ
		}

	}
	return out
}

type exString string

func (this exString) ToBytes() []byte {
	return []byte(this)
}

func (this exString) Compare(b Interface) Balance {
	var out Balance
	switch that := b.(type) {
	case exString:
		a := string(this)
		b := string(that)
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
			out = GT
		case result < 0:
			out = LT
		case result == 0:
			out = EQ
		}
	case exInt:
		a, _ := strconv.Atoi(string(this))
		switch result := a - int(that); {
		case result > 0:
			out = GT
		case result < 0:
			out = LT
		case result == 0:
			out = EQ
		}

	}
	return out
}

type exStruct struct {
	M int
	S string
}

func (this exStruct) Compare(b Interface) Balance {
	var out Balance
	switch that := b.(type) {
	case exStruct:
		switch result := int(this.M - that.M); {
		case result > 0:
			out = GT
		case result < 0:
			out = LT
		case result == 0:
			out = EQ
		}
	case exInt:
		switch result := this.M - int(that); {
		case result > 0:
			out = GT
		case result < 0:
			out = LT
		case result == 0:
			out = EQ
		}
	}
	return out
}

/* Helpers for tree traversal and testing tree properties */
func printNode(n Interface) {
	fmt.Println("Elem:", n)
}

// Binary tree property tests
func inc(t *testing.T) func(n Interface) {
	var prior int = -1
	return func(n Interface) {
		if prior < int(n.(exInt)) {
			//fmt.Println("VALUE: ", value.(int))
			prior = int(n.(exInt))
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, n)
		}
	}
}

func TestMax(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		if check := tree.Max(); check != nil {
			t.Errorf("Not minding nill tree Got %T", check)
		}

		for i := 0; i < iters; i++ {
			tree.Insert(exInt(i))
			if tree.Max() != exInt(i) {
				t.Errorf("Max not updateing")
			}
		}
		for i := iters; i > 0; i-- {
			tree.Remove(exInt(i))
			if tree.Max() != exInt(i-1) {
				t.Errorf("Max not updateing")
			}
		}
	}
}
func TestMin(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		var old Interface

		if check := tree.Min(); check != nil {
			t.Errorf("Not minding nill tree Got %T", check)
		}
		for i := iters; i >= 0; i-- {
			tree.Insert(exInt(i))
			if tree.Min() != exInt(i) {
				t.Errorf("Min not updateing")
			}
		}
		for i := 0; i < iters; i++ {

			old = tree.Remove(exInt(i))
			if old == nil {
				t.Errorf("Not giving back old value")
			}
			if tree.Min() != exInt(i+1) {
				t.Errorf("Min not updateing")
			}
		}

		tree.Remove(exInt(iters))
		if tree.Min() != nil {
			t.Errorf("Min not working")
		}
	}
}

func TestSize(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		for i := 0; i <= iters; i++ {
			tree.Insert(exInt(i))
			if tree.Size() != i+1 {
				t.Errorf("Size not correctly updateing")
			}
		}

		for i := 0; i < iters; i++ {
			tree.Remove(exInt(i))
			if tree.Size() != (iters - i) {
				t.Errorf("Size on remove not working")
			}

		}
	}
}

func TestInsert(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		var old Interface
		items := []exStruct{exStruct{0, "0"},
			exStruct{2, "2"}, exStruct{2, "3"}}

		old = tree.Insert(items[0])
		if old != nil {
			t.Errorf("First check on input is messed!")
		}

		var fake Interface
		old = tree.Insert(fake)
		if old != nil {
			t.Errorf("Should not accept nil")
		}
		tree.Insert(items[1])
		old = tree.Insert(items[2])
		if tree.Size() != 2 {
			t.Errorf("Problems tracking Size")
		}
		if old != items[1] {
			t.Errorf("old input is messed!")
		}
		old = tree.Insert(items[1])
		if tree.Size() != 2 {
			t.Errorf("Problems tracking Size")
		}
		if old != items[2] {
			t.Errorf("old input is messed!")
		}
	}
}

// assumes inorder implemented
func TestRandomInsert(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		r := rand.New(rand.NewSource(int64(5)))

		for i := 0; i < iters; i++ {
			item := r.Int()
			tree.Insert(exInt(item))
			tree.Map(InOrder, inc(t))
		}
	}
}

func TestClear(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		r := rand.New(rand.NewSource(int64(5)))

		for i := 0; i < iters; i++ {
			item := r.Int()
			tree.Insert(exInt(item))
		}
		prevHeight := tree.Height()
		prevSize := tree.Size()
		tree.Clear()
		if tree.Height() == prevHeight {
			t.Errorf("Height didn't reset")
		}
		if tree.Size() == prevSize {
			t.Errorf("Size didn't reset")
		}
	}
}

func TestSearch(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		x := tree.Search(nil)
		if x != nil {
			t.Errorf("Not minding empty tree")
		}
		for i := 0; i < iters; i++ {
			tree.Insert(exInt(i))
		}
		x = tree.Search(nil)
		if x != nil {
			t.Errorf("Not minding nil key's")
		}

		for i := 0; i < iters; i++ {
			value := tree.Search(exInt(i))
			if int(value.(exInt)) != i {
				t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
			}
		}

		for i := iters; i < iters+1000; i++ {
			value := tree.Search(exInt(i))
			if value != nil {
				t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	for _, v := range trees {
		tree := v
		tree.Clear()

		c := tree.Remove(exInt(0))
		if c != nil {
			t.Errorf("Not respecting empty tree.")
		}
		items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
		for _, v := range items {
			tree.Insert(exString(v))
		}
		for k, v := range items {
			check := tree.Remove(exString(v))
			if check != exString(v) {
				fmt.Println(check)
				t.Errorf("Not getting removed item back")
			}
			check = tree.Search(exString(v))
			if check != nil {
				t.Errorf("Didn't really remove")
			}

			for i, n := k+1, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {
				if items[i] != string(n.(exString)) {
					t.Errorf("Other elems deleted Got:%s, Exp: %s", n, items[i])
				}
			}
		}
	}
}

func TestRandomRemove(t *testing.T) {

	for _, v := range trees {
		tree := v
		tree.Clear()
		r := rand.New(rand.NewSource(int64(5)))
		m := make(map[int]int)
		for i := 0; i < iters; i++ {
			a := r.Intn(searchSpace)
			m[a] = a
			tree.Insert(exInt(a))
		}
		for _, value := range m {
			tree.Remove(exInt(value))

			if check := tree.Search(exInt(value)); check != nil {
				t.Errorf("Didn't really remove")
			}
		}
	}
}

func TestIterInorder(t *testing.T) {

	for _, v := range trees {
		tree := v
		tree.Clear()
		items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
		for _, v := range items {
			tree.Insert(exString(v))
		}
		var check Interface
		if check = tree.Next(); check != nil {
			t.Errorf("Didn't avoid a non intialized next call")
		}

		count := 0
		for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {

			count++
			if items[i] != string(n.(exString)) {
				t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
			}

		}
		if count != len(items) {
			t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
		}

		if check = tree.Next(); check != nil {
			t.Errorf("Didn't avoid a non intialized next call")
		}

		count = 0
		exit := 3
		for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {
			count++
			if items[i] != string(n.(exString)) {
				t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
			}
			if i == exit {
				break
			}
		}
		// restart same iterator
		for i, n := exit+1, tree.Next(); n != nil; i, n = i+1, tree.Next() {
			count++
			if items[i] != string(n.(exString)) {
				t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
			}
		}
	}

}

func TestIterRevorder(t *testing.T) {

	for _, v := range trees {
		tree := v
		tree.Clear()
		//items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
		items := []string{"i", "h", "g", "f", "e", "d", "c", "b", "a"}
		for _, v := range items {
			tree.Insert(exString(v))
		}
		var check Interface
		if check = tree.Next(); check != nil {
			t.Errorf("Didn't avoid a non intialized next call")
		}

		count := 0
		for i, n := 0, tree.IterInit(RevOrder); n != nil; i, n = i+1, tree.Next() {

			count++
			if items[i] != string(n.(exString)) {
				t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
			}

		}
		if count != len(items) {
			t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
		}

		if check = tree.Next(); check != nil {
			t.Errorf("Didn't avoid a non intialized next call")
		}

		count = 0
		exit := 3
		for i, n := 0, tree.IterInit(RevOrder); n != nil; i, n = i+1, tree.Next() {
			count++
			if items[i] != string(n.(exString)) {
				t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
			}
			if i == exit {
				break
			}
		}
		// restart same iterator
		for i, n := exit+1, tree.Next(); n != nil; i, n = i+1, tree.Next() {
			count++
			if items[i] != string(n.(exString)) {
				t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
			}
		}
	}

}

func benchSearch(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {
		b.StopTimer()
		r := rand.New(rand.NewSource(int64(1234)))
		m := make(map[int]int)
		for i := 0; i < searchTotal; i++ {
			a := r.Intn(searchSpace)
			m[a] = a
			tree.Insert(exInt(a))
			//tree.Insert(a, a)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			a := r.Intn(searchSpace)
			tree.Search(exInt(a))
		}
	}
}

func benchSearchDense(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {
		b.StopTimer()
		r := rand.New(rand.NewSource(int64(1234)))
		m := make(map[int]exInt)
		for i := 0; i < b.N; i++ {
			a := r.Intn(searchSpace)
			m[i] = exInt(a)
			tree.Insert(exInt(a))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			tree.Search(m[i])
		}
	}
}

func benchInsert(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.Insert(exInt(b.N - i))
		}
	}
}

func benchRandomInsert(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.Insert(exInt(rand.Int()))
		}
	}
}

func benchRemove(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {

		b.StopTimer()
		for i := 0; i < b.N; i++ {
			tree.Insert(exInt(b.N - i))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			tree.Remove(exInt(b.N - i))
		}
	}
}

func benchRandomRemove(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {
		b.StopTimer()
		r := rand.New(rand.NewSource(int64(1234)))
		for i := 0; i < searchTotal; i++ {
			a := r.Intn(searchSpace)
			tree.Insert(exInt(a))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			a := r.Intn(searchSpace)
			tree.Remove(exInt(a))
		}

	}
}

func benchIterInorder(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {

		b.StopTimer()
		r := rand.New(rand.NewSource(int64(5)))
		for i := 0; i < 1000; i++ {
			tree.Insert(exInt(r.Int()))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sum := 0

			for n := tree.IterInit(InOrder); n != nil; n = tree.Next() {
				i := int(n.(exInt))
				sum += i
			}
		}
	}
}

func recurse() func(n Interface) {
	var sum int = 0
	return func(n Interface) {
		sum += int(n.(exInt))
	}
}

func benchMap(tree Tree) func(b *testing.B) {
	tree.Clear()
	return func(b *testing.B) {

		b.StopTimer()
		r := rand.New(rand.NewSource(int64(5)))
		for i := 0; i < 1000; i++ {
			tree.Insert(exInt(r.Int()))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			f := recurse()
			tree.Map(InOrder, f)

		}
	}
}

func benchText(tree Tree) func(b *testing.B) {

	tree.Clear()
	return func(b *testing.B) {
		b.StopTimer()
		content, err := ioutil.ReadFile("testText.txt")
		if err != nil {
			panic("Couldn't read in file to benchmark on")
		}
		data := strings.Fields(string(content))
		fmt.Println(len(data))
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for _, e := range data {
				tree.Insert(exString(e))
			}
			fmt.Println("Tree Size", tree.Size())
			tree.Clear()
		}
	}
}

func TestEfficiency(t *testing.T) {
	var result testing.BenchmarkResult
	fmt.Println("Starting Efficiency Functions")

	fmt.Println("\nExample text file insert")
	for _, v := range trees {
		result = testing.Benchmark(benchText(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

	fmt.Println("\nRandom Search")
	for _, v := range trees {
		result = testing.Benchmark(benchSearch(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

	fmt.Println("\nRandom Search Dense")
	for _, v := range trees {
		result = testing.Benchmark(benchSearchDense(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

	fmt.Println("\nLinear Insert")
	for _, v := range trees {
		result = testing.Benchmark(benchInsert(v))
		fmt.Printf("%-30T %s\n", v, result)
	}
	fmt.Println("\nRandom Insert")
	for _, v := range trees {
		result = testing.Benchmark(benchRandomInsert(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

	fmt.Println("\nLinear Remove")
	for _, v := range trees {
		result = testing.Benchmark(benchRemove(v))
		fmt.Printf("%-30T %s\n", v, result)
	}
	fmt.Println("\nRandom Remove")
	for _, v := range trees {
		result = testing.Benchmark(benchRandomRemove(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

	fmt.Println("\nInorder traverse")
	for _, v := range trees {
		result = testing.Benchmark(benchIterInorder(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

	fmt.Println("\nInorder map")
	for _, v := range trees {
		result = testing.Benchmark(benchMap(v))
		fmt.Printf("%-30T %s\n", v, result)
	}

}
