package gotree

import ()

type Tree interface {
	Search(item Interface) (found Interface)
	Insert(item Interface) (old Interface)
	Remove(item Interface) (old Interface)

	Clear()
	Size() int
	Height() int

	Min() Interface
	Max() Interface

	IterInit(order TravOrder) Interface
	Next() Interface
	Map(order TravOrder, f IterFunc)
}

// IterFunc is function we can give to our iterators to work with our stored types.
// EX:
//     func printRBNode(n *RBNode}) {
//         fmt.Printf("ElementType: %T, ElementValue: %v\n", n.Elem,n.Elem)
//     }
type IterFunc func(Interface)

// Our possible tree traversal ablitites
type TravOrder int

const (
	InOrder TravOrder = iota
	PreOrder
	PostOrder
	LevelOrder
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
	}
	return s
}

// Possible directions our path down the tree may take
type Balance int

const (
	GT Balance = iota
	EQ
	LT
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
//
// Think of Compare as asking the question, "what is the calle's relationship to arg?"
// Is the calle less than arg? Is it greater than, etc.

// To Insert a type into the RBTree it must implement this one method interface.
type Interface interface {
	Compare(Interface) Balance
}
type Byte interface {
	Interface
	ToBytes() []byte
}
