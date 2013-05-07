package gotree

import (
	"fmt"
)

// Our possible tree traversal ablitites
type TravOrder int

const (
	InOrder    TravOrder = iota // a,b,c,d,e,f,g,h,i
	PreOrder                    // d,b,a,c,h,f,e,g,i
	PostOrder                   // a,c,b,e,g,f,i,h,d
	LevelOrder                  // dependent on tree state
)

// implement the String interface for human readable names of traversal ablities
// used for debug and error reporting
func (order TravOrder) String() string {
	var s string
	switch order {
	case PreOrder:
		s = "pre-order traversal"
	case PostOrder:
		s = "post-order traversal"
	case InOrder:
		s = "in-order traversal"
	case LevelOrder:
		s = "level-order traversal"
	default:
		s = "unkown traversal"
	}
	return s
}

// Possible directions our path down the tree may take
type Balance int

const (
	GT Balance = iota
	EQ
	LT
	NP
)

// human readable representation of Balance values
// used for debug and error reporting
func (d Balance) String() string {
	var s string
	switch d {
	case GT:
		s = "greater than"
	case EQ:
		s = "equal to"
	case LT:
		s = "less than"
	case NP:
		s = "not possible"
	}
	return s
}

// Comparer is a one method interface type.
// Types which implement a Compare method can be inserted into a Tree.
// Compare returns a Balance with respect to the total order between the method calle and its given
// arguement.
// Given:
//
//      bal = calle.Compare(arg):
//
// The result of bal should follow this logic:
//
//      if calle < arg {
//          bal == LT
//      }
//      if calle > arg {
//          bal == GT
//      }
//      if calle == arg {
//          bal == EQ
//      }
//      if (calle can't be compared to arg) {
//          bal == NP
//      }
// Think of Compare as asking the question, "what is the calle's relationship to arg?"
type Comparer interface {
	Compare(Comparer) Balance
}

type UncompareableTypeError struct {
	this Comparer
	that Comparer
}

func (e UncompareableTypeError) Error() string {
	return fmt.Sprintf("gotree: Can not compare %T with the unkown type of %T", e.this, e.that)
}

type InvalidComparerError string

func (e InvalidComparerError) Error() string {
	return ("gotree: Can't use nil as item to search for.")
}

type NonexistentElemError string

func (e NonexistentElemError) Error() string {
	return "gotree: Remove could not find Elem."
}

/*
Our function we can give to our iterators to work with our stored types.
EX:
    func printNode(n *Node}) {
        fmt.Printf("ElementType: %T, ElementValue: %v\n", n.Elem,n.Elem)
    }
*/

type IterFunc func(*Node)
