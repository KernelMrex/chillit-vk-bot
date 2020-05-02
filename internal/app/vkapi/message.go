package vkapi

// Message structure of message data in request
type Message struct {
	RecieverID int    `url:"user_id"`
	RandomID   int64  `url:"random_id"`
	Text       string `url:"message"`
}
