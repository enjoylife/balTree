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

// Interface is a one method interface type.
// Types which implement a Compare method can be inserted into a Tree.
// Compare returns a Balance with respect to the total order between the method calle and its given arguement.
// Given:
//
//      bal = calle.Compare(arg):
//
// The result of bal must follow this logic:
//
//      if calle < arg {
//          // ...
//          return // bal == LT
//      }
//
//      if calle > arg {
//          // ...
//          return // bal == GT
//      }
//      //bal == GT
//
//      if calle == arg {
//          // ...
//          return // bal == EQ
//      }
//
//      if (calle can't be compared to arg) {
//          // ...
//          return // bal == NP
//      }
//
// Think of Compare as asking the question, "what is the calle's relationship to arg?"
// Is the calle less than arg? Is it greater than, etc.

// To Insert a type into the RBTree it must implement this one method interface.
type Interface interface {
	Compare(Interface) Balance
}

// This is returned when a Interface's compare method returns a NP case.```
type UncompareableTypeError struct {
	this Interface
	that Interface
}

func (e UncompareableTypeError) Error() string {
	return fmt.Sprintf("gotree: Can not compare %T with the unkown type of %T", e.this, e.that)
}

// Returned when the Interface to be inserted, searched, removed, etc is nil or something uncomparable
type InvalidInterfaceError string

func (e InvalidInterfaceError) Error() string {
	return ("gotree: Can't use nil as item to search for.")
}

type NonexistentElemError string

func (e NonexistentElemError) Error() string {
	return "gotree: Could not find requested Elem."
}
