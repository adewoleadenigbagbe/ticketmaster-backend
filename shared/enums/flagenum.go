package enums

type FlaggedEnum[T ~uint64] uint64

func (FlaggedEnum[T]) New() FlaggedEnum[T] {
	return FlaggedEnum[T](0)
}

func (f FlaggedEnum[T]) Add(values ...T) T {
	t := new(T)
	for _, value := range values {
		*t |= value
	}
	return *t
}

func (f FlaggedEnum[T]) Remove(source T, values ...T) T {
	for _, value := range values {
		source &= ^value
	}
	return source
}

func (f FlaggedEnum[T]) Has(source T, value T) bool {
	return source&value != 0
}

func (f FlaggedEnum[T]) IsEqualTo(val1 T, val2 T) bool {
	return val1 == val2
}

func (f FlaggedEnum[T]) IsGreaterThan(val1 T, val2 T) bool {
	return val1 > val2
}

func (f FlaggedEnum[T]) IsLesserThan(val1 T, val2 T) bool {
	return val1 < val2
}

func (f FlaggedEnum[T]) TypePossibleValues() []T {
	values := new([]T)
	var s uint64
	for s = 0; s < 64; s++ {
		if uint64(1<<s) > 0 {
			t := T(1 << s)
			*values = append(*values, t)
		}
	}
	return *values
}

func (f FlaggedEnum[T]) GetType() T {
	return *new(T)
}

func (f FlaggedEnum[T]) Max(values []T) T {
	return *new(T)
}

func (f FlaggedEnum[T]) Min(values []T) T {
	return *new(T)
}



