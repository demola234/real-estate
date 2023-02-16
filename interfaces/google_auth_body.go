package interfaces

type GoogleAuth struct {
	Token string `json:"token" binding:"required"`
	Code      string `json:"code" binding:"required"`
}
