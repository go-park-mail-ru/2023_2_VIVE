package serviceErrors

import "fmt"

var (
	ErrOpenConn          = fmt.Errorf("could not open websocket connection")
	ErrInvalidUserID     = fmt.Errorf("invalid user_id")
	ErrConnAlreadyExists = fmt.Errorf("connection already exists")
	ErrNoConn            = fmt.Errorf("no connection")
	ErrInvalidConnection = fmt.Errorf("invalid stored connection")
)
