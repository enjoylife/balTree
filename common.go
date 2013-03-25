package gotree

type Tree interface {
	Search(interface{}) (interface{}, bool)
	Insert(interface{}) interface{}
	Delete(interface{}) interface{}
	Traverse() interface{}
	Compare(interface{}, interface{}) int
}
