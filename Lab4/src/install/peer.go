package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mgutz/logxi/v1"
	"github.com/skorobogatov/input"
)

type Peer struct {
	name, ip, port string
	parent         *Peer
	children       []Peer
}

// Client - состояние клиента.
type Client struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений
	res    string        // Текущая сбалансированность круглых скобок
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.

func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		res:    "",
	}
}

type Request struct {
	// Поле Command может принимать семь значений:
	// * "quit" - прощание с пиром(после этого пир рвёт соединение);
	// * "descendants" - показать имена всех потомки
	// * "parent" - показать имя родителя
	// вспомогательные команды (командной строкой не обрабатываются)
	// * "address" - ребенок отправляет родителю ip и порт, через которые он "слушает"
	// * "childQuit" - пир рвет соединение с родителем
	// * "becomeMe" - пир ставит на место себя своего ребенка и рвет с ним соединение
	// * "parentQuit" - пир рвет соединение с остельными детьми(дети переподключаются к новому родителю)
	Command string `json:"command"`

	// Если Command == "quit", "descendants", "parent" поле Data пустое.
	// Если Command == "childQuit", "becomeMe", "parentQuit", "address" в поле Data строка
	Data *json.RawMessage `json:"data"`
}

type Response struct {
	// Поле Status может принимать два значения:
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - результат выполнения запроса.
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	// Если Status == "result", в поле Data должен лежать результат: строка из имен потомков
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

func interactWithFamily(peerh Peer, command string) string {
	var addrStr string
	addrStr = peerh.ip + ":" + peerh.port
	res := ""
	// Разбор адреса, установка соединения с сервером и
	// запуск цикла взаимодействия с сервером.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		fmt.Printf("error %s: %v\n", addrStr, err)
	} else if conn, err := net.DialTCP("tcp", nil, addr); err != nil {
		fmt.Printf("error %s: %v\n", addrStr, err)
	} else {

		defer conn.Close()
		encoder, decoder := json.NewEncoder(conn), json.NewDecoder(conn)
		switch command {

		case "childQuit":
			send_request(encoder, command, peer.name)
		case "becomeMe":
			if peer.name != peer.parent.name {
				send_request(encoder, command, peer.parent.name+" "+peer.parent.ip+":"+peer.parent.port)
			} else {
				send_request(encoder, command, peerh.name+" "+peerh.ip+":"+peerh.port)
			}

		case "descendants":

			send_request(encoder, command, nil)

		case "parentQuit":
			p := peer.children[0]
			send_request(encoder, command, p.name+" "+p.ip+":"+p.port)
		}

		// Получение ответа.
		var resp Response
		if err := decoder.Decode(&resp); err != nil {
			fmt.Printf("error: %v\n", err)
			return ""
		}

		// Вывод ответа в стандартный поток вывода.
		switch resp.Status {

		case "failed":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var errorMsg string
				if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("failed: %s\n", errorMsg)
				}
			}
		case "result":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var string string
				if err := json.Unmarshal(*resp.Data, &string); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					res += string
					return res
				}
			}
		default:
			fmt.Printf("error: server reports unknown status %q\n", resp.Status)
		}
	}
	return ""
}

