package gorm

import (
	"os"
	"time"

	"github.com/hhk7734/zapx.go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLiteGormDBForTest() (*gorm.DB, func()) {
	tempDir, _ := os.MkdirTemp(os.TempDir(), "test-*")
	tempFile, _ := os.CreateTemp(tempDir, "*.db")

	db, _ := gorm.Open(sqlite.Open(tempFile.Name()), &gorm.Config{
		Logger: &zapx.GormLogger{
			Config: logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				Colorful:                  false,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      false,
				LogLevel:                  logger.Warn,
			},
		},
	})

	return db, func() {
		defer os.RemoveAll(tempDir)
	}
}
