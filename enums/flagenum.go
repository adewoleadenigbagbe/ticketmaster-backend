package enums

type FlaggedEnum[T ~uint64] uint64

func (FlaggedEnum[T]) New() FlaggedEnum[T] {
	return FlaggedEnum[T](0)
}

func (f FlaggedEnum[T]) Add(enums ...T) T {
	t := new(T)
	for _, enum := range enums {
		*t |= enum
	}
	return *t
}

func (f FlaggedEnum[T]) Remove(from T, enums ...T) T {
	for _, enum := range enums {
		from &= ^enum
	}
	return from
}

func (f FlaggedEnum[T]) Has(from T, enum T) bool {
	t := new(T)
	return *t&enum != 0
}

func (f FlaggedEnum[T]) And(enum T) bool {
	t := new(T)
	return *t&enum != 0
}
