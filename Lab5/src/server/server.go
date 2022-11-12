package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

func find(str string, symb byte) int {
	for i := 0; i < len(str); i++ {
		if str[i] == symb {
			return (i)
		}
	}
	return -1
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		params := strings.Split(string(message), ")")
		log.Printf("recv: %s", string(message))
		x := 0
		y := 0
		z := 0
		den := 0
		for i := 0; i < (len(params) - 1); i++ {
			nums := strings.Split(params[i][find(params[i], '(')+1:], ",")
			mass, _ := strconv.Atoi(strings.TrimSpace(nums[0]))
			den += mass
			temp, _ := strconv.Atoi(strings.TrimSpace(nums[1]))
			x += temp * mass
			temp, _ = strconv.Atoi(strings.TrimSpace(nums[2]))
			y += temp * mass
			temp, _ = strconv.Atoi(strings.TrimSpace(nums[3]))
			z += temp * mass
		}
		res := "(" + fmt.Sprintf("%v", float32(x)/float32(den)) + ", " + fmt.Sprintf("%v", float32(y)/float32(den)) + ", " + fmt.Sprintf("%v", float32(z)/float32(den)) + ")"
		//fmt.Println(x, den,  float32(x)/float32(den))
		fmt.Printf("Центр масс: %s\n", res)
		err = c.WriteMessage(mt, []byte(res))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Println(http.ListenAndServe(*addr, nil))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
