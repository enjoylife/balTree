package gotree

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkBurstText(b *testing.B) {

	containerMax = 256
	b.StopTimer()
	burst := &BurstTree{}
	content, err := ioutil.ReadFile("misc/testText.txt")
	if err != nil {
		panic("Couldn't read in file to benchmark on")
	}
	data := strings.Fields(string(content))
	//fmt.Println(len(data))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range data {
			burst.Insert(exString(e))
		}
		burst.Clear()
	}
}

func TestBurstText(t *testing.T) {

	containerMax = 100
	burst := &BurstTree{}
	content, err := ioutil.ReadFile("misc/testText.txt")
	if err != nil {
		panic("Couldn't read in file to benchmark on")
	}
	data := strings.Fields(string(content))
	m := map[string]bool{}
	for _, e := range data {
		burst.Insert(exString(e))
		m[e] = true
	}
	for _, e := range data {
		//fmt.Println("TestLoop", i)
		x := exString(e)
		if found := burst.Search(exString(e)); found != x {
			t.Errorf("Not Found %v", exString(e))
		}
		m[e] = true
	}

	if len(m) != burst.Size() {
		t.Errorf("Sizes don't match")
	}
	for _, e := range data {
		burst.Remove(exString(e))
		delete(m, e)
	}
	if len(m) != burst.Size() || burst.Size() != 0 {
		t.Errorf("Sizes don't match")
	}
}

var _ = spew.Dump

func init() {
	containerMax = 1
}

type exByte struct {
	id string
}

func (this exByte) ToBytes() []byte {
	return []byte(this.id)
}

func testListContainer(t *testing.T) {
	containerMax = 4
	x := &listContainer{list.New(), nil}
	if old, parent := x.insert([]byte{1}, exByte{"1"}); old != nil || parent != nil {
		t.Errorf("inital insert wrong")
	}
	if old, parent := x.insert([]byte{2}, exByte{"2"}); old != nil || parent != nil {
		t.Errorf("2 element wrong insert")
	}
	if old, parent := x.insert([]byte{3}, exByte{"3"}); old != nil || parent != nil {
		t.Errorf("3 element wrong insert")
	}
	spew.Dump(x)
}
func TestCompactArry(t *testing.T) {
	containerMax = 10
	x := &compactArray{}

	data := rand.Perm(containerMax)
	for _, a := range data {
		s := fmt.Sprintf("%d", a)
		x.insert([]byte{byte(a)}, exByte{s})
	}
	x.sort(InOrder)
	spew.Dump(x)
	for i, a := range x.items {
		s := fmt.Sprintf("%d", i)
		b := exByte{s}

		if a != b {
			t.Errorf("Wrong order")
		}
	}
	x.sort(RevOrder)
	for i, a := range x.items {
		s := fmt.Sprintf("%d", containerMax-i-1)
		b := exByte{s}

		if a != b {
			t.Errorf("Wrong order")
		}
	}

}

func TestBurstInsertPrimary(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}
	var old Byte
	old = burst.Insert(nil)
	if old != nil {
		t.Errorf("Should not accept nil")
	}
	old = burst.Insert(exByte{""})
	if old != nil {
		t.Errorf("Should avoid empty string")
	}
	burst.Insert(exByte{"a"})
	if _, ok := burst.root.(*accessContainer); !ok {
		t.Errorf("Didnt not correctly handle nil root")
	}
	if _, ok := burst.root.(*accessContainer).records['a'].(container); !ok {
		t.Errorf("Didn't create new container")

	}
	if s := burst.Size(); s != 1 {
		t.Errorf("Size isn't proper")
	}
	old = burst.Insert(exByte{"aa"})

	if s := burst.Size(); s != 2 {
		t.Errorf("Size isn't proper")
	}
	old = burst.Insert(exByte{"aab"})
	//spew.Dump(burst)
	if _, ok := burst.root.(*accessContainer).records['a'].(*accessContainer); !ok {
		t.Errorf("Didn't create new accessContainer")

	}

	if s := burst.Size(); s != 3 {
		t.Errorf("Size isn't proper")
	}
	a2 := burst.root.(*accessContainer).records['a'].(*accessContainer).single
	a := exByte{"a"}
	if a2 != a {
		t.Errorf("Didn't add empty string single record")
	}

	if s := burst.Size(); s != 3 {
		t.Errorf("Size isn't proper")
	}
}

