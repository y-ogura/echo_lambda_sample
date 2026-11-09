package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// 標準ライブラリのlog.Panicfと同様の動作
func Panic(msg string, err error, keyvals ...any) {
	if err != nil {
		keyvals = append([]any{"error", err}, keyvals...)
		Log.Error(msg, keyvals...)
		panic(err)
	}
}
