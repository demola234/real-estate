package interfaces

type PushNotification struct {
	To     string   `json:"to"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Tokens []string `json:"tokens"`
}