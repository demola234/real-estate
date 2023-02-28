package interfaces

type PushNotificationToAll struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Tokens []string `json:"tokens"`
}


type PushNotificationToUser struct {

	Title string `json:"title"`
	Body  string `json:"body"`
	To    string `json:"to"`
}