package vkbot

// Config contains configuration for bot
type Config struct {
	GroupID      int    `json:"group_id" yaml:"group_id"`
	Confirmation string `json:"confirmation_code" yaml:"confirmation_code"`
	Token        string `json:"token" yaml:"token"`
	Host         string `json:"host" yaml:"host"`
	DialogsPath  string `json:"dialogs_path" yaml:"dialogs_path"`
}
