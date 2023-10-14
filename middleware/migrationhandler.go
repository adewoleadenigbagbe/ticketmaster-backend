package middlewares

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

const (
	RootFolderPath = "ticketmaster-backend"
	//If migration folder is moved the path needs to change too
	TargetFolderPath = "database\\migrations"
)

var ()

type MigrationChanges struct {
	App *core.BaseApp
}

// Loop through all the files getting the timestamp of each and check it match the latest migration row in _migration tables
func (mc MigrationChanges) CheckMigrationCompatibility(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != "GET" {
			if mc.App.IsMigrationChecked {
				data := mc.getRecentMigrationData()
				if !reflect.ValueOf(data).IsZero() {
					latestTimestamp := mc.getRecentMigrationFile()
					if latestTimestamp != data.VersionId {
						log.Fatal("Database model has changed. Please pull the recent migration changes")
					}
				}
				mc.App.IsMigrationChecked = true
			}
		}
		return nil
	}
}

func (mc MigrationChanges) getRecentMigrationFile() int64 {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	index := strings.Index(currentWorkingDirectory, RootFolderPath)
	if index == -1 {
		log.Fatal("Target folder not found")
	}

	path := filepath.Join(currentWorkingDirectory[:index], RootFolderPath, TargetFolderPath, "\\*.sql")
	files, err := filepath.Glob(path)

	if err != nil {
		log.Fatalln(err)
	}

	var maxTimeStamp int64 = -1
	for _, file := range files {
		if filepath.Ext(file) == ".sql" {
			baseFileName := filepath.Base(file)
			idx := strings.Index(baseFileName, "_")
			if idx != -1 {
				timestamp, _ := strconv.ParseInt(baseFileName[0:idx], 10, 64)
				if timestamp != 0 && timestamp > maxTimeStamp {
					maxTimeStamp = timestamp
				}
			}
		}
	}

	return maxTimeStamp
}

func (mc MigrationChanges) getRecentMigrationData() MigrationModel {
	model := &MigrationModel{}
	mc.App.DB.Table("_migrations").
		Where("_migrations.is_applied = ?", 1).
		Order("id DESC").
		First(model)

	return *model
}

type MigrationModel struct {
	Id        int64
	VersionId int64
	IsApplied bool
}
