package vkbot

// Config contains configuration for bot
type Config struct {
	GroupID      int    `json:"group_id"`
	Confirmation string `json:"confirmation"`
	Host         string `json:"hostname"`
}
