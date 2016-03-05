package packet

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPacket_EncodingDecoding(t *testing.T) {
	var buf12 bytes.Buffer
	var buf3 bytes.Buffer
	p1 := &Packet{
		ID: "my-packet",
		Payload: Payload{
			"val1": "str",
			"val2": float64(-1),
			"val3": true,
		},
	}
	p2 := &Packet{
		ID: "my-new-packet",
		Payload: Payload{
			"vala": "text",
			"valb": false,
			"valc": float64(3.14),
		},
	}
	p3 := &Packet{
		ID: "000",
		Payload: Payload{
			"val_i":   false,
			"val_ii":  float64(2.71),
			"val_iii": "line",
		},
	}

	_, err := p1.WriteTo(&buf12)
	if err != nil {
		t.Fatalf("p1.WriteTo: %v", err)
	}
	_, err = p2.WriteTo(&buf12)
	if err != nil {
		t.Fatalf("p2.WriteTo: %v", err)
	}
	_, err = p3.WriteTo(&buf3)
	if err != nil {
		t.Fatalf("p3.WriteTo: %v", err)
	}

	compare := func(pa, pb *Packet) {
		if pb.ID != pa.ID {
			t.Errorf("ID не совпали: %q получено, %q отправлено", pb.ID, pa.ID)
		}

		for k, v := range pa.Payload {
			v2, ok := pb.Payload[k]
			if !ok {
				t.Errorf("Отсутствует ключ %q", k)
			}
			if v2 != v {
				t.Errorf("Значение по ключу %q не совпадает: %v получено, %v отправлено", k, v2, v)
			}
		}
	}

	r1, err := ReadFrom(&buf12)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if r1 == nil {
		t.Fatal("ReadFrom: нулевой указатель")
	}
	compare(p1, r1)

	r2, err := ReadFrom(&buf12)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if r2 == nil {
		t.Fatal("ReadFrom: нулевой указатель")
	}
	compare(p2, r2)

	r3, err := ReadFrom(&buf3)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if r3 == nil {
		t.Fatal("ReadFrom: нулевой указатель")
	}
	compare(p3, r3)
}

func TestPacket_EdgeCases(t *testing.T) {
	brokenJSON := `{"val1": 123, "val2":}`
	buf1 := bytes.NewBufferString(brokenJSON)
	p, err := ReadFrom(buf1)
	if p != nil || err == nil {
		t.Error("ReadFrom: должен сообщать об ошибках парсинга")
	}

	p = &Packet{
		ID: "my-packet",
		Payload: Payload{
			"val1": "str",
			"val2": float64(-1),
			"val3": true,
		},
	}
	jb, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	buf2 := bytes.NewBuffer(jb)
	p, err = ReadFrom(buf2)
	if p != nil || err == nil {
		t.Error("ReadFrom: должен сообщать о неожиданном конце потока")
	}
}
