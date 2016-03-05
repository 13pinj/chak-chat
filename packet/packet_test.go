package packet

import (
	"bytes"
	"encoding/json"
	"testing"
)

func comparePackets(t *testing.T, pa, pb *Packet) {
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

	r1, err := ReadFrom(&buf12)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if r1 == nil {
		t.Fatal("ReadFrom: нулевой указатель")
	}
	comparePackets(t, p1, r1)

	r2, err := ReadFrom(&buf12)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if r2 == nil {
		t.Fatal("ReadFrom: нулевой указатель")
	}
	comparePackets(t, p2, r2)

	r3, err := ReadFrom(&buf3)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if r3 == nil {
		t.Fatal("ReadFrom: нулевой указатель")
	}
	comparePackets(t, p3, r3)
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

func compareRequests(t *testing.T, pa, pb *Request) {
	if pa.Head != pb.Head {
		t.Errorf("Не совпадают заголовки: %q получено, %q отправлено", pb.Head, pa.Head)
	}
	ppa, ppb := pa.Packet, pb.Packet
	comparePackets(t, &ppa, &ppb)
}

func TestRequest_EncodingDecoding(t *testing.T) {
	p1 := &Request{
		Head: "hello",
		Packet: Packet{
			ID: "my-packet",
			Payload: Payload{
				"val1": "str",
				"val2": float64(-1),
				"val3": true,
			},
		},
	}
	p2 := &Request{
		Head: "hello",
		Packet: Packet{
			ID: "my-new-packet",
			Payload: Payload{
				"vala": "text",
				"valb": false,
				"valc": float64(3.14),
			},
		},
	}
	p3 := &Request{
		Head: "hello",
		Packet: Packet{
			ID: "000",
			Payload: Payload{
				"val_i":   false,
				"val_ii":  float64(2.71),
				"val_iii": "line",
			},
		},
	}

	buf1, buf2, buf3 := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	p1.WriteTo(buf1)
	p2.WriteTo(buf2)
	p3.WriteTo(buf3)

	rp1, _ := ReadFrom(buf1)
	r1, err := ToRequest(rp1)
	if r1 == nil || err != nil {
		t.Fatalf("ToRequest: неожиданная ошибка: %v", err)
	}
	compareRequests(t, p1, r1)

	rp2, _ := ReadFrom(buf2)
	r2, err := ToRequest(rp2)
	if r2 == nil || err != nil {
		t.Fatalf("ToRequest: неожиданная ошибка: %v", err)
	}
	compareRequests(t, p2, r2)

	rp3, _ := ReadFrom(buf3)
	r3, err := ToRequest(rp3)
	if r3 == nil || err != nil {
		t.Fatalf("ToRequest: неожиданная ошибка: %v", err)
	}
	compareRequests(t, p3, r3)
}

func TestRequest_NonReqPacket(t *testing.T) {
	p := &Packet{
		ID: "my-packet",
		Payload: Payload{
			"val1": "str",
			"val2": float64(-1),
			"val3": true,
		},
	}
	buf := &bytes.Buffer{}
	p.WriteTo(buf)

	r, _ := ReadFrom(buf)
	rr, err := ToRequest(r)
	if rr != nil || err == nil {
		t.Error("ToRequest: ожидается ошибка для пакета-незапроса")
	}
}

func compareResponses(t *testing.T, pa, pb *Response) {
	if pa.Status != pb.Status {
		t.Errorf("Не совпадают статусы: %q получено, %q отправлено", pb.Status, pa.Status)
	}
	ppa, ppb := pa.Packet, pb.Packet
	comparePackets(t, &ppa, &ppb)
}

func TestResponse_EncodingDecoding(t *testing.T) {
	p1 := &Response{
		Status: "ok",
		Packet: Packet{
			ID: "my-packet",
			Payload: Payload{
				"val1": "str",
				"val2": float64(-1),
				"val3": true,
			},
		},
	}
	p2 := &Response{
		Status: "err",
		Packet: Packet{
			ID: "my-new-packet",
			Payload: Payload{
				"vala": "text",
				"valb": false,
				"valc": float64(3.14),
			},
		},
	}
	p3 := &Response{
		Status: "internal",
		Packet: Packet{
			ID: "000",
			Payload: Payload{
				"val_i":   false,
				"val_ii":  float64(2.71),
				"val_iii": "line",
			},
		},
	}

	buf1, buf2, buf3 := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	p1.WriteTo(buf1)
	p2.WriteTo(buf2)
	p3.WriteTo(buf3)

	rp1, _ := ReadFrom(buf1)
	r1, err := ToResponse(rp1)
	if r1 == nil || err != nil {
		t.Fatalf("ToResponse: неожиданная ошибка: %v", err)
	}
	compareResponses(t, p1, r1)

	rp2, _ := ReadFrom(buf2)
	r2, err := ToResponse(rp2)
	if r2 == nil || err != nil {
		t.Fatalf("ToResponse: неожиданная ошибка: %v", err)
	}
	compareResponses(t, p2, r2)

	rp3, _ := ReadFrom(buf3)
	r3, err := ToResponse(rp3)
	if r3 == nil || err != nil {
		t.Fatalf("ToResponse: неожиданная ошибка: %v", err)
	}
	compareResponses(t, p3, r3)
}

func TestResponse_NonResPacket(t *testing.T) {
	p := &Packet{
		ID: "my-packet",
		Payload: Payload{
			"val1": "str",
			"val2": float64(-1),
			"val3": true,
		},
	}
	buf := &bytes.Buffer{}
	p.WriteTo(buf)

	r, _ := ReadFrom(buf)
	rr, err := ToResponse(r)
	if rr != nil || err == nil {
		t.Error("ToResponse: ожидается ошибка для пакета-незапроса")
	}
}
