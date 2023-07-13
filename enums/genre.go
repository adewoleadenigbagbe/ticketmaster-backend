package enums

type Genre int

const (
	None           Genre = iota + 0
	Action               = 1
	Adventure            = 1 << 1
	Animation            = 1 << 2
	Comedy               = 1 << 3
	Crime                = 1 << 4
	Documentary          = 1 << 5
	Drama                = 1 << 6
	Family               = 1 << 7
	Fantasy              = 1 << 8
	History              = 1 << 9
	Horror               = 1 << 10
	Music                = 1 << 11
	Mystery              = 1 << 12
	Romance              = 1 << 13
	ScienceFiction       = 1 << 14
	TVMovie              = 1 << 15
	Thriller             = 1 << 16
	War                  = 1 << 17
	Western              = 1 << 18
)

var MovieGenre Genre

func (movieGenre *Genre) AddGenre(genres ...Genre) {
	for _, genre := range genres {
		*movieGenre |= genre
	}

}

func (movieGenre *Genre) RemoveGenre(genres ...Genre) {
	for _, genre := range genres {
		*movieGenre &= ^genre
	}
}

func (movieGenre *Genre) HasGenre(genre Genre) bool {
	return *movieGenre&genre != 0
}
