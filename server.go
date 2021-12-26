package main

import (
	"net/http"
	"os"
	"github.com/labstack/echo/v4"
	"io/ioutil"
)

// この関数を追加
func port() string {

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "1234"
	}

	return ":" + port
}


func main() {

	url := "https://quiet-stream-64429.herokuapp.com/url"
	


	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		req, _ := http.NewRequest("GET", url, nil)
		client := new(http.Client)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
  
		byteArray, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(byteArray))
		return c.String(http.StatusOK, string(byteArray))
	})
	// Port番号を関数から取得
	e.Logger.Fatal(e.Start(port()))
}