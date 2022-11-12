package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

const (
	host     = "students.yss.su"
	database = "iu9networkslabs"
	user     = "iu9networkslabs"
	password = "Je2dTYr6"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type lett struct {
	User    string
	Mail    string
	Message string
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func main() {
	auth := smtp.PlainAuth("", "kateduzheeva@mail.ru", "aFEejxp2H44BFkGeyzyJ", "smtp.mail.ru")

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)

	// Initialize connection object.
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")

	rows, err := db.Query("select * from letters")
	checkError(err)
	defer rows.Close()

	letts := []lett{}

	for rows.Next() {
		p := lett{}
		err := rows.Scan(&p.User, &p.Mail, &p.Message)
		if err != nil {
			fmt.Println(err)
			continue
		}
		letts = append(letts, p)
	}
	c := 0
	l := len(letts)
	for _, p := range letts {
		c++
		body := template.Must(template.New("data").Parse(`
			<table bgcolor="#FBCEB1" border="0" cellpadding="0" cellspacing="0" style="margin:0; padding:0">
				<tr>
					<td>
						<center style="max-width: 600px; width: 100%;">
							<p><b>Здравствуйте, {{.User}}!</b></p>
							<p><i>{{.Message}}</i></p>
						</center>   
					</td>
				</tr>
			</table>`))
		buf := new(bytes.Buffer)
		body.Execute(buf, p)

		request := Mail{
			Sender:  "kateduzheeva@mail.ru",
			To:      /*[]string{"katy2jf@gmail.com"},*/[]string{p.Mail},
			Subject: "Дужеева Катя, ИУ9-31б",
			Body:    buf.String(),
		}

		msg := BuildMessage(request)

		err = smtp.SendMail("smtp.mail.ru:25", auth, "kateduzheeva@mail.ru", []string{p.Mail}, []byte(msg))
		checkError(err)
		fmt.Println("Письмо успешно отправлено!")

		//fmt.Println(buf.String())
		if c != l {
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(3)+1
			fmt.Printf("Sleeping %d minutes...\n", n)
			time.Sleep(time.Duration(n) * time.Minute)
			fmt.Println("Done")
		}
	}

}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
