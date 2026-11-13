package config

import (
	"os"
	"time"
	"y-ogura/echo_lambda_sample/src/logger"
	"y-ogura/echo_lambda_sample/src/utils"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	c := mysqlDriver.Config{
		DBName:               getEnv("DB_NAME", "myapp"),
		User:                 getEnv("DB_USER", "root"),
		Passwd:               getEnv("DB_PASS", "password"),
		Addr:                 getEnv("DB_HOST", "127.0.0.1") + ":" + getEnv("DB_PORT", "3306"),
		Net:                  "tcp",
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  time.UTC,
		AllowNativePasswords: true,
	}

	db, err := gorm.Open(mysql.Open(c.FormatDSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		logger.Panic("Failed to connect to database", err)
	}

	if utils.IsDebugEnv() {
		db = db.Debug()
	}

	return db
}

func Migrate(db *gorm.DB) {
	logger.Log.Info("Start auto migrate")
	// モデルをマイグレーション（初期DB作成）
	err := db.AutoMigrate()
	if err != nil {
		logger.Panic("Failed to migrate database", err)
	}
}

// 環境変数取得（なければデフォルト）
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
