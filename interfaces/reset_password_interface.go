package interfaces

type ResendOtp struct {
	Email string `json:"email" binding:"required"`
}