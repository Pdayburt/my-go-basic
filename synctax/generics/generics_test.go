package generics

type List[T any] interface {
	Add(idx int, v T)
	Append(val T)
}

func UserList() {
	//var l List[int]
	//l.Append("jack")
}

func Sum[T Number](vals ...T) T {
	var res T
	for _, val := range vals {
		res += val
	}
	return res
}

type Number interface {
	~int | float32 | float64
}
