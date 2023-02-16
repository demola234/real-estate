package interfaces

type ForgetPassword struct {
	Email string `json:"email" binding:"required"`
}