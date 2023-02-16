package interfaces

type FaceResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Picture Picture `json:"picture"`
}

type Picture struct {
	Data Data `json:"data"`
}

type Data struct {
	Height       int64  `json:"height"`
	IsSilhouette bool   `json:"is_silhouette"`
	URL          string `json:"url"`
	Width        int64  `json:"width"`
}
