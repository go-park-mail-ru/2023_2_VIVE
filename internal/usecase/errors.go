package usecase

import "fmt"

var (
	ErrInapropriateRole = fmt.Errorf("Such a request is not possible with your role")
	ErrReadAvatar       = fmt.Errorf("An error occurred while loading the avatar")
	BadAvatarSize       = fmt.Errorf("The uploaded file must be 2MB or less in size")
	BadAvatarType       = fmt.Errorf("The uploaded file must be of the jpeg, png or gif type")
	ErrForbidden        = fmt.Errorf("Forbidden request")
)
