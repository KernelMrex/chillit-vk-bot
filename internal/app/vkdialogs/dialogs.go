package vkdialogs

import "fmt"

// Dialogs provides map to auto load and store user dialogs
type Dialogs map[string]*Message

// Message Provides stucture for custom messages
type Message struct {
	Text string `json:"text"`
}

// GetText provides safe access to a text field for dialog
func (d Dialogs) GetText(dialogName string) (string, error) {
	dialog, ok := d[dialogName]
	if !ok {
		return "", fmt.Errorf("no such dialog '%s'", dialogName)
	}

	if dialog.Text == "" {
		return "", fmt.Errorf("dialog '%s' has empty text field", dialogName)
	}

	return dialog.Text, nil
}
