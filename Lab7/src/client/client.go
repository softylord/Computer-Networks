package main

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func main() {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "kateduzheeva@mail.ru", "aFEejxp2H44BFkGeyzyJ", "smtp.mail.ru")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	var t string
	fmt.Println("Message To:")
	fmt.Scanln(&t)
	fmt.Println("Message Subject:")
	in := bufio.NewReader(os.Stdin)
	subj, _ := in.ReadString('\n')
	fmt.Println("Message Body:")
	body, _ := in.ReadString('\n')
	to := []string{t}
	msg := []byte("To: " + t + "\r\n" +
		"Subject: " + subj + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail("smtp.mail.ru:25", auth, "kateduzheeva@mail.ru", to, msg)
	if err != nil {
		log.Fatal(err)
	} else{
		fmt.Println("Письмо успешно отправлено!")
	}
}
