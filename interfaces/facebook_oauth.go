package interfaces

type FacebookAuth struct {
	Token string `json:"token" binding:"required"`
	ID    string `json:"id" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
