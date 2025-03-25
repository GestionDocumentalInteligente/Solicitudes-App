package transport

type Req struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}
