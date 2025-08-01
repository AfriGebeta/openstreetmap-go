package v1

type UserCreate struct {
	Email                string `json:"email" validate:"required,email"`
	DisplayName          string `json:"display_name" validate:"required,max=60"`
	Password             string `json:"password" validate:"required,max=60"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,max=60"`
}
