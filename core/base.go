package core

import "gorm.io/gorm"

//BaseApp implements core.App and defines the structure of the whole application
type BaseApp struct {
	db *gorm.DB
}

func NewBaseApp()*BaseApp{
	app := &BaseApp{

	}
	return app
}

func (app *BaseApp) DB() *gorm.DB {
	return app.db
}