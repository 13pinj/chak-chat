# Архитектура

Структура пакетов:
+ `packet` - отправка и чтение пакетов по сети
+ `server` - серверные исходники
+ `client` - клиентские исходники

## `packet`

### Пакет

```go
type Packet struct {
  ID string
  // Полезная нагрузка.
  // Не содержит поле "id" и игнорирует его, если оно будет записано.
  Payload map[string]interface{}
}
```
Packet - структура необозначенного пакета.

```go
func ReadFrom(r io.Reader) (*Packet, error)
```
ReadFrom декодирует пакет из потока и возвращает его структуру.
Возвращает ошибки I/O и парсинга.

```go
func (p *Packet) WriteTo(w io.Writer) (n int64, err error)
```
WriteTo кодирует пакет и записывает в `w`. Реализует интерфейс `io.WriterTo`.

### Запрос и ответ

```go
type Request struct {
  Packet
  // Поле "req". Редактирование заголовка запроса должно происходить
  // только через это поле: значение Payload["req"] будет игнорировано.
  Head string
}
```
Request - структура запроса.

```go
func ToRequest(p *Packet) (*Request, error)
```
ToRequest проверяет, является ли пакет запросом. В случае успеха, возвращает
структуру запроса, эквивалентную пакету. В случае неуспеха, возвращает в тексте
ошибки критерий, по которому пакет не прошел проверку.

```go
func (req *Request) WriteTo(w io.Writer) (n int64, err error)
```
WriteTo кодирует запрос и записывает в `w`. Реализует интерфейс `io.WriterTo`.

```go
type Response struct {
  Packet
  // Редактирование статуса должно происходить
  // только через это поле: значение Payload["status"] будет игнорировано.
  Status string
}
```
Response - структура ответа.

```go
func ToResponse(p *Packet) (*Response, error)
```
ToResponse проверяет, является ли пакет ответом. В случае успеха, возвращает
структуру ответа, эквивалентную пакету. В случае неуспеха, возвращает в тексте
ошибки критерий, по которому пакет не прошел проверку.

```go
func (p *Response) WriteTo(w io.Writer) (n int64, err error)
```
WriteTo кодирует ответ и записывает в `w`. Реализует интерфейс `io.WriterTo`.
