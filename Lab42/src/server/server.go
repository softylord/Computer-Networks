package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gliderlabs/ssh"
	"log"
	"io"
	"os/exec"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		for {

			scanner := bufio.NewScanner(s)
			for scanner.Scan() { // use `for scanner.Scan()` to keep reading
				line := scanner.Text()
				fmt.Println("captured:", line)
				//com := strings.Split(line, "\n")
				cmd := exec.Command("sh", "-c", line)

				var out bytes.Buffer
				cmd.Stdout = &out

				err2 := cmd.Run()

				if err2 != nil {
					log.Fatal(err2)
				}

				// Print the output
				str := out.String()
				io.WriteString(s, fmt.Sprintf(str))
			}

		}
	})

	log.Println("starting ssh server on port 6060")
	log.Fatal(ssh.ListenAndServe(":6060", nil))
}
