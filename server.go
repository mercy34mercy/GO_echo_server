package main

import (
	"net/http"
	"os"
	"github.com/labstack/echo/v4"
)

// この関数を追加
func port() string {

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	return ":" + port
}


func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// Port番号を関数から取得
	e.Logger.Fatal(e.Start(port()))
}