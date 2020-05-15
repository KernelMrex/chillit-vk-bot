package vkapi

// Keyboard keyboard struct
type Keyboard struct {
	OneTime bool           `json:"one_time" url:"one_time"`
	Buttons []*KeyboardRow `json:"buttons" url:"buttons"`
}

// KeyboardRow keyboard row array of KeyboardButtons
type KeyboardRow []*KeyboardButton

// KeyboardButton keyboard button struct
type KeyboardButton struct {
	Action *ButtonAction `json:"action" url:"action"`
	Color  string        `json:"color" url:"color"`
}

// ButtonAction keyboard button action struct
type ButtonAction struct {
	Type    string `json:"type" url:"type"`
	Label   string `json:"label" url:"label"`
	Payload string `json:"payload" url:"payload"`
}
