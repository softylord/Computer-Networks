package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"bufio"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	var (
		username, password, hostname, port string
	)
	fmt.Print("ftp-host: ")
	fmt.Scan(&hostname)
	fmt.Print("port: ")
	fmt.Scan(&port)
	fmt.Print("login: ")
	fmt.Scan(&username)
	fmt.Print("password: ")
	fmt.Scan(&password)

	c, err := ftp.Dial(hostname+":"+port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		coms := scanner.Text()
		cmd := strings.Split(coms, " ")

		switch cmd[0] {
		case "QUIT":
			if err := c.Quit(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)

		case "STOR":
			f, err := os.Open(cmd[1])
			if err != nil {
				panic(err)
			}
			defer f.Close()

			name := findName(cmd[1])
			err = c.Stor(name, f)
			if err != nil {
				panic(err)

			}

		case "RETR":
			r, err := c.Retr(cmd[1])
			if err != nil {
				panic(err)
			}
			defer r.Close()

			buf, err := io.ReadAll(r)
			file, err := os.Create("./src/" + cmd[1])
			if err != nil {
				fmt.Println("Unable to create file:", err)
			}
			defer file.Close()
			file.WriteString(string(buf))

		case "MKD":
			err := c.MakeDir(cmd[1])
			if err != nil {
				panic(err)
			}

		case "DELE":
			err := c.Delete(cmd[1])
			if err != nil {
				panic(err)
			}
		case "RMD":
			err:=c.RemoveDir(cmd[1])
			if err!=nil{
				panic(err)
			}

		case "LIST":
			list := ""
			r, err := c.List(cmd[1])
			if err != nil {
				panic(err)
			}
			fmt.Println(len(r))
			for _, elem := range r {
				list = list + elem.Name + " "
			}
			fmt.Println(list)

		case "NEWS":
			ti := time.Now().Format(time.RFC1123Z)
			t, _ := time.Parse(time.RFC1123Z, ti)
			news:=""
			r, _ := c.NameList("./")
			for i, elem := range r {
				temp := strings.Split(r[i], " ")
				if temp[0] == "Дужеева" {
					reader, err := c.Retr(elem)
					if err != nil {
						panic(err)
					}
					buf, err := io.ReadAll(reader)
					news+=string(buf)
					reader.Close()
				}
			}
			text := ""

			fp := gofeed.NewParser()
			feed, _ := fp.ParseURL(cmd[1])
			for _, item := range feed.Items {

				if !Find(news, item.Title) {
					pub, _ := time.Parse(time.RFC1123Z, item.Published)
					text += "Title: " + item.Title + "\n" + "Publication date: " + pub.String()[:10] + " " + pub.String()[11:16] + "\n" +
						"Autor: " + item.Author.Name + "\n" + item.Description + "\n\n"

				}
			}
			if text!=""{

			file, err := os.Create("./src/news " + t.String() + ".txt")
			if err != nil {
				panic(err)
			}
			file.WriteString(text)
			file.Close()
			f, _:= os.Open("./src/news " + t.String() + ".txt")

			c.Stor("Дужеева Катя"+" "+t.String()[:10]+" "+t.String()[11:16]+".txt", f)
		}

		default:
			fmt.Println("Unknown command, try again")
		}

	}
}

func Find(text, title string) bool {
	data := strings.Split(text, "\n")
	for i := 0; i < len(data); i++ {
		if len(data[i]) > 5 && data[i][:5] == "Title" {
			if data[i][7:] == title {
				return true
			}
		}
	}
	return false
}

func findName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}