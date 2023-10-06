package main

func Sum[T Number](vals ...T) T {
	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

type Number interface {
	~int | int64 | float64
}

type Integer int