func (client *Client) serve() {
	defer client.conn.Close()

	decoder := json.NewDecoder(client.conn)

	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			if client.handleRequest2(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest2(req *Request) bool {
	switch req.Command {
	case "childQuit":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var name string
			if err := json.Unmarshal(*req.Data, &name); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("deliting child", "child", name)
				for i := 0; i < len(peer.children); i++ {
					if peer.children[i].name == name {
						peer.children = append(peer.children[:i], peer.children[i+1:]...)
						break
					}
				}

			}
		}
		if errorMsg == "" {
			//client.respond("ok", nil)
			return true
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "becomeMe":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var peerP string
			if err := json.Unmarshal(*req.Data, &peerP); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("reconnecting to", "peer", peerP)
				addrP := peerP[find(peerP, ' ')+1:]
				nameP := peerP[:find(peerP, ' ')]
				parent := Peer{name: nameP, ip: addrP[:find(addrP, ':')], port: addrP[find(addrP, ':')+1:]}
				peer.parent = &parent
				//addrStr := peer.ip + ":" + peer.port
				//подключение к родителю
				if peer.name != peer.parent.name {
					var addrPar string
					addrPar = peer.parent.ip + ":" + peer.parent.port
					if addrP, err := net.ResolveTCPAddr("tcp", addrPar); err != nil {
						fmt.Printf("ResolveTCPAddr error: %v\n", err)
					} else if conn, err := net.DialTCP("tcp", nil, addrP); err != nil {
						fmt.Printf("DialTCP error: %v\n", err)
					} else {
						for i := 0; i < 1; i++ {
							log.Info("sending address")
							address := peer.name + " " + peer.ip + ":" + peer.port + ")"
							encoder := json.NewEncoder(conn)

							send_request(encoder, "address", address)
						}
						log.Info("connected to parent", "adderess", addrP)

					}
				}

			}
		}
		if errorMsg == "" {
			//client.respond("ok", nil)
			return true
		} else {
			//client.logger.Error("addition failed", "reason", errorMsg)
			//client.respond("failed", errorMsg)
		}
	case "parentQuit":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var peerP string
			if err := json.Unmarshal(*req.Data, &peerP); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("reconnecting to", "peer", peerP)
				addrP := peerP[find(peerP, ' ')+1:]
				nameP := peerP[:find(peerP, ' ')]
				parent := Peer{name: nameP, ip: addrP[:find(addrP, ':')], port: addrP[find(addrP, ':')+1:]}
				peer.parent = &parent
				var addrPar string
				addrPar = peer.parent.ip + ":" + peer.parent.port
				if addrP, err := net.ResolveTCPAddr("tcp", addrPar); err != nil {
					fmt.Printf("ResolveTCPAddr error: %v\n", err)
				} else if conn, err := net.DialTCP("tcp", nil, addrP); err != nil {
					fmt.Printf("DialTCP error: %v\n", err)
				} else {
					for i := 0; i < 1; i++ {
						log.Info("sending address")
						address := peer.name + " " + peer.ip + ":" + peer.port + ")"
						encoder := json.NewEncoder(conn)

						send_request(encoder, "address", address)
					}
					log.Info("connected to parent", "adderess", addrP)

				}

			}
		}
		if errorMsg == "" {
			return true
		}
	case "address":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var address string
			if err := json.Unmarshal(*req.Data, &address); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("reciving address", "address", address)
				if len(address) != 0 {
					name := address[:find(address, ' ')]
					addr := address[find(address, ' ')+1 : find(address, ')')]
					par := peer.parent.ip + ":" + peer.parent.port
					if addr != par && find(addr, ':') != 0 {

						var ch Peer
						ch.name = name
						ch.ip = addr[:find(addr, ':')]
						ch.port = addr[find(addr, ':')+1:]
						peer.children = append(peer.children, ch)
					}
				}
			}
		}
		if errorMsg == "" {
		}
	case "descendants":
		client.res = client.res + "\n" + peer.name
		fl = 1
		for i := 0; i < len(peer.children); i++ {
			client.res += interactWithFamily(peer.children[i], "descendants")
		}
		client.respond("result", client.res)
		return true
	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

var fl int

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&Response{status, &raw})
}

