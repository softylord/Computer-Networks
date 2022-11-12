package main

import (
	"database/sql"
	"fmt"
	
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/mmcdole/gofeed"
)

const (
	host     = "students.yss.su"
	database = "iu9networkslabs"
	user     = "iu9networkslabs"
	password = "Je2dTYr6"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Find(a []string, title string) bool{
	for i:=0; i<len(a); i++{
		if a[i]==title{
			return true
		}
	}
	return false
}

//func Month(m []string) int

func main() {

	// Initialize connection string.
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)

	// Initialize connection object.
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")

	// Modify some data in table.
	sqlStatement, err := db.Prepare("INSERT INTO iu9duzheeva (title, text, date, time, author) VALUES (?, ?, ?, ?, ?);")

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://news.rambler.ru/rss/technology/")
	for _, item := range feed.Items {
		rows, err := db.Query("SELECT title from iu9duzheeva;")
		checkError(err)
		defer rows.Close()
		var arr []string
		for rows.Next() {
			var title string
			err := rows.Scan(&title)
			arr = append(arr, title)
			checkError(err)
		}
		if !Find(arr, item.Title) {
			pub, _:=time.Parse(time.RFC1123Z, item.Published)
			res, err := sqlStatement.Exec(item.Title, item.Description, pub.String()[:10], pub.String()[11:18], item.Author.Name)

			checkError(err)

			rowCount, err := res.RowsAffected()
			fmt.Printf("Inserted %d row(s) of data.\n", rowCount)

		}
	}


}
