package vkbot

type messageActionObject struct {
	Message *messageObject `json:"message"`
}

type messageObject struct {
	ID      int    `json:"id"`
	Time    uint64 `json:"date"`
	From    int    `json:"from_id"`
	Text    string `json:"text"`
	Payload string `json:"payload"`
}
