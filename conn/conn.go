package conn

type Connection interface {
	ReadLoop() error
	WriteLoop() error
	Write(bytes []byte) error
	GetChanOut() chan<- []byte
	Ping() error
	IsClosed() bool
}
