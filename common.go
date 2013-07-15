// Provides reference based containers using tree or tree like data structures.
// Two types of trees are available based upon the abilities of the to be inserted data.
// A general type that has a total order and a type which can be represented by a byte array.
package gotree

import ()

type ByteTree interface {
	Search(item Byte) (found Byte)
	Insert(item Byte) (found Byte)
	Remove(item Byte) (found Byte)
	Size() int
	Clear()

	IterInit(order TravOrder) Byte
	Next() Byte
	Map(order TravOrder, f ByteIterFunc)
}

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

type ByteIterFunc func(Byte)

// Our possible tree traversal abilities
type TravOrder int

// InOrder items are visted from smallest to largest, while RevOrder visits them from largest to smallest.
// LevelOrder, PreOrder, PostOrder are dependent on the shape and layout of the underlying tree.
// AnyOrder is where the algorithm is chosen for performance reasons.
// RandOrder, items are visted in a uniformly random order.
const (
	InOrder TravOrder = iota
	RevOrder
	PreOrder
	PostOrder
	LevelOrder
	AnyOrder
	RandOrder
)

// implement the String interface for human readable names of traversal abilities
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
	case RevOrder:
		s = "reverse-order traversal"
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
type Interface interface {
	// Compare returns a Balance with respect to the total order between the method calle and its given arguement.
	// Given:
	//
	//      bal = calle.Compare(arg):
	//
	// The result of bal must follow this logic:
	//
	//      if calle < arg {
	//          return  LT
	//      }
	//
	//      if calle > arg {
	//          return  GT
	//      }
	//
	//      if calle == arg {
	//          return  EQ
	//      }
	// Think of Compare as asking the question, "what is the calle's relationship to arg?"
	// Is the calle less than arg? Is it greater than? Or equal to?
	Compare(Interface) Balance
}

type Byte interface {
	ToBytes() []byte
}
