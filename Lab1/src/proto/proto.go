package proto

import "encoding/json"

// Request -- запрос клиента к серверу.
type Request struct {
	// Поле Command может принимать три значения:
	// * "quit" - прощание с сервером (после этого сервер рвёт соединение);
	// * "check" - сервер принимает строку и проверяет ее на сбалансированность круглых скобок;
	Command string `json:"command"`

	// Если Command == "check", в поле Data должна лежать строка
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Response -- ответ сервера клиенту.
type Response struct {
	// Поле Status может принимать три значения:
	// * "ok" - успешное выполнение команды "quit";
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - сбалансированность круглых скобок вычислена.
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	// Если Status == "result", в поле Data должен лежать результат: "ok", "smth is wrong"
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Fraction -- дробь.
/*type Fraction struct {
	// Числитель (в десятичной системе, разрешён знак).
	Numerator string `json:"num"`

	// Знаменатель (в десятичной системе).
	Denominator string `json:"denom"`
}*/