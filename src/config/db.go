package config

import (
	"os"
	"time"
	"y-ogura/echo_lambda_sample/src/logger"
	"y-ogura/echo_lambda_sample/src/utils"
	"y-ogura/echo_lambda_sample/src/utils/aws/secrets_manager"

	"github.com/aws/aws-sdk-go-v2/aws"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg aws.Config) *gorm.DB {
	DBuser, DBpass := getDBCredentials(cfg)

	c := mysqlDriver.Config{
		DBName:               getEnv("DB_NAME", "myapp"),
		User:                 DBuser,
		Passwd:               DBpass,
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

// SECRETS_MANAGERとREGIONが設定されている場合はSecretsManagerから取得
func getDBCredentials(cfg aws.Config) (string, string) {
	DBuser := getEnv("DB_USER", "root")
	DBpass := getEnv("DB_PASS", "password")
	smName := getEnv("SECRETS_MANAGER", "")
	region := getEnv("REGION", "")

	if smName == "" || region == "" {
		return DBuser, DBpass
	}

	secretsManager := secrets_manager.NewSecretsManager(cfg)
	sm, err := secretsManager.GetRDS(smName, region)
	if err == nil {
		DBuser = sm.Username
		DBpass = sm.Password
	}

	return DBuser, DBpass
}
