package usecase

import "fmt"

var (
	ErrInapropriateRole = fmt.Errorf("such a request is not possible with your role")
	ErrReadAvatar       = fmt.Errorf("an error occurred while loading the avatar")
)
