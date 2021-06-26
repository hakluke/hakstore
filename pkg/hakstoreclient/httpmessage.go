package hakstoreclient

// Message is a structure to use for success/error json response messages
type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
