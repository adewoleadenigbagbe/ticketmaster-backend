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

func (f FlaggedEnum[T]) GetEnumValues() []T {
	values := new([]T)
	var s uint64
	for i := 64; s > 0; i-- {
		if 1<<i > 0 {
			t := T(1 << i)
			*values = append(*values, t)
		}
	}
	return *values
}

func (f FlaggedEnum[T]) GetCount(values []T) int {
	return len(values)
}

func (f FlaggedEnum[T]) GetStringValues(from T, enum T) bool {
	return from&enum != 0
}
