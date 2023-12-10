package serviceErrors

import "fmt"

var (
	ErrOpenConn          = fmt.Errorf("could not open websocket connection")
	ErrHandshakeMsg      = fmt.Errorf("invalid handshake message")
	ErrConnAlreadyExists = fmt.Errorf("connection already exists")
	ErrNoConn            = fmt.Errorf("no connection")
)
