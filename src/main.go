package main

import (
	"context"
	"os"
	"strings"
	"y-ogura/echo_lambda_sample/src/config"
	"y-ogura/echo_lambda_sample/src/logger"
	"y-ogura/echo_lambda_sample/src/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func main() {
	logger.Log.Info("Start Lambda server")
	godotenv.Load()
	db := config.InitDB() // DB初期化
	e := echo.New()

	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(originsEnv, ",")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Refresh-Token",
			"X-Logged-In-Type",
		},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Auth-Token", "Content-Disposition"},
	}))

	router.Init(e, db) // ルーティング

	if sqlDb, err := db.DB(); err != nil {
		logger.Panic("Failed to get DB connection", err)
	} else {
		defer sqlDb.Close()
		logger.Log.Info("DB connection established")
	}

	echoLambda = echoadapter.New(e)
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}
