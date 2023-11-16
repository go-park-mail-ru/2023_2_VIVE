package usecase

import "fmt"

var (
	INAPPROPRIATE_ROLE  = fmt.Errorf("Such a request is not possible with your role")
	CAN_NOT_READ_AVATAR = fmt.Errorf("An error occurred while loading the avatar")
)
