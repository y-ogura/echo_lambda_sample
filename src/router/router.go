package router

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Init(e *echo.Echo, db *gorm.DB) {
	urlPrefix := os.Getenv("PAGE_URL_PREFIX")
	g := e.Group(urlPrefix)
	g.GET("/health", func(ctx echo.Context) error {
		message := "hello, world!"
		return ctx.String(http.StatusOK, message)
	})
}
