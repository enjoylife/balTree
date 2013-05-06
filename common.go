package gotree

// Our possible tree traversal ablitites
type TravOrder int

const (
	PreOrder TravOrder = iota
	PostOrder
	InOrder
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

/*
Needed function to determine insertion order, and compare our stored types equality. Param a and param b are keys of nodes within a tree.
We compare the first param to the second, so if first param is bigger then GT, if equal EQ, etc.

EX:
    func testCmp(a interface{}, b interface{}) gotree.Balance {
        // we assume we are handling int's
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
*/
type CompareFunc func(interface{}, interface{}) Balance

// A Comparable is a type that can be inserted into a Tree or used as a range
// or equality query on the tree,
type Comparable interface {
	Compare(Comparable) Balance
}

/*
Our function we can give to our iterators to work with our stored types.
EX:
    func printNode(key interface{}, value interface{}) {
        fmt.Printf("keyType: %T, valueType: %T \n", key, value)
    }
*/
type IterFunc func(*Node)

type TreeNode interface {
	Key() interface{}
	Value() interface{}
	leftChild() *TreeNode
	lightChild() *TreeNode
}