func TestBurstSearchPrimary(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}
	if check := burst.Search(nil); check != nil {
		t.Errorf("Should not accept nil")
	}
	if check := burst.Search(exByte{""}); check != nil {
		t.Errorf("Should not accept empty string")
	}
	if check := burst.Search(exByte{"a"}); check != nil {
		t.Errorf("Did not handle nil root")
	}
	burst.Insert(exByte{"a"})

	a := exByte{"a"}
	if check := burst.Search(exByte{"a"}); check == nil || check != a {
		t.Errorf("Did not handle single entry in container")
	}
	burst.Insert(exByte{"aa"})
	burst.Insert(exByte{"aab"})
	if check := burst.Search(exByte{"a"}); check == nil || check != a {
		t.Errorf("Did not handle single entry in accessContainer")
	}

}

func TestBurstInsertSwitch(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}
	size := 10
	for i := 1; i < size+1; i++ {
		s := fmt.Sprintf("%d", i)
		burst.Insert(exByte{s})
	}

	if s := burst.Size(); s != size {
		t.Errorf("Size isn't proper")
	}
	for i := 1; i < size+1; i++ {
		s := fmt.Sprintf("%d", i)
		check := burst.Search(exByte{s})
		x := exByte{s}
		if check != x {
			t.Errorf("Should have found something")
			t.Error(exByte{s})
			t.Error(check)
		}
	}
}

func TestBurstInsertContainerInsert(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}
	size := 1000
	start := 10
	for i := start; i < size+1; i++ {
		s := fmt.Sprintf("%d", i)
		burst.Insert(exByte{s})
	}

	if s := burst.Size(); s != size-start+1 {
		t.Errorf("Size isn't proper", burst.Size())
	}
	for i := start; i < size+1; i++ {
		s := fmt.Sprintf("%d", i)
		check := burst.Search(exByte{s})
		x := exByte{s}
		if check != x {
			t.Errorf("Should have found something")
			t.Error(exByte{s})
			t.Error(check)
		}
	}
}

func TestBurstInsertAndSearchRand(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}
	size := 2000

	data := rand.Perm(size)
	for _, i := range data {
		i++ // avoid 0 aka nil string
		s := fmt.Sprintf("%d", i)

		burst.Insert(exByte{s})
	}
	if s := burst.Size(); s != size {
		t.Errorf("Size isn't proper")
		t.Error("True:", size, " Found:", s)
	}

	if check := burst.root.(*accessContainer).single; check != nil {
		t.Error("Should be nil in root record")
	}

	for _, i := range data {
		if i != 0 {
			s := fmt.Sprintf("%d", i)
			check := burst.Search(exByte{s})
			if check == nil {
				t.Errorf("Should have found something")
				t.Error(s)
				t.Error(exByte{s})
			}
		}
	}
}

func TestBurstRemove(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}

	if check := burst.Remove(nil); check != nil {
		t.Errorf("Should not accept nil")
	}
	if check := burst.Remove(exByte{""}); check != nil {
		t.Errorf("Should not accept empty string")
	}
	if check := burst.Remove(exByte{"a"}); check != nil {
		t.Errorf("Did not handle nil root")
	}

	max := 200
	for i := 1; i < max; i++ {
		burst.Insert(exByte{fmt.Sprintf("%d", i)})
	}
	//spew.Dump(burst)
	for i := 1; i < max; i++ {
		s := fmt.Sprintf("%d", i)
		a := exByte{s}
		if check := burst.Remove(exByte{s}); check == nil || a != check {
			t.Errorf("Should be Elem: %v be matched element removed", a)
			t.Errorf(" Elem: %v,  Current: %v", check, a)
		}
		check := burst.Search(exByte{s})
		if check != nil || check == a {
			t.Errorf("Should not have found something")
		}
	}
	for _, v := range burst.root.(*accessContainer).records {
		if v != nil {
			t.Errorf("Should have empty root")
		}
	}
}

func TestBurstIter(t *testing.T) {
	containerMax = 1
	burst := &BurstTree{}

	max := 200

	data := rand.Perm(max)

	for _, x := range data {
		burst.Insert(exByte{fmt.Sprintf("%d", x+1)})
	}
	prior := []byte{0}
	for i, x := 1, burst.IterInit(InOrder); x != nil; i, x = i+1, burst.Next() {
		if x.ToBytes()[0] < prior[0] {
			t.Errorf("Wrong Order. ")
		}
		if len(x.ToBytes()) > len(prior) {
			if !bytes.Equal(x.ToBytes()[:len(prior)], prior) {
				t.Errorf("Empty byte in wrong order")
			}
		}
		prior = x.ToBytes()
	}
}
