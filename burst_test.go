package gotree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkBurstText(b *testing.B) {

	listContainerMax = 256
	b.StopTimer()
	burst := &BurstTree{}
	content, err := ioutil.ReadFile("testText.txt")
	if err != nil {
		panic("Couldn't read in file to benchmark on")
	}
	data := strings.Fields(string(content))
	fmt.Println(len(data))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range data {
			burst.Insert(exString(e))
		}
		burst.Clear()
	}
}

var _ = spew.Dump

func init() {
	listContainerMax = 1
}

type exByte struct {
	id string
}

func (this exByte) ToBytes() []byte {
	return []byte(this.id)
}

func TestBurstInsert(t *testing.T) {
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
	old = burst.Insert(exByte{"aab"})
	if _, ok := burst.root.(*accessContainer).records['a'].(*accessContainer); !ok {
		t.Errorf("Didn't create new accessContainer")

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

func TestBurstSearch(t *testing.T) {
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

	if check := burst.Search(exByte{"a"}); check == nil {
		t.Errorf("Did not handle single entry in container")
	}
	burst.Insert(exByte{"aa"})
	burst.Insert(exByte{"aab"})
	a := exByte{"a"}
	if check := burst.Search(exByte{"a"}); check == nil || check != a {
		t.Errorf("Did not handle single entry in accessContainer")
	}

}

func TestBurstInsertAndSearchAlot(t *testing.T) {
	burst := &BurstTree{}
	size := 1000
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
		if check == nil {
			t.Errorf("Should have found something")
		}
	}
}

func TestBurstInsertAndSearchRand(t *testing.T) {
	burst := &BurstTree{}
	size := 2000

	data := rand.Perm(size)
	size = 0
	for _, i := range data {
		if i != 0 {
			size++
			s := fmt.Sprintf("%d", i)

			burst.Insert(exByte{s})
		}
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

func TestBurstSize(t *testing.T) {
	//TODO: size increases for insert, and remove
	// mind access containers, bursts, and containers
}

func TestBurstRemove(t *testing.T) {
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

	max := 100
	for i := 1; i < max; i++ {
		burst.Insert(exByte{fmt.Sprintf("%d", i)})
	}
	for i := 1; i < max; i++ {
		s := fmt.Sprintf("%d", i)
		a := exByte{s}
		if check := burst.Remove(exByte{s}); check == nil || a != check {
			t.Errorf("Should Elem: %v be matched element removed", a)
			t.Errorf(" Elem: %v,  Current: %v", check, a)
			//spew.Dump(burst)
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
	//spew.Dump(burst)
}
