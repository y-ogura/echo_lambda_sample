package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Init(e *echo.Echo, db *gorm.DB) {
	e.GET("/health", func(ctx echo.Context) error {
		message := "hello, world!"
		return ctx.String(http.StatusOK, message)
	})
}
