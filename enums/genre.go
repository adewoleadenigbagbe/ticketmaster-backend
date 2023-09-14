package enums

type Genre uint64

const (
	None           Genre = iota + 0
	Action         Genre = 1
	Adventure      Genre = 1 << 1
	Animation      Genre = 1 << 2
	Comedy         Genre = 1 << 3
	Crime          Genre = 1 << 4
	Documentary    Genre = 1 << 5
	Drama          Genre = 1 << 6
	Family         Genre = 1 << 7
	Fantasy        Genre = 1 << 8
	History        Genre = 1 << 9
	Horror         Genre = 1 << 10
	Music          Genre = 1 << 11
	Mystery        Genre = 1 << 12
	Romance        Genre = 1 << 13
	ScienceFiction Genre = 1 << 14
	TVMovie        Genre = 1 << 15
	Thriller       Genre = 1 << 16
	War            Genre = 1 << 17
	Western        Genre = 1 << 18
	Western2        Genre = 1 << 63
)
