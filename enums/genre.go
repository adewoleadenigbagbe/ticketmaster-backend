package enums

type Genre uint64

const (
	None           Genre = iota + 0
	Action         Genre = 1 << 0
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
)

// TODO:If you ever need to add a new Genre constant value above, dont forget to add the key value to the map
// as there is not way in golang to interate and get the string values of that
func (genre Genre) GetKeyValues() map[Genre]string {
	kv := map[Genre]string{}

	kv[Action] = "Action"
	kv[Adventure] = "Adventure"
	kv[Animation] = "Animation"
	kv[Comedy] = "Comedy"
	kv[Crime] = "Crime"
	kv[Documentary] = "Documentary"
	kv[Drama] = "Drama"
	kv[Family] = "Family"
	kv[Fantasy] = "Fantasy"
	kv[History] = "History"
	kv[Horror] = "Horror"
	kv[Music] = "Music"
	kv[Mystery] = "Mystery"
	kv[Romance] = "Romance"
	kv[ScienceFiction] = "ScienceFiction"
	kv[TVMovie] = "TVMovie"
	kv[Thriller] = "Thriller"
	kv[War] = "War"
	kv[Western] = "Western"

	return kv
}

func (genre Genre) GetValues() []string {
	kv := genre.GetKeyValues()
	values := []string{}
	for _, v := range kv {
		values = append(values, v)
	}

	return values
}

func (genre Genre) GetValue() string {
	kv := genre.GetKeyValues()

	return kv[genre]
}

func (genre Genre) GetFlaggedValues() []string {
	return []string{}
}
