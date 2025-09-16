package v1

type CreateUserBody struct {
	Email                string `json:"email" validate:"required,email"`
	DisplayName          string `json:"display_name"`
	Password             string `json:"password"`
	Passwordconfirmation string `json:"password_confirmation"`
}

type LoginUserBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
