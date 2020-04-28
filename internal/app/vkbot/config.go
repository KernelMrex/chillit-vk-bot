package vkbot

// Config contains configuration for bot
type Config struct {
	GroupID      int    `json:"group_id" yaml:"group_id"`
	Confirmation string `json:"confirmation_code" yaml:"confirmation_code"`
	Host         string `json:"host" yaml:"host"`
}
