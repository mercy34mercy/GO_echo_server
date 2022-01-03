package main

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	// "github.com/heroku/x/hmetrics/onload"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// この関数を追加
func port() string {

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	return ":" + port
}


//from sql server
type sqlimage struct {
	Alt string `db:"alt"`
	Src string `db:"url`
}

/** JSONデコード用に構造体定義 */
type data struct {
	Alt string `json:"alt"`
	Src string `json:"url"`
}

type image struct {
	Data []data `json:"data"`
}

/*デコードするJson*/


func main() {
	//heroku psqlのURL
	DATABASE_URL := `postgres://vxuvjuiftslcyt:adc9ab1d3939a492978975d987ab1fb58e853dc8991496bd62a3257eef646de3@ec2-52-70-205-234.compute-1.amazonaws.com:5432/d6h6nu23bkqeej`
	//美女検索APIのURL
	url := "https://quiet-stream-64429.herokuapp.com/url"

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		//APIにGETリクエストを送信
		req, _ := http.NewRequest("GET", url, nil)
		client := new(http.Client)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		
		//sql Open
		db, err := sql.Open("postgres", DATABASE_URL)
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}

		// JSONデコード
		var images image
		if err := json.Unmarshal(byteArray, &images); err != nil {
			log.Fatal(err)
			fmt.Printf("erroo")
		}
		// デコードしたデータを表示
		for _, p := range images.Data {
			fmt.Printf("%s : %s\n", p.Alt, p.Src)
			if _, err := db.Exec("insert into beautifulimage (alt,src) values($1,$2);", p.Alt, p.Src); err != nil {

				fmt.Printf("Error insert database table: %q", err)
			}

		}

		return c.String(http.StatusOK, "success")
	})


	//DB作成
	e.GET("/createDB", func(c echo.Context) error {
		db, err := sql.Open("postgres", os.Getenv(DATABASE_URL))
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}

		if _, err := db.Exec("CREATE TABLE beautifulimage (alt text,src text unique);"); err != nil {

			fmt.Printf("Error creating database table: %q", err)
			return c.String(http.StatusOK, "Error creating database table")
		}
		return c.String(http.StatusOK, "success")
	})

	//DBからrandomに一つ
	e.GET("/images", func(c echo.Context) error {
		//psql Open
		db, err := sql.Open("postgres", DATABASE_URL)
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}

		//sqlにrequest送信
		if result, err := db.Query("SELECT * FROM beautifulimage ORDER BY random() LIMIT 1;"); err != nil {
			fmt.Printf("Error get database table: %q", err)
			return c.String(http.StatusOK, "Error get database table")
		} else {
			var a sqlimage
			for result.Next() {
				result.Scan(&a.Alt, &a.Src)
				fmt.Printf("alt: %s, src: %s\n", a.Alt, a.Src)
				r := data{string(a.Alt),string(a.Src)}
				
				//構造体→Json
				res, _ := json.Marshal(r)

				return c.String(http.StatusOK, string(res))
			}
		}

		return c.String(http.StatusOK, "success")
	})
	// Port番号を関数から取得
	e.Logger.Fatal(e.Start(port()))
}
