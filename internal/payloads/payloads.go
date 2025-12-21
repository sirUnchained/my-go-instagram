package payloads

type UserPayload struct {
	Username string `json:"username" validate:"required,min=8,max=255"`
	Fullname string `json:"fullname" validate:"required,min=8,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}
