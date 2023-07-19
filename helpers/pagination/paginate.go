package paginate

import (
	"gorm.io/gorm"
)

func Paginate(page, pagelength int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pagelength
		return db.Offset(offset).Limit(pagelength)
	}
}

func GetEntityCount(db *gorm.DB, entity interface{}, count *int64) {
	db.Find(entity).Count(count)
}

func FilterFields(entity interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(entity)
	}
}

func OrderBy(sortOrder string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(sortOrder)
	}
}
