package packet

import (
	"bytes"
	"testing"
)

func TestPackageEncoding(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 16*1024))
	p := &Packet{
		ID: "my-packet",
		Payload: Payload{
			"val1": "str",
			"val2": float64(123),
			"val3": true,
		},
	}
	_, err := p.WriteTo(buf)
	if err != nil {
		t.Fatal(err)
	}

	p2, err := ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	if p2 == nil {
		t.Fatal("ReadFrom: получен нулевой указатель")
	}

	if p2.ID != p.ID {
		t.Errorf("ID не совпали: %q получено, %q отправлено", p2.ID, p.ID)
	}

	for k, v := range p.Payload {
		v2, ok := p2.Payload[k]
		if !ok {
			t.Errorf("Отсутствует ключ %q", k)
		}
		if v2 != v {
			t.Errorf("Значение по ключу %q не совпадает: %v получено, %v отправлено", k, v2, v)
		}
	}
}
