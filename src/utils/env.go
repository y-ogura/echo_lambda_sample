package utils

import "os"

const (
	envLocal = "local"
	envDev   = "dev"
	envStg   = "stg"
	envProd  = "prod"
	envCI    = "ci"
)

// 実行環境がlocalまたはdevならDebug環境
func IsDebugEnv() bool {
	env := os.Getenv("ENV")
	return env == envLocal || env == envDev
}

// 実行環境がlocal環境かを確認
func IsLocal() bool {
	env := os.Getenv("ENV")
	return env == envLocal
}

// 実行環境がDocker上か
func IsDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

// Gitlab CI/CD環境判定
func IsCI() bool {
	return os.Getenv("IS_CI") == envCI
}
