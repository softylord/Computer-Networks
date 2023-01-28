package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"html/template"
	"net/http"
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
		msg := string(message)
		if msg == "1" {
			go func() {
				c.WriteMessage(mt, []byte("Searching for new files..."))

			}()
		} else if msg == "2" {
			go func() {
				
				c.WriteMessage(mt, []byte("hi"))
			}()
		} else {
			go func() {
				c.WriteMessage(mt, []byte("hahahahahahah"))
			}()
		}

	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
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
        window.addEventListener("load", function (evt) {
            var ws1;
            var ws2;
            var ws3;
            var print1 = function (message) {
                var d = document.createElement("div");
                d.textContent = message;
                if (output1.hasChildNodes()) {

                    output1.removeChild(output1.childNodes[0]);
                }
                output1.appendChild(d);

            };
            var print2 = function (message) {
                var d = document.createElement("div");
                d.textContent = message;
                if (output2.hasChildNodes()) {

                    output2.removeChild(output2.childNodes[0]);
                }
                output2.appendChild(d);

            };
            var print3 = function (message) {
                var d = document.createElement("div");
                d.textContent = message;
                if (output3.hasChildNodes()) {

                    output3.removeChild(output3.childNodes[0]);
                }
                output3.appendChild(d);

            };
            ws1 = new WebSocket("{{.}}");
            ws1.onopen = function (evt) {
                while (1 == 1) {
                    ws1.send("1");
                    return false;
                }
            }
            ws1.onclose = function (evt) {
                print1("CLOSE");
                ws1 = null;
            }
            ws1.onmessage = function (evt) {
                print1(evt.data);

                ws1.send("1");
            }
            ws2 = new WebSocket("{{.}}");
            ws2.onopen = function (evt) {
                while (1 == 1) {
                    ws2.send("2");
                    return false;
                }

            }
            ws2.onclose = function (evt) {
                print2("CLOSE");
                ws2 = null;
            }
            ws2.onmessage = function (evt) {
                print2(evt.data);
            }
            ws3 = new WebSocket("{{.}}");
            ws3.onopen = function (evt) {
                while (1 == 1) {
                    ws3.send("3");
                    return false;
                }
            }
            ws3.onclose = function (evt) {
                print3("CLOSE");
                ws3 = null;
            }
            ws3.onmessage = function (evt) {
                print3(evt.data);
            }
            return false;
        });
    </script>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
      <title>Dashboard</title>
      <style type="text/css">
        .layout {
            overflow: hidden;
            /* Отмена обтекания */
        }

        .col1,
        .col2,
        .col3 {
            width: 33.33%;
            /* Ширина колонок */
            float: left;
            /* Создаем колонки */
        }

        .layout div div {
            margin: 0 10px;
            /* Отступы */
            padding: 10px;
            /* Поля */
            height: 1000px;
            /* Высота колонок */
            background: #4f703e;
            /* Цвет фона */
            color: #f5e8d0;
            /* Цвет текста */
            overflow: auto;
        }
    </style>
     
</head>
 

<body>
      <div class="layout">
           <div class="col1">
            <div id="output1" style="max-height: 70vh;"></div>
               
        </div>
           <div class="col2">
            <div id="output2" style="max-height: 70vh;"></div>
               
        </div>
           <div class="col3">
            <div id="output3" style="max-height: 70vh;"></div>
               
        </div>
          </div>
     </body>

</html>
`))