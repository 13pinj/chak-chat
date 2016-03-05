package packet

import (
	"encoding/json"
	"io"
)

type Payload map[string]interface{}

// Packet - структура необозначенного пакета.
type Packet struct {
	ID string
	// Полезная нагрузка.
	// Не содержит поле "id" и игнорирует его, если оно будет записано.
	Payload Payload
}

//TODO:  Json функция marchal,unmarchal; их использовааать
// ReadFrom декодирует пакет из потока и возвращает его структуру.
// Возвращает ошибки I/O и парсинга.
func ReadFrom(r io.Reader) (*Packet, error) {
	var buf []byte
	slice := make([]byte, 1)
	for {
		_, err := r.Read(slice)
		if err != nil {
			return nil, err
		}
		if slice[0] == 0 {
			break
		}
		buf = append(buf, slice[0])
	}
	p := &Packet{}
	err := json.Unmarshal(buf, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// WriteTo кодирует пакет и записывает в `w`. Реализует интерфейс `io.WriterTo`.
func (p *Packet) WriteTo(w io.Writer) (n int64, err error) {
	buf, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}
	buf = append(buf, 0)
	v, err := w.Write(buf)
	return int64(v), err
}

// Request - структура запроса.
type Request struct {
	Packet
	// Поле "req". Редактирование заголовка запроса должно происходить
	// только через это поле: значение Payload["req"] будет игнорировано.
	Head string
}

// ToRequest проверяет, является ли пакет запросом. В случае успеха, возвращает
// структуру запроса, эквивалентную пакету. В случае неуспеха, возвращает в тексте
// ошибки критерий, по которому пакет не прошел проверку.
func ToRequest(p *Packet) (*Request, error) {
	return nil, nil
}

// WriteTo кодирует запрос и записывает в `w`. Реализует интерфейс `io.WriterTo`.
func (req *Request) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}
