package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func find(str string, ch byte) int {
	for i := 0; i < len(str); i++ {
		if str[i] == ch {
			return i
		}
	}
	return -1
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	str := strings.Split(r.URL.Path, "/")

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if len(str) > 1 {
		name := []string{}
		//p := ""
		pars := []string{}

		if len(str) > 2 {
			name = strings.Split(str[1], ".")
			pars = str[2:]
		} else {
			name = strings.Split(str[1], ".")
			//fmt.Fprintln(w, str[1]+" a")
			values := r.URL.Query()
			for _, v := range values {
				pars = append(pars, v[0])
			}
			/*for i, j := 0, len(pars)-1; i < j; i, j = i+1, j-1 {
				pars[i], pars[j] = pars[j], pars[i]
			}*/
		}

		switch name[1] {
		case "txt":
			fileBytes, err := ioutil.ReadFile("./" + name[1] + "/" + str[1])
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/txt; charset=utf-8")
			w.Write(fileBytes)
		case "html":
			fileBytes, err := ioutil.ReadFile("./" + name[1] + "/" + str[1])
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, string(fileBytes))

		case "jpeg":
			fileBytes, err := ioutil.ReadFile("./" + name[1] + "/" + str[1])
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(fileBytes)
		case "png":
			fileBytes, err := ioutil.ReadFile("./" + name[1] + "/" + str[1])
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "image/png")
			w.Write(fileBytes)

		case "gif":
			fileBytes, err := ioutil.ReadFile("./" + name[1] + "/" + str[1])
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "image/gif")
			w.Write(fileBytes)

		case "cpp":
			cmd := exec.Command("g++", "-o", "./"+name[1]+"/"+name[0], "./"+name[1]+"/"+str[1])
			cmd.Output()
			//fmt.Fprint(w, "g++ -o "+"./"+name[1]+"/"+name[0]+" "+"./"+name[1]+"/"+str[1])
			cmd = exec.Command("sh", "-c", "./"+name[1]+"/"+name[0])

			var out bytes.Buffer
			cmd.Stdout = &out
			stdin, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal(err)
			}

			go func() {
				defer stdin.Close()
				for i := 0; i < len(pars); i++ {
					io.WriteString(stdin, pars[i]+"\n")

				}

				//io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
			}()

			err2 := cmd.Run()

			if err2 != nil {
				log.Fatal(err2)
			}

			w.Write(out.Bytes())
		}
	}

}

func main() {

	http.HandleFunc("/", helloHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
