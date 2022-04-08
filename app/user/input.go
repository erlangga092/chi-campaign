package user

type RegisterUserInput struct {
	Name       string `json:"name" validate:"required"`
	Occupation string `json:"occupation" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
}

type CheckEmailAvailableInput struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
