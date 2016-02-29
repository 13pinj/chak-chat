package packet

import "io"

type Payload map[string]interface{}

// Packet - структура необозначенного пакета.
type Packet struct {
	ID string
	// Полезная нагрузка.
	// Не содержит поле "id" и игнорирует его, если оно будет записано.
	Payload Payload
}

// ReadFrom декодирует пакет из потока и возвращает его структуру.
// Возвращает ошибки I/O и парсинга.
func ReadFrom(r io.Reader) (*Packet, error) {
	return nil, nil
}

// WriteTo кодирует пакет и записывает в `w`. Реализует интерфейс `io.WriterTo`.
func (p *Packet) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}
