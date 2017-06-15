package conn

type Connection interface {
	ReadLoop() error
	Write(bytes []byte) error
}
