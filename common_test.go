package gotree

import (
	"fmt"
	"strconv"
)

const (
	searchTotal = 1000000
	searchSpace = searchTotal / 2
	iters       = 10000
)

type exInt int

func (this exInt) Compare(b Interface) Balance {
	switch that := b.(type) {
	case exInt:
		switch result := int(this - that); {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}
	case exStruct:
		switch result := int(this) - that.M; {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}

	default:
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}

}

type exString string

func (this exString) Compare(b Interface) Balance {
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
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}
	case exInt:
		a, _ := strconv.Atoi(string(this))
		switch result := a - int(that); {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}

	default:
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}
}

type exStruct struct {
	M int
	S string
}

func (this exStruct) Compare(b Interface) Balance {
	switch that := b.(type) {
	case exStruct:
		switch result := int(this.M - that.M); {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}
	case exInt:
		switch result := this.M - int(that); {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}
	default:
		return NP
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}
}