// interact - функция, содержащая цикл взаимодействия с сервером.
func interact() {
	fl = 0
	decoder := json.NewDecoder(os.Stdout)
	for {
		// Чтение команды из стандартного потока ввода
		command := input.Gets()

		// Отправка запроса.
		switch command {
		case "quit":
			if len(peer.children) != 0 {
				interactWithFamily(peer.children[0], "becomeMe")
				for i := 1; i < len(peer.children); i++ {
					log.Info("quiting is in progress")
					interactWithFamily(peer.children[i], "parentQuit")
				}
			}
			if peer.name != peer.parent.name {
				interactWithFamily(*peer.parent, "childQuit")
			}
			os.Exit(0)
		case "parent":
			fmt.Println(peer.parent.name)
			interact()
		case "descendants":
			if len(peer.children) == 0 {
				fmt.Println("No descendants")
			} else {
				res := ""
				for i := 0; i < len(peer.children); i++ {

					res += interactWithFamily(peer.children[i], command)
				}
				if fl == 0 {
					fmt.Printf("Descendants are: %s\n", res)
				}
			}
			interact()
		default:
			fmt.Printf("error: unknown command\n")
			interact()
		}

		// Получение ответа.
		var resp Response
		if err := decoder.Decode(&resp); err != nil {
			fmt.Printf("error: %v\n", err)
			//break
		}

		// Вывод ответа в стандартный поток вывода.
		switch resp.Status {
		case "ok":
			fmt.Printf("okk\n")
		case "failed":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var errorMsg string
				if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("failed: %s\n", errorMsg)
				}
			}
			interact()
		case "result":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var string string
				if err := json.Unmarshal(*resp.Data, &string); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("result: %s\n", string)
				}
			}
			interact()
		default:
			fmt.Printf("error: server reports unknown status %q\n", resp.Status)
		}
	}

}

// send_request - вспомогательная функция для передачи запроса с указанной командой
// и данными. Данные могут быть пустыми (data == nil).
func send_request(encoder *json.Encoder, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&Request{command, &raw})
}

func find(str string, a byte) int {
	for i := 0; i < len(str); i++ {
		if str[i] == a {
			return i
		}
	}
	return 0
}

var peer Peer

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	var name, ip, port, nameP, ipP, portP string
	fmt.Println("Name, IP and port:")
	fmt.Scanf("%s\n%s\n%s\n", &name, &ip, &port)
	fmt.Println("Name, IP and port of Parent:")
	fmt.Scanf("%s\n%s\n%s\n", &nameP, &ipP, &portP)
	peerP := Peer{name: nameP, ip: ipP, port: portP}
	peer = Peer{name: name, ip: ip, port: port, parent: &peerP}
	peer.children = make([]Peer, 0, 100)

	flag.StringVar(&addrStr, "addr", peer.ip+":"+peer.port, "specify ip address and port")
	flag.Parse()
	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		//подключение к родителю
		if peer.name != peer.parent.name {
			var addrPar string
			addrPar = peer.parent.ip + ":" + peer.parent.port
			if addrP, err := net.ResolveTCPAddr("tcp", addrPar); err != nil {
				fmt.Printf("ResolveTCPAddr error: %v\n", err)
			} else if conn, err := net.DialTCP("tcp", nil, addrP); err != nil {
				fmt.Printf("DialTCP error: %v\n", err)
			} else {
				for i := 0; i < 1; i++ {
					log.Info("sending address")
					address := peer.name + " " + peer.ip + ":" + peer.port + ")"
					encoder := json.NewEncoder(conn)

					send_request(encoder, "address", address)
				}
				log.Info("connected to parent", "adderess", addrP)

			}
		}

		go func() {
			for {
				interact()
			}
		}()
		/*go func() {
			if peer.name != peer.parent.name {
				interactWithFamily(*peer.parent, "smile")
			}
		}()*/
		// Инициация слушания сети на заданном адресе.
		go func() {
			if listener, err := net.ListenTCP("tcp", addr); err != nil {
				log.Error("listening failed", "reason", err)
			} else {
				// Цикл приёма входящих соединений.
				for {
					if conn, err := listener.AcceptTCP(); err != nil {
						log.Error("cannot accept connection", "reason", err)
					} else {
						str := conn.RemoteAddr().String()
						log.Info("accepted connection", "address", str)
						//go handleConnection(conn)

						// Запуск go-программы для обслуживания клиентов.
						go NewClient(conn).serve()
					}
				}
			}
		}()
	}
	<-time.After(time.Minute * 120)
}
