package interfaces

type Password struct {
	Password string `json:"password" binding:"required" min:"8"`
	Email    string `json:"email" binding:"required"`
}
