package connmngr

type Connection interface {
	Send(msg string) error
	Receive() (msg string, err error)
	Close() error
}

type ConnectionFactory interface {
	CreateConnection() (Connection, error)
}
